package featuretoggle

import "testing"

func TestIsEnabled(t *testing.T) {
	t.Run("Should return the default value if the local memory is empty", func(t *testing.T) {
		defaultVal := true
		actual := IsEnabled("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key value is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey": "",
		}
		defaultVal := true

		actual := IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{}

		actual = IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key type is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "1",
			"MyKey.type": "",
		}
		defaultVal := true

		actual := IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey": "1",
		}

		actual = IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key type is not 'boolean'", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "1",
			"MyKey.type": "not boolean",
		}
		defaultVal := true

		actual := IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was not 'boolean', actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key value is not a valid boolean", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "not a boolean value",
			"MyKey.type": "boolean",
		}

		defaultVal := true
		actual := IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was not a boolean, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the found value if the key represents a valid boolean", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "0",
			"MyKey.type": "boolean",
		}

		expected := false
		actual := IsEnabled("MyKey", true)
		if actual != expected {
			t.Errorf(
				"Should have returned the found value if the key was a valid boolean, actualy returned %v",
				actual,
			)
		}
	})
}

func TestGetString(t *testing.T) {
	t.Run("Should return the default value if the local memory is empty", func(t *testing.T) {
		localMemory = nil
		defaultVal := "MyDefaultVal"
		actual := GetString("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the local memory was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key value is empty", func(t *testing.T) {
		localMemory = map[string]string{}
		defaultVal := "MyDefaultVal"

		actual := GetString("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was empty actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey": "",
		}

		actual = GetString("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key type is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey": "myval",
		}
		defaultVal := "MyDefaultVal"

		actual := GetString("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey":      "myval",
			"MyKey.type": "",
		}

		actual = GetString("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key type is not 'string'", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "myval",
			"MyKey.type": "not a string",
		}

		defaultVal := "MyDefaultVal"
		actual := GetString("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was not 'string', actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the found value if the key represents a valid string", func(t *testing.T) {
		expected := "MyReturn"
		localMemory = map[string]string{
			"MyKey":      expected,
			"MyKey.type": "string",
		}

		actual := GetString("MyKey", "MyDefaultVal")
		if actual != expected {
			t.Errorf(
				"Should have returned the default value if the value was a valid string, actualy returned %v",
				actual,
			)
		}
	})
}

func TestGetNumber(t *testing.T) {
	t.Run("Should return the default value if the local memory is empty", func(t *testing.T) {
		localMemory = nil
		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the local memory was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key value is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey": "",
		}
		defaultVal := 14.78

		actual := GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{}

		actual = GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key type is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "1",
			"MyKey.type": "",
		}
		defaultVal := 14.78

		actual := GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey": "1",
		}

		actual = GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key value is a non number value", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "not a number",
			"MyKey.type": "number",
		}

		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key value was not a number, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the key type is not 'number'", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "100.5",
			"MyKey.type": "not number",
		}

		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the key type was not 'number', actualy returned %v",
				actual,
			)
		}

	})
	t.Run("Should return the found value if the client returns a valid number", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "2000.76",
			"MyKey.type": "number",
		}

		expected := 2000.76
		actual := GetNumber("MyKey", 100)
		if actual != expected {
			t.Errorf(
				"Should have returned the found value if the client returned a valid number, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey":      "10",
			"MyKey.type": "number",
		}
		expected = 10.0
		actual = GetNumber("MyKey", 100)
		if actual != expected {
			t.Errorf(
				"Should have returned the found value if the client returned a valid number, actualy returned %v",
				actual,
			)
		}
	})
}

func TestIsEnabledByPercent(t *testing.T) {
	t.Run("Should return false if the local memory is empty", func(t *testing.T) {
		localMemory = nil

		actual := IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the local memory was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return false if the key value is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey": "",
		}
		actual := IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key value was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{}
		actual = IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key value was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return false if the key type is empty", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "1",
			"MyKey.type": "",
		}
		actual := IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key type was empty, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey": "1",
		}
		actual = IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key type was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return false if the key value is a non number value", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "not a number",
			"MyKey.type": "number",
		}
		actual := IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key value was not a number, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return false if the key type is not 'number'", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "100.5",
			"MyKey.type": "not number",
		}
		actual := IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key value was not a number, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return false if the key value is a non percentage value", func(t *testing.T) {
		localMemory = map[string]string{
			"MyKey":      "101",
			"MyKey.type": "number",
		}
		actual := IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key value was a non percentage value, actualy returned %v",
				actual,
			)
		}

		localMemory = map[string]string{
			"MyKey":      "-1",
			"MyKey.type": "number",
		}
		actual = IsEnabledByPercent("MyKey")
		if actual {
			t.Errorf(
				"Should have returned false if the key value was a non percentage value, actualy returned %v",
				actual,
			)
		}
	})
}
