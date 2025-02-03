package mac_prefs

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"fmt"
)

// UserType represents the type of user for preferences
type UserType string

// HostType represents the type of host for preferences
type HostType string

// PreferenceScope defines the scope for preferences
type PreferenceScope struct {
	User UserType
	Host HostType
}

var (
	// CurrentUser represents the current user's preferences
	CurrentUser UserType = "kCFPreferencesCurrentUser"
	// AnyUser represents preferences for any user
	AnyUser UserType = "kCFPreferencesAnyUser"
	// CurrentHost represents the current host's preferences
	CurrentHost HostType = "kCFPreferencesCurrentHost"
	// AnyHost represents preferences for any host
	AnyHost HostType = "kCFPreferencesAnyHost"

	// CurrentUserCurrentHost represents preferences for the current user on the current host
	CurrentUserCurrentHost = PreferenceScope{User: CurrentUser, Host: CurrentHost}
	// CurrentUserAnyHost represents preferences for the current user on any host
	CurrentUserAnyHost = PreferenceScope{User: CurrentUser, Host: AnyHost}
	// AnyUserCurrentHost represents preferences for any user on the current host
	AnyUserCurrentHost = PreferenceScope{User: AnyUser, Host: CurrentHost}
	// AnyUserAnyHost represents preferences for any user on any host
	AnyUserAnyHost = PreferenceScope{User: AnyUser, Host: AnyHost}
)

// Set sets a preference value for the given key, application ID, and preference scope.
func Set(key string, value interface{}, applicationID string, scope PreferenceScope) error {
	cKey, err := StringToCFString(key)
	if err != nil {
		return fmt.Errorf("error creating CFString for key: %v", err)
	}
	defer Release(C.CFTypeRef(cKey))

	cValue, err := ConvertToCFType(value)
	if err != nil {
		return fmt.Errorf("error converting value to CFType: %v", err)
	}
	if cValue != NilCFType {
		defer Release(cValue)
	}

	cAppID, err := StringToCFString(applicationID)
	if err != nil {
		return fmt.Errorf("error creating CFString for applicationID: %v", err)
	}
	defer Release(C.CFTypeRef(cAppID))

	var cUserName C.CFStringRef
	switch scope.User {
	case CurrentUser:
		cUserName = C.kCFPreferencesCurrentUser
	case AnyUser:
		cUserName = C.kCFPreferencesAnyUser
	default:
		return fmt.Errorf("invalid user type in scope: must be CurrentUser or AnyUser")
	}

	var cHostName C.CFStringRef
	switch scope.Host {
	case CurrentHost:
		cHostName = C.kCFPreferencesCurrentHost
	case AnyHost:
		cHostName = C.kCFPreferencesAnyHost
	default:
		return fmt.Errorf("invalid host type in scope: must be CurrentHost or AnyHost")
	}

	C.CFPreferencesSetValue(cKey, cValue, cAppID, cUserName, cHostName)

	success := C.CFPreferencesSynchronize(cAppID, cUserName, cHostName)
	if success == C.false {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// SetValue sets a preference value for the given key and application ID.
// This function uses the CurrentUserAnyHost scope.
func SetApp(key string, value interface{}, appID string) error {
	cKey, err := StringToCFString(key)
	if err != nil {
		return fmt.Errorf("error creating CFString for key: %v", err)
	}
	defer Release(C.CFTypeRef(cKey))

	cValue, err := ConvertToCFType(value)
	if err != nil {
		return fmt.Errorf("error converting value to CFType: %v", err)
	}
	if cValue != NilCFType {
		defer Release(cValue)
	}

	cAppID, err := StringToCFString(appID)
	if err != nil {
		return fmt.Errorf("error creating CFString for applicationID: %v", err)
	}
	defer Release(C.CFTypeRef(cAppID))

	C.CFPreferencesSetAppValue(cKey, cValue, cAppID)

	success := C.CFPreferencesAppSynchronize(cAppID)
	if success == C.false {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// Get retrieves a preference value for the given key, application ID, and preference scope.
func Get(key string, applicationID string, scope PreferenceScope) (interface{}, error) {
	cKey, err := StringToCFString(key)
	if err != nil {
		return nil, fmt.Errorf("error creating CFString for key: %v", err)
	}
	defer Release(C.CFTypeRef(cKey))

	cAppID, err := StringToCFString(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error creating CFString for applicationID: %v", err)
	}
	defer Release(C.CFTypeRef(cAppID))

	var cUserName C.CFStringRef
	switch scope.User {
	case CurrentUser:
		cUserName = C.kCFPreferencesCurrentUser
	case AnyUser:
		cUserName = C.kCFPreferencesAnyUser
	default:
		return nil, fmt.Errorf("invalid user type in scope: must be CurrentUser or AnyUser")
	}

	var cHostName C.CFStringRef
	switch scope.Host {
	case CurrentHost:
		cHostName = C.kCFPreferencesCurrentHost
	case AnyHost:
		cHostName = C.kCFPreferencesAnyHost
	default:
		return nil, fmt.Errorf("invalid host type in scope: must be CurrentHost or AnyHost")
	}

	value := C.CFPreferencesCopyValue(cKey, cAppID, cUserName, cHostName)
	if value == NilCFType {
		return nil, nil // Preference not found
	}
	defer Release(C.CFTypeRef(value))

	return ConvertFromCFType(value)
}

// GetValue retrieves a preference value for the given key and application ID.
// This function uses the CurrentUserAnyHost scope.
func GetApp(key string, appID string) (interface{}, error) {
	cKey, err := StringToCFString(key)
	if err != nil {
		return nil, fmt.Errorf("error creating CFString for key: %v", err)
	}
	defer Release(C.CFTypeRef(cKey))

	cAppID, err := StringToCFString(appID)
	if err != nil {
		return nil, fmt.Errorf("error creating CFString for applicationID: %v", err)
	}
	defer Release(C.CFTypeRef(cAppID))

	return Get(key, appID, CurrentUserAnyHost)

	value := C.CFPreferencesCopyAppValue(cKey, cAppID)
	if value == NilCFType {
		return nil, nil // Preference not found
	}
	defer Release(C.CFTypeRef(value))

	return ConvertFromCFType(value)
}
