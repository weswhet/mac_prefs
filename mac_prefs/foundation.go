//go:build darwin

package mac_prefs

/*
#cgo LDFLAGS: -framework CoreFoundation

#include <CoreFoundation/CoreFoundation.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"time"
	"unsafe"
)

const (
	NilCFData       C.CFDataRef       = 0
	NilCFString     C.CFStringRef     = 0
	NilCFDictionary C.CFDictionaryRef = 0
	NilCFArray      C.CFArrayRef      = 0
	NilCFType       C.CFTypeRef       = 0
)

// BytesToCFData converts a byte slice to a CFDataRef.
func BytesToCFData(b []byte) (C.CFDataRef, error) {
	if uint64(len(b)) > math.MaxUint32 {
		return C.CFDataRef(0), errors.New("data is too large")
	}
	var p *C.UInt8
	if len(b) > 0 {
		p = (*C.UInt8)(&b[0])
	}
	cfData := C.CFDataCreate(C.kCFAllocatorDefault, p, C.CFIndex(len(b)))
	if cfData == C.CFDataRef(0) {
		return C.CFDataRef(0), fmt.Errorf("CFDataCreate failed")
	}
	return cfData, nil
}

// CFDataToBytes converts CFData to bytes.
func CFDataToBytes(cfData C.CFDataRef) ([]byte, error) {
	return C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtr(cfData)), C.int(C.CFDataGetLength(cfData))), nil
}

// StringToCFString converts a Go string to a CFStringRef.
func StringToCFString(s string) (C.CFStringRef, error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	cfStr := C.CFStringCreateWithCString(C.kCFAllocatorDefault, cstr, C.kCFStringEncodingUTF8)
	if cfStr == NilCFString {
		return NilCFString, errors.New("CFStringCreateWithCString failed")
	}
	return cfStr, nil
}

// CFStringToString converts a CFStringRef to a Go string.
func CFStringToString(cfStr C.CFStringRef) string {
	length := C.CFStringGetLength(cfStr)
	if length == 0 {
		return ""
	}
	cfRange := C.CFRange{location: 0, length: length}
	enc := C.CFStringEncoding(C.kCFStringEncodingUTF8)
	var usedBufLen C.CFIndex
	if C.CFStringGetBytes(cfStr, cfRange, enc, 0, C.false, nil, 0, &usedBufLen) == 0 {
		return ""
	}
	buffer := make([]byte, usedBufLen)
	C.CFStringGetBytes(cfStr, cfRange, enc, 0, C.false, (*C.UInt8)(&buffer[0]), C.CFIndex(len(buffer)), &usedBufLen)
	return string(buffer)
}

// MapToCFDictionary converts a Go map to a CFDictionaryRef.
func MapToCFDictionary(m map[C.CFTypeRef]C.CFTypeRef) (C.CFDictionaryRef, error) {
	keys := make([]unsafe.Pointer, 0, len(m))
	values := make([]unsafe.Pointer, 0, len(m))
	for k, v := range m {
		keys = append(keys, unsafe.Pointer(k))
		values = append(values, unsafe.Pointer(v))
	}
	cfDict := C.CFDictionaryCreate(C.kCFAllocatorDefault, &keys[0], &values[0], C.CFIndex(len(m)), &C.kCFTypeDictionaryKeyCallBacks, &C.kCFTypeDictionaryValueCallBacks)
	if cfDict == NilCFDictionary {
		return NilCFDictionary, fmt.Errorf("CFDictionaryCreate failed")
	}
	return cfDict, nil
}

// CFDictionaryToMap converts a CFDictionaryRef to a Go map.
func CFDictionaryToMap(cfDict C.CFDictionaryRef) map[C.CFTypeRef]C.CFTypeRef {
	count := C.CFDictionaryGetCount(cfDict)
	if count == 0 {
		return nil
	}
	keys := make([]C.CFTypeRef, count)
	values := make([]C.CFTypeRef, count)
	C.CFDictionaryGetKeysAndValues(cfDict, (*unsafe.Pointer)(unsafe.Pointer(&keys[0])), (*unsafe.Pointer)(unsafe.Pointer(&values[0])))
	m := make(map[C.CFTypeRef]C.CFTypeRef, count)
	for i := C.CFIndex(0); i < count; i++ {
		m[keys[i]] = values[i]
	}
	return m
}

// ConvertMapToCFDictionary converts a map[string]interface{} to a CFDictionaryRef.
func ConvertMapToCFDictionary(attr map[string]interface{}) (C.CFDictionaryRef, error) {
	m := make(map[C.CFTypeRef]C.CFTypeRef)
	for key, value := range attr {
		keyRef, err := StringToCFString(key)
		if err != nil {
			return NilCFDictionary, fmt.Errorf("error converting key to CFString: %v", err)
		}

		valueRef, err := ConvertToCFType(value)
		if err != nil {
			C.CFRelease(C.CFTypeRef(keyRef))
			return NilCFDictionary, fmt.Errorf("error converting value for key %s: %v", key, err)
		}

		m[C.CFTypeRef(keyRef)] = valueRef
	}

	cfDict, err := MapToCFDictionary(m)
	if err != nil {
		for k, v := range m {
			C.CFRelease(k)
			C.CFRelease(v)
		}
		return NilCFDictionary, err
	}
	return cfDict, nil
}

// Release releases a CFTypeRef.
func Release(ref C.CFTypeRef) {
	if ref != NilCFType {
		C.CFRelease(ref)
	}
}

// TimeToCFDate converts a Go time.Time to a CFDateRef.
func TimeToCFDate(t time.Time) C.CFDateRef {
	seconds := float64(t.Unix()) - 978307200 // Subtract seconds between 1970 and 2001
	return C.CFDateCreate(C.kCFAllocatorDefault, C.CFAbsoluteTime(seconds))
}

// CFDateToTime converts a CFDateRef to a Go time.Time.
func CFDateToTime(dateRef C.CFDateRef) time.Time {
	seconds := float64(C.CFDateGetAbsoluteTime(dateRef))
	return time.Unix(int64(seconds+978307200), 0).UTC() // Add seconds between 1970 and 2001
}

// ConvertToCFType converts a Go value to its corresponding CFTypeRef.
func ConvertToCFType(value interface{}) (C.CFTypeRef, error) {
	switch v := value.(type) {
	case string:
		cfStr, err := StringToCFString(v)
		if err != nil {
			return NilCFType, err
		}
		return C.CFTypeRef(cfStr), nil
	case []byte:
		cfData, err := BytesToCFData(v)
		if err != nil {
			return NilCFType, err
		}
		return C.CFTypeRef(cfData), nil
	case bool:
		if v {
			return C.CFTypeRef(C.kCFBooleanTrue), nil
		}
		return C.CFTypeRef(C.kCFBooleanFalse), nil
	case time.Time:
		return C.CFTypeRef(TimeToCFDate(v)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		// Convert numeric types to CFNumber
		var int64Value int64
		var floatValue float64
		switch num := v.(type) {
		case int:
			int64Value = int64(num)
		case int8:
			int64Value = int64(num)
		case int16:
			int64Value = int64(num)
		case int32:
			int64Value = int64(num)
		case int64:
			int64Value = num
		case uint:
			int64Value = int64(num)
		case uint8:
			int64Value = int64(num)
		case uint16:
			int64Value = int64(num)
		case uint32:
			int64Value = int64(num)
		case uint64:
			int64Value = int64(num)
		case float32:
			floatValue = float64(num)
		case float64:
			floatValue = num
		}

		var numRef C.CFNumberRef
		if floatValue != 0 {
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberDoubleType, unsafe.Pointer(&floatValue))
		} else {
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberLongLongType, unsafe.Pointer(&int64Value))
		}
		return C.CFTypeRef(numRef), nil
	case []interface{}:
		// Convert slice of interface{} to CFArrayRef
		cfValues := make([]C.CFTypeRef, len(v))
		for i, item := range v {
			cfItem, err := ConvertToCFType(item)
			if err != nil {
				return NilCFType, fmt.Errorf("error converting array item at index %d: %v", i, err)
			}
			cfValues[i] = cfItem
		}
		cfArray := C.CFArrayCreate(C.kCFAllocatorDefault, (*unsafe.Pointer)(unsafe.Pointer(&cfValues[0])), C.CFIndex(len(cfValues)), &C.kCFTypeArrayCallBacks)
		return C.CFTypeRef(cfArray), nil
	case map[string]interface{}:
		// Convert map[string]interface{} to CFDictionaryRef
		cfDict, err := ConvertMapToCFDictionary(v)
		if err != nil {
			return NilCFType, fmt.Errorf("error converting map to CFDictionary: %v", err)
		}
		return C.CFTypeRef(cfDict), nil
	default:
		return NilCFType, fmt.Errorf("unsupported type: %T", value)
	}
}

// ConvertFromCFType converts a CFTypeRef to its corresponding Go value.
func ConvertFromCFType(cfType C.CFTypeRef) (interface{}, error) {
	typeID := C.CFGetTypeID(cfType)
	switch typeID {
	case C.CFStringGetTypeID():
		return CFStringToString(C.CFStringRef(cfType)), nil
	case C.CFDataGetTypeID():
		return CFDataToBytes(C.CFDataRef(cfType))
	case C.CFBooleanGetTypeID():
		return C.CFBooleanGetValue(C.CFBooleanRef(cfType)) != 0, nil
	case C.CFDateGetTypeID():
		return CFDateToTime(C.CFDateRef(cfType)), nil
	case C.CFNumberGetTypeID():
		var int64Value int64
		var floatValue float64
		numberType := C.CFNumberGetType(C.CFNumberRef(cfType))
		switch numberType {
		case C.kCFNumberSInt8Type, C.kCFNumberSInt16Type, C.kCFNumberSInt32Type, C.kCFNumberSInt64Type,
			C.kCFNumberCharType, C.kCFNumberShortType, C.kCFNumberIntType, C.kCFNumberLongType, C.kCFNumberLongLongType,
			C.kCFNumberCFIndexType, C.kCFNumberNSIntegerType:
			C.CFNumberGetValue(C.CFNumberRef(cfType), C.kCFNumberLongLongType, unsafe.Pointer(&int64Value))
			return int64Value, nil
		case C.kCFNumberFloat32Type, C.kCFNumberFloat64Type, C.kCFNumberFloatType, C.kCFNumberDoubleType:
			C.CFNumberGetValue(C.CFNumberRef(cfType), C.kCFNumberDoubleType, unsafe.Pointer(&floatValue))
			return floatValue, nil
		default:
			return nil, fmt.Errorf("unsupported CFNumber type")
		}
	case C.CFArrayGetTypeID():
		cfArray := C.CFArrayRef(cfType)
		count := C.CFArrayGetCount(cfArray)
		result := make([]interface{}, count)
		for i := C.CFIndex(0); i < count; i++ {
			item := C.CFArrayGetValueAtIndex(cfArray, i)
			convertedItem, err := ConvertFromCFType(C.CFTypeRef(item))
			if err != nil {
				return nil, fmt.Errorf("error converting array item at index %d: %v", i, err)
			}
			result[i] = convertedItem
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported CFTypeRef type")
	}
}