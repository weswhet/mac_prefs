# mac_prefs

mac_prefs is a Go library for reading and writing macOS preferences using the CoreFoundation Preferences API. It provides a simple interface to interact with macOS system preferences and application-specific preferences.

## Features

- Read and write macOS system preferences
- Read and write application-specific preferences
- Support for various data types (string, integer, float, slice, map, date)
- Easy-to-use API

## Installation

To use mac_prefs in your Go project, run:

```bash
go get github.com/weswhet/mac_prefs
```

## Usage

Here are some examples of how to use the mac_prefs library:

```go
package main

import (
    "fmt"
    "github.com/weswhet/mac_prefs"
)

func main() {
    // Set an application-specific preference
    err = mac_prefs.SetApp("AppKey", 42, "com.example.app")
    if err != nil {
        fmt.Printf("Error setting app preference: %v\n", err)
    }

    // Get an application-specific preference
    appValue, err := mac_prefs.GetApp("AppKey", "com.example.app")
    if err != nil {
        fmt.Printf("Error getting app preference: %v\n", err)
    } else {
        fmt.Printf("App Value: %v\n", appValue)
    }

    // Delete a preference by setting it to nil
    err = mac_prefs.SetApp("AppKey", nil, "com.example.app")
    if err != nil {
        fmt.Printf("Error deleting preference: %v\n", err)
    }
}
```

## API Reference

### Functions

- `Set(key string, value interface{}, applicationID string, scope PreferenceScope) error`
- `Get(key string, applicationID string, scope PreferenceScope) (interface{}, error)`
- `SetApp(key string, value interface{}, applicationID string) error`
- `GetApp(key string, applicationID string) (interface{}, error)`

### Types

- `PreferenceScope`: An enum representing the scope of the preference (e.g., `CurrentUserCurrentHost`, `CurrentUserAnyHost`, `AnyUserCurrentHost`, `AnyUserAnyHost`)

#### The above preference scopes will be written to permanent storage at the following locations.

- `CurrentUserCurrentHost`
  - `[UserHomeDir]/Library/Preferences/ByHost/[applicationID].xxxx.plist`
- `CurrentUserAnyHost`
  - `[UserHomeDir]/Library/Preferences/[applicationID].plist`

#### Requires root privileges to set.

- `AnyUserCurrentHost`
  - `/Library/Preferences/[applicationID].plist`
- `AnyUserAnyHost`
  - `/var/root/Library/Preferences/ByHost/[applicationID].xxxx.plist`

### Notes

This pkg tries to mimic the usage as you would with the [CoreFoundation Preferences](https://developer.apple.com/documentation/corefoundation/preferences_utilities) library in swift. As per the documentation it is highly recommended to use higher level functions of `GetApp()` and `SetApp()` and only use the `Set()` and `Get()` functions if you absolutely have too.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License Version 2.0, - see the LICENSE file for details.
