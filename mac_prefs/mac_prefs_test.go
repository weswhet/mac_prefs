package mac_prefs

import (
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
