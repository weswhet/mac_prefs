package mac_prefs

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"fmt"
)

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