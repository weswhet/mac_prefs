//go:build darwin

package mac_prefs

import (
	"os/user"
	"reflect"
	"testing"
	"time"
)

const testAppID = "com.github.weswhet.mac_prefs.test"

func TestSet(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         interface{}
		applicationID string
		scope         PreferenceScope
		wantErr       bool
	}{
		{
			name:          "Set string value",
			key:           "TestKey",
			value:         "TestValue",
			applicationID: testAppID,
			scope:         CurrentUserCurrentHost,
			wantErr:       false,
		},
		{
			name:          "Set integer value",
			key:           "TestIntKey",
			value:         42,
			applicationID: testAppID,
			scope:         CurrentUserCurrentHost,
			wantErr:       false,
		},
		{
			name:          "Set float value",
			key:           "TestFloatKey",
			value:         3.14,
			applicationID: testAppID,
			scope:         CurrentUserCurrentHost,
			wantErr:       false,
		},
		{
			name:          "Set slice value",
			key:           "TestSliceKey",
			value:         []interface{}{"apple", "banana", "cherry"},
			applicationID: testAppID,
			scope:         CurrentUserCurrentHost,
			wantErr:       false,
		},
		{
			name:          "Set map value",
			key:           "TestMapKey",
			value:         map[string]interface{}{"name": "John", "age": 30, "city": "New York"},
			applicationID: testAppID,
			scope:         CurrentUserCurrentHost,
			wantErr:       false,
		},
		{
			name:          "Set date value",
			key:           "TestDateKey",
			value:         time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC),
			applicationID: testAppID,
			scope:         CurrentUserCurrentHost,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Set(tt.key, tt.value, tt.applicationID, tt.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify the value was set correctly
			got, err := Get(tt.key, tt.applicationID, tt.scope)
			if err != nil {
				t.Errorf("Get() error = %v", err)
			}
			t.Logf("Get() returned type: %T, value: %v", got, got)
			t.Logf("Expected type: %T, value: %v", tt.value, tt.value)
			if !reflect.DeepEqual(got, tt.value) {
				t.Errorf("Get() got = %v (%T), want %v (%T)", got, got, tt.value, tt.value)
			}
		})
	}
}

func TestSetApp(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "SetApp string value",
			key:     "TestAppKey",
			value:   "TestAppValue",
			wantErr: false,
		},
		{
			name:    "SetApp integer value",
			key:     "TestAppIntKey",
			value:   100,
			wantErr: false,
		},
		{
			name:    "SetApp float value",
			key:     "TestAppFloatKey",
			value:   2.718,
			wantErr: false,
		},
		{
			name:    "SetApp slice value",
			key:     "TestAppSliceKey",
			value:   []interface{}{"foo", "bar", "baz"},
			wantErr: false,
		},
		{
			name:    "SetApp map value",
			key:     "TestAppMapKey",
			value:   map[string]interface{}{"x": 10, "y": 20, "label": "point"},
			wantErr: false,
		},
		{
			name:    "SetApp date value",
			key:     "TestAppDateKey",
			value:   time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetApp(tt.key, tt.value, testAppID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetApp() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify the value was set correctly
			got, err := GetApp(tt.key, testAppID)
			if err != nil {
				t.Errorf("GetApp() error = %v", err)
			}
			t.Logf("GetApp() returned type: %T, value: %v", got, got)
			t.Logf("Expected type: %T, value: %v", tt.value, tt.value)
			if !reflect.DeepEqual(got, tt.value) {
				t.Errorf("GetApp() got = %v (%T), want %v (%T)", got, got, tt.value, tt.value)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Set up test data
	testKey := "TestGetKey"
	testValue := "TestGetValue"
	testScope := CurrentUserCurrentHost

	err := Set(testKey, testValue, testAppID, testScope)
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	tests := []struct {
		name          string
		key           string
		applicationID string
		scope         PreferenceScope
		want          interface{}
		wantErr       bool
	}{
		{
			name:          "Get existing value",
			key:           testKey,
			applicationID: testAppID,
			scope:         testScope,
			want:          testValue,
			wantErr:       false,
		},
		{
			name:          "Get non-existent value",
			key:           "NonExistentKey",
			applicationID: testAppID,
			scope:         testScope,
			want:          nil,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.key, tt.applicationID, tt.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApp(t *testing.T) {
	// Set up test data
	testKey := "TestGetAppKey"
	testValue := "TestGetAppValue"

	err := SetApp(testKey, testValue, testAppID)
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	tests := []struct {
		name    string
		key     string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "GetApp existing value",
			key:     testKey,
			want:    testValue,
			wantErr: false,
		},
		{
			name:    "GetApp non-existent value",
			key:     "NonExistentAppKey",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetApp(tt.key, testAppID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveUserName(t *testing.T) {
	tests := []struct {
		name          string
		user          UserType
		wantRelease   bool
		wantValue     string
		checkCFString bool
	}{
		{
			name:        "current user constant",
			user:        CurrentUser,
			wantRelease: false,
		},
		{
			name:        "any user constant",
			user:        AnyUser,
			wantRelease: false,
		},
		{
			name:          "literal username",
			user:          UserType("alice"),
			wantRelease:   true,
			wantValue:     "alice",
			checkCFString: true,
		},
		{
			name:          "empty literal username",
			user:          UserType(""),
			wantRelease:   true,
			wantValue:     "",
			checkCFString: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, releaseRef, err := resolveUserName(tt.user)
			if err != nil {
				t.Fatalf("resolveUserName() error = %v", err)
			}
			if releaseRef != tt.wantRelease {
				t.Fatalf("resolveUserName() releaseRef = %v, want %v", releaseRef, tt.wantRelease)
			}
			if tt.checkCFString {
				defer releaseCFString(got)
				if gotValue := cfStringToString(got); gotValue != tt.wantValue {
					t.Fatalf("resolveUserName() CFString = %q, want %q", gotValue, tt.wantValue)
				}
			}
		})
	}
}

func TestSetGetSupportsLiteralCurrentUser(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current() error = %v", err)
	}
	if currentUser.Username == "" {
		t.Fatal("user.Current() returned an empty username")
	}

	key := "TestLiteralCurrentUserKey"
	want := "TestLiteralCurrentUserValue"
	scope := PreferenceScope{
		User: UserType(currentUser.Username),
		Host: AnyHost,
	}

	if err := Set(key, want, testAppID, scope); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	defer func() {
		if err := Set(key, nil, testAppID, scope); err != nil {
			t.Fatalf("cleanup Set() error = %v", err)
		}
	}()

	got, err := Get(key, testAppID, scope)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != want {
		t.Fatalf("Get() got = %v, want %v", got, want)
	}
}

func TestSetDeleteValue(t *testing.T) {
	testKey := "TestSetDeleteKey"
	testValue := "TestSetDeleteValue"
	testScope := CurrentUserCurrentHost

	// Set the initial value
	err := Set(testKey, testValue, testAppID, testScope)
	if err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}

	// Verify the value was set correctly
	got, err := Get(testKey, testAppID, testScope)
	if err != nil {
		t.Fatalf("Failed to get initial value: %v", err)
	}
	if got != testValue {
		t.Errorf("Initial value not set correctly. Got %v, want %v", got, testValue)
	}

	// Delete the value by setting it to nil
	err = Set(testKey, nil, testAppID, testScope)
	if err != nil {
		t.Fatalf("Failed to delete value: %v", err)
	}

	// Verify the value has been deleted
	got, err = Get(testKey, testAppID, testScope)
	if err != nil {
		t.Fatalf("Failed to get deleted value: %v", err)
	}
	if got != nil {
		t.Errorf("Value not deleted. Got %v, want nil", got)
	}
}

func TestSetAppSupportsUnsignedIntegers(t *testing.T) {
	const key = "TestAppUintKey"

	if err := SetApp(key, uint(42), testAppID); err != nil {
		t.Fatalf("SetApp() error = %v", err)
	}

	got, err := GetApp(key, testAppID)
	if err != nil {
		t.Fatalf("GetApp() error = %v", err)
	}
	if got != 42 {
		t.Fatalf("GetApp() got = %v (%T), want 42 (int)", got, got)
	}
}

func TestSetAppRejectsOverflowingUnsignedIntegers(t *testing.T) {
	const key = "TestAppUintOverflowKey"

	if err := SetApp(key, uint64(1<<63), testAppID); err == nil {
		t.Fatal("SetApp() expected overflow error for uint64 value above MaxInt64")
	}
}

func TestSetAppSupportsEmptyCollections(t *testing.T) {
	const sliceKey = "TestAppEmptySliceKey"
	const mapKey = "TestAppEmptyMapKey"

	if err := SetApp(sliceKey, []string{}, testAppID); err != nil {
		t.Fatalf("SetApp() empty slice error = %v", err)
	}

	gotSlice, err := GetApp(sliceKey, testAppID)
	if err != nil {
		t.Fatalf("GetApp() empty slice error = %v", err)
	}
	if !reflect.DeepEqual(gotSlice, []interface{}{}) {
		t.Fatalf("GetApp() empty slice got = %#v (%T), want []interface{}{}", gotSlice, gotSlice)
	}

	if err := SetApp(mapKey, map[string]interface{}{}, testAppID); err != nil {
		t.Fatalf("SetApp() empty map error = %v", err)
	}

	gotMap, err := GetApp(mapKey, testAppID)
	if err != nil {
		t.Fatalf("GetApp() empty map error = %v", err)
	}
	if !reflect.DeepEqual(gotMap, map[string]interface{}{}) {
		t.Fatalf("GetApp() empty map got = %#v (%T), want map[string]interface{}{}", gotMap, gotMap)
	}
}

func TestSetAppSupportsTypedStringMaps(t *testing.T) {
	const key = "TestAppTypedMapKey"

	if err := SetApp(key, map[string]string{"name": "mac_prefs"}, testAppID); err != nil {
		t.Fatalf("SetApp() typed map error = %v", err)
	}

	got, err := GetApp(key, testAppID)
	if err != nil {
		t.Fatalf("GetApp() typed map error = %v", err)
	}

	want := map[string]interface{}{"name": "mac_prefs"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetApp() typed map got = %#v (%T), want %#v (%T)", got, got, want, want)
	}
}

func TestSetAppPreservesTimePrecision(t *testing.T) {
	const key = "TestAppTimePrecisionKey"

	want := time.Date(2023, 5, 1, 12, 0, 0, 123456789, time.UTC)
	if err := SetApp(key, want, testAppID); err != nil {
		t.Fatalf("SetApp() time precision error = %v", err)
	}

	got, err := GetApp(key, testAppID)
	if err != nil {
		t.Fatalf("GetApp() time precision error = %v", err)
	}

	gotTime, ok := got.(time.Time)
	if !ok {
		t.Fatalf("GetApp() time precision got = %v (%T), want time.Time", got, got)
	}

	diff := gotTime.Sub(want)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Microsecond {
		t.Fatalf("GetApp() time precision got = %s, want %s within %s", gotTime.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano), time.Microsecond)
	}
}

func TestIsForcedApp(t *testing.T) {
	const existingKey = "TestAppForcedKey"
	const missingKey = "TestAppForcedMissingKey"

	if err := SetApp(existingKey, "not-managed", testAppID); err != nil {
		t.Fatalf("SetApp() setup error = %v", err)
	}

	for _, tc := range []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "existing unmanaged preference",
			key:  existingKey,
			want: false,
		},
		{
			name: "missing preference",
			key:  missingKey,
			want: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IsForcedApp(tc.key, testAppID)
			if err != nil {
				t.Fatalf("IsForcedApp() error = %v", err)
			}
			if got != tc.want {
				t.Fatalf("IsForcedApp() got = %v, want %v", got, tc.want)
			}
		})
	}
}
