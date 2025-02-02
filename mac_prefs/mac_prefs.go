package mac_prefs

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"fmt"
	"strings"
)

const (
	CurrentUser = "kCFPreferencesCurrentUser"
	AnyUser     = "kCFPreferencesAnyUser"
	CurrentHost = "kCFPreferencesCurrentHost"
	AnyHost     = "kCFPreferencesAnyHost"
)

// SetValue sets a preference value for the given key, application ID, user, and host.
func SetValue(key string, value interface{}, applicationID string, userName string, hostName string) error {
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
	switch strings.ToLower(userName) {
	case strings.ToLower(CurrentUser):
		cUserName = C.kCFPreferencesCurrentUser
	case strings.ToLower(AnyUser):
		cUserName = C.kCFPreferencesAnyUser
	default:
		return fmt.Errorf("invalid userName: must be CurrentUser or AnyUser")
	}

	var cHostName C.CFStringRef
	switch strings.ToLower(hostName) {
	case strings.ToLower(CurrentHost):
		cHostName = C.kCFPreferencesCurrentHost
	case strings.ToLower(AnyHost):
		cHostName = C.kCFPreferencesAnyHost
	default:
		return fmt.Errorf("invalid hostName: must be CurrentHost or AnyHost")
	}

	C.CFPreferencesSetValue(cKey, cValue, cAppID, cUserName, cHostName)

	success := C.CFPreferencesAppSynchronize(cAppID)
	if success == C.false {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// Set sets a preference value for the given key and application ID.
func Set(key string, value interface{}, appID string) error {
	cKey, err := StringToCFString(key)
	if err != nil {
		return fmt.Errorf("error creating CFString for key: %v", err)
	}
	defer Release(C.CFTypeRef(cKey))

	cAppID, err := StringToCFString(appID)
	if err != nil {
		return fmt.Errorf("error creating CFString for appID: %v", err)
	}
	defer Release(C.CFTypeRef(cAppID))

	cValue, err := ConvertToCFType(value)
	if err != nil {
		return fmt.Errorf("error converting value to CFType: %v", err)
	}
	defer Release(cValue)

	C.CFPreferencesSetAppValue(cKey, cValue, cAppID)
	success := C.CFPreferencesAppSynchronize(cAppID)
	if success == C.false {
		return fmt.Errorf("failed to synchronize preferences")
	}

	return nil
}

// GetValue retrieves a preference value for the given key, application ID, user, and host.
func GetValue(key string, applicationID string, userName string, hostName string) (interface{}, error) {
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
	switch strings.ToLower(userName) {
	case strings.ToLower(CurrentUser):
		cUserName = C.kCFPreferencesCurrentUser
	case strings.ToLower(AnyUser):
		cUserName = C.kCFPreferencesAnyUser
	default:
		return nil, fmt.Errorf("invalid userName: must be CurrentUser or AnyUser")
	}

	var cHostName C.CFStringRef
	switch strings.ToLower(hostName) {
	case strings.ToLower(CurrentHost):
		cHostName = C.kCFPreferencesCurrentHost
	case strings.ToLower(AnyHost):
		cHostName = C.kCFPreferencesAnyHost
	default:
		return nil, fmt.Errorf("invalid hostName: must be CurrentHost or AnyHost")
	}

	value := C.CFPreferencesCopyValue(cKey, cAppID, cUserName, cHostName)
	if value == NilCFType {
		return nil, nil // Preference not found
	}
	defer Release(C.CFTypeRef(value))

	return ConvertFromCFType(value)
}

// Get retrieves a preference value for the given key and application ID.
func Get(key string, appID string) (interface{}, error) {
	cKey, err := StringToCFString(key)
	if err != nil {
		return nil, fmt.Errorf("error creating CFString for key: %v", err)
	}
	defer Release(C.CFTypeRef(cKey))

	cAppID, err := StringToCFString(appID)
	if err != nil {
		return nil, fmt.Errorf("error creating CFString for appID: %v", err)
	}
	defer Release(C.CFTypeRef(cAppID))

	value := C.CFPreferencesCopyAppValue(cKey, cAppID)
	if value == NilCFType {
		return nil, nil // Preference not found
	}
	defer Release(C.CFTypeRef(value))

	return ConvertFromCFType(value)
}