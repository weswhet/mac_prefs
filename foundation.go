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
	"reflect"
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

// bytesToCFData converts a byte slice to a CFDataRef.
func bytesToCFData(b []byte) (C.CFDataRef, error) {
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

// cfDataToBytes converts CFData to bytes.
func cfDataToBytes(cfData C.CFDataRef) ([]byte, error) {
	return C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtr(cfData)), C.int(C.CFDataGetLength(cfData))), nil
}

// stringToCFString converts a Go string to a CFStringRef.
func stringToCFString(s string) (C.CFStringRef, error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	cfStr := C.CFStringCreateWithCString(C.kCFAllocatorDefault, cstr, C.kCFStringEncodingUTF8)
	if cfStr == NilCFString {
		return NilCFString, errors.New("CFStringCreateWithCString failed")
	}
	return cfStr, nil
}

// cfStringToString converts a CFStringRef to a Go string.
func cfStringToString(cfStr C.CFStringRef) string {
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

// mapToCFDictionary converts a Go map to a CFDictionaryRef.
func mapToCFDictionary(m map[C.CFTypeRef]C.CFTypeRef) (C.CFDictionaryRef, error) {
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

// cfDictionaryToMap converts a CFDictionaryRef to a Go map.
func cfDictionaryToMap(cfDict C.CFDictionaryRef) map[C.CFTypeRef]C.CFTypeRef {
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

// convertMapToCFDictionary converts a map[string]interface{} to a CFDictionaryRef.
func convertMapToCFDictionary(attr map[string]interface{}) (C.CFDictionaryRef, error) {
	m := make(map[C.CFTypeRef]C.CFTypeRef)
	for key, value := range attr {
		keyRef, err := stringToCFString(key)
		if err != nil {
			return NilCFDictionary, fmt.Errorf("error converting key to CFString: %v", err)
		}

		valueRef, err := convertToCFType(value)
		if err != nil {
			C.CFRelease(C.CFTypeRef(keyRef))
			return NilCFDictionary, fmt.Errorf("error converting value for key %s: %v", key, err)
		}

		m[C.CFTypeRef(keyRef)] = valueRef
	}

	cfDict, err := mapToCFDictionary(m)
	if err != nil {
		for k, v := range m {
			C.CFRelease(k)
			C.CFRelease(v)
		}
		return NilCFDictionary, err
	}
	return cfDict, nil
}

// release releases a CFTypeRef.
func release(ref C.CFTypeRef) {
	if ref != NilCFType {
		C.CFRelease(ref)
	}
}

// timeToCFDate converts a Go time.Time to a CFDateRef.
func timeToCFDate(t time.Time) C.CFDateRef {
	seconds := float64(t.Unix()) - 978307200 // Subtract seconds between 1970 and 2001
	return C.CFDateCreate(C.kCFAllocatorDefault, C.CFAbsoluteTime(seconds))
}

// cfDateToTime converts a CFDateRef to a Go time.Time.
func cfDateToTime(dateRef C.CFDateRef) time.Time {
	seconds := float64(C.CFDateGetAbsoluteTime(dateRef))
	return time.Unix(int64(seconds+978307200), 0).UTC() // Add seconds between 1970 and 2001
}

// convertToCFType converts a Go value to its corresponding CFTypeRef.
func convertToCFType(value interface{}) (C.CFTypeRef, error) {
	if value == nil {
		return NilCFType, nil
	}

	switch v := value.(type) {
	case string:
		cfStr, err := stringToCFString(v)
		if err != nil {
			return NilCFType, err
		}
		return C.CFTypeRef(cfStr), nil
	case []byte:
		cfData, err := bytesToCFData(v)
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
		return C.CFTypeRef(timeToCFDate(v)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		var numRef C.CFNumberRef
		switch num := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			int64Value := reflect.ValueOf(num).Int()
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberLongLongType, unsafe.Pointer(&int64Value))
		case float32, float64:
			floatValue := reflect.ValueOf(num).Float()
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberDoubleType, unsafe.Pointer(&floatValue))
		}
		return C.CFTypeRef(numRef), nil
	default:
		// Handle generic slices
		if slice, ok := value.([]interface{}); ok {
			return convertSliceToCFArray(slice)
		}
		sliceValue := reflect.ValueOf(value)
		if sliceValue.Kind() == reflect.Slice {
			return convertSliceToCFArray(sliceValue.Interface())
		}

		// Handle generic maps
		if m, ok := value.(map[string]interface{}); ok {
			cfDict, err := convertMapToCFDictionary(m)
			if err != nil {
				return NilCFType, err
			}
			return C.CFTypeRef(cfDict), nil
		}
		mapValue := reflect.ValueOf(value)
		if mapValue.Kind() == reflect.Map && mapValue.Type().Key().Kind() == reflect.String {
			cfDict, err := convertMapToCFDictionary(mapValue.Interface().(map[string]interface{}))
			if err != nil {
				return NilCFType, err
			}
			return C.CFTypeRef(cfDict), nil
		}

		return NilCFType, fmt.Errorf("unsupported type: %T", value)
	}
}

func convertSliceToCFArray(slice interface{}) (C.CFTypeRef, error) {
	sliceValue := reflect.ValueOf(slice)
	cfValues := make([]C.CFTypeRef, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		cfItem, err := convertToCFType(sliceValue.Index(i).Interface())
		if err != nil {
			return NilCFType, fmt.Errorf("error converting array item at index %d: %v", i, err)
		}
		cfValues[i] = cfItem
	}
	cfArray := C.CFArrayCreate(C.kCFAllocatorDefault, (*unsafe.Pointer)(unsafe.Pointer(&cfValues[0])), C.CFIndex(len(cfValues)), &C.kCFTypeArrayCallBacks)
	return C.CFTypeRef(cfArray), nil
}

// convertFromCFType converts a CFTypeRef to its corresponding Go value.
func convertFromCFType(cfType C.CFTypeRef) (interface{}, error) {
	typeID := C.CFGetTypeID(cfType)
	switch typeID {
	case C.CFStringGetTypeID():
		return cfStringToString(C.CFStringRef(cfType)), nil
	case C.CFDataGetTypeID():
		return cfDataToBytes(C.CFDataRef(cfType))
	case C.CFBooleanGetTypeID():
		return C.CFBooleanGetValue(C.CFBooleanRef(cfType)) != 0, nil
	case C.CFDateGetTypeID():
		return cfDateToTime(C.CFDateRef(cfType)), nil
	case C.CFNumberGetTypeID():
		var intValue int
		var floatValue float64
		numberType := C.CFNumberGetType(C.CFNumberRef(cfType))
		switch numberType {
		case C.kCFNumberSInt8Type, C.kCFNumberSInt16Type, C.kCFNumberSInt32Type, C.kCFNumberSInt64Type,
			C.kCFNumberCharType, C.kCFNumberShortType, C.kCFNumberIntType, C.kCFNumberLongType, C.kCFNumberLongLongType,
			C.kCFNumberCFIndexType, C.kCFNumberNSIntegerType:
			C.CFNumberGetValue(C.CFNumberRef(cfType), C.kCFNumberLongLongType, unsafe.Pointer(&intValue))
			return intValue, nil
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
			convertedItem, err := convertFromCFType(C.CFTypeRef(item))
			if err != nil {
				return nil, fmt.Errorf("error converting array item at index %d: %v", i, err)
			}
			result[i] = convertedItem
		}
		return result, nil
	case C.CFDictionaryGetTypeID():
		cfDict := C.CFDictionaryRef(cfType)
		count := C.CFDictionaryGetCount(cfDict)
		keys := make([]C.CFTypeRef, count)
		values := make([]C.CFTypeRef, count)
		C.CFDictionaryGetKeysAndValues(cfDict, (*unsafe.Pointer)(unsafe.Pointer(&keys[0])), (*unsafe.Pointer)(unsafe.Pointer(&values[0])))
		result := make(map[string]interface{}, count)
		for i := C.CFIndex(0); i < count; i++ {
			key := cfStringToString(C.CFStringRef(keys[i]))
			value, err := convertFromCFType(values[i])
			if err != nil {
				return nil, fmt.Errorf("error converting dictionary value for key %s: %v", key, err)
			}
			result[key] = value
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported CFTypeRef type")
	}
}
