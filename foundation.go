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

var cfAbsoluteTimeEpoch = time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)

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
	keys := make([]C.CFTypeRef, 0, len(m))
	values := make([]C.CFTypeRef, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}

	var keyPtr *unsafe.Pointer
	var valuePtr *unsafe.Pointer
	if len(keys) > 0 {
		keyPtr = (*unsafe.Pointer)(unsafe.Pointer(&keys[0]))
		valuePtr = (*unsafe.Pointer)(unsafe.Pointer(&values[0]))
	}

	cfDict := C.CFDictionaryCreate(C.kCFAllocatorDefault, keyPtr, valuePtr, C.CFIndex(len(m)), &C.kCFTypeDictionaryKeyCallBacks, &C.kCFTypeDictionaryValueCallBacks)
	if cfDict == NilCFDictionary {
		return NilCFDictionary, fmt.Errorf("CFDictionaryCreate failed")
	}
	return cfDict, nil
}

// cfDictionaryToMap converts a CFDictionaryRef to a Go map.
func cfDictionaryToMap(cfDict C.CFDictionaryRef) map[C.CFTypeRef]C.CFTypeRef {
	count := C.CFDictionaryGetCount(cfDict)
	if count == 0 {
		return map[C.CFTypeRef]C.CFTypeRef{}
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

// convertMapToCFDictionary converts a string-keyed Go map to a CFDictionaryRef.
func convertMapToCFDictionary(attr interface{}) (C.CFDictionaryRef, error) {
	mapValue := reflect.ValueOf(attr)
	if mapValue.Kind() != reflect.Map || mapValue.Type().Key().Kind() != reflect.String {
		return NilCFDictionary, fmt.Errorf("unsupported map type: %T", attr)
	}

	m := make(map[C.CFTypeRef]C.CFTypeRef, mapValue.Len())
	releaseMapEntries := func() {
		for k, v := range m {
			release(k)
			release(v)
		}
	}
	for _, keyValue := range mapValue.MapKeys() {
		key := keyValue.String()

		keyRef, err := stringToCFString(key)
		if err != nil {
			releaseMapEntries()
			return NilCFDictionary, fmt.Errorf("error converting key to CFString: %v", err)
		}

		valueRef, err := convertToCFType(mapValue.MapIndex(keyValue).Interface())
		if err != nil {
			release(C.CFTypeRef(keyRef))
			releaseMapEntries()
			return NilCFDictionary, fmt.Errorf("error converting value for key %s: %v", key, err)
		}

		m[C.CFTypeRef(keyRef)] = valueRef
	}

	cfDict, err := mapToCFDictionary(m)
	if err != nil {
		releaseMapEntries()
		return NilCFDictionary, err
	}

	releaseMapEntries()
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
	seconds := t.UTC().Sub(cfAbsoluteTimeEpoch).Seconds()
	return C.CFDateCreate(C.kCFAllocatorDefault, C.CFAbsoluteTime(seconds))
}

// cfDateToTime converts a CFDateRef to a Go time.Time.
func cfDateToTime(dateRef C.CFDateRef) time.Time {
	seconds := float64(C.CFDateGetAbsoluteTime(dateRef))
	nanos := int64(math.Round(seconds * float64(time.Second)))
	return cfAbsoluteTimeEpoch.Add(time.Duration(nanos))
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
		numberValue := reflect.ValueOf(v)
		switch numberValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			int64Value := numberValue.Int()
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberLongLongType, unsafe.Pointer(&int64Value))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uint64Value := numberValue.Uint()
			if uint64Value > math.MaxInt64 {
				return NilCFType, fmt.Errorf("unsigned integer %d overflows signed 64-bit CFNumber", uint64Value)
			}
			int64Value := int64(uint64Value)
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberLongLongType, unsafe.Pointer(&int64Value))
		case reflect.Float32, reflect.Float64:
			floatValue := numberValue.Float()
			numRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberDoubleType, unsafe.Pointer(&floatValue))
		}
		if numRef == 0 {
			return NilCFType, fmt.Errorf("CFNumberCreate failed")
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
			cfDict, err := convertMapToCFDictionary(value)
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
			for _, value := range cfValues[:i] {
				release(value)
			}
			return NilCFType, fmt.Errorf("error converting array item at index %d: %v", i, err)
		}
		cfValues[i] = cfItem
	}

	var valuePtr *unsafe.Pointer
	if len(cfValues) > 0 {
		valuePtr = (*unsafe.Pointer)(unsafe.Pointer(&cfValues[0]))
	}

	cfArray := C.CFArrayCreate(C.kCFAllocatorDefault, valuePtr, C.CFIndex(len(cfValues)), &C.kCFTypeArrayCallBacks)
	if cfArray == NilCFArray {
		for _, value := range cfValues {
			release(value)
		}
		return NilCFType, fmt.Errorf("CFArrayCreate failed")
	}

	for _, value := range cfValues {
		release(value)
	}
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
		if count == 0 {
			return map[string]interface{}{}, nil
		}
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
