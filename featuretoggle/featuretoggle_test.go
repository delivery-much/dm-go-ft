package featuretoggle

import "testing"

var (
	mock    = redisDBMock{}
	mockErr = redisDBMock{throwErr: true}
)

func TestIsEnabled(t *testing.T) {
	t.Run("Should return the default value if the client is empty", func(t *testing.T) {
		defaultVal := true
		actual := IsEnabled("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns an error", func(t *testing.T) {
		client = &mockErr

		defaultVal := true
		actual := IsEnabled("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an error, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns an empty value", func(t *testing.T) {
		mock.setGetResult("")
		client = &mock

		defaultVal := true
		actual := IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an empty value, actualy returned %v",
				actual,
			)
		}

		mock.setGetResult(" ")
		actual = IsEnabled("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an empty value, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns a non boolean value", func(t *testing.T) {
		mock.setGetResult("thisIsNotABoolean")
		client = &mock

		defaultVal := true
		actual := IsEnabled("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned a non boolean value, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the found value if the client returns a valid boolean", func(t *testing.T) {
		mock.setGetResult("false")
		client = &mock

		expected := false
		actual := IsEnabled("MyKey", true)

		if actual != expected {
			t.Errorf(
				"Should have returned the found value if the client returned a valid boolean, actualy returned %v",
				actual,
			)
		}
	})
}

func TestGetString(t *testing.T) {
	t.Run("Should return the default value if the client is empty", func(t *testing.T) {
		client = nil
		defaultVal := "MyDefaultVal"
		actual := GetString("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns an error", func(t *testing.T) {
		client = &mockErr

		defaultVal := "MyDefaultVal"
		actual := GetString("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an error, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns an empty value", func(t *testing.T) {
		mock.setGetResult("")
		client = &mock

		defaultVal := "MyDefaultVal"
		actual := GetString("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an empty value, actualy returned %v",
				actual,
			)
		}

		mock.setGetResult("  ")
		actual = GetString("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an empty value, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the found value if the client returns a valid string", func(t *testing.T) {
		expected := "MyReturn"
		mock.setGetResult(expected)
		client = &mock

		defaultVal := "MyDefaultVal"
		actual := GetString("MyKey", defaultVal)

		if actual != expected {
			t.Errorf(
				"Should have returned the default value if the client returned a non boolean value, actualy returned %v",
				actual,
			)
		}
	})
}

func TestGetNumber(t *testing.T) {
	t.Run("Should return the default value if the client is empty", func(t *testing.T) {
		client = nil
		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client was empty, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns an error", func(t *testing.T) {
		client = &mockErr

		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an error, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns an empty value", func(t *testing.T) {
		mock.setGetResult("")
		client = &mock

		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an empty value, actualy returned %v",
				actual,
			)
		}

		mock.setGetResult("  ")
		actual = GetNumber("MyKey", defaultVal)
		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned an empty value, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the default value if the client returns a non number value", func(t *testing.T) {
		mock.setGetResult("thisIsNotANumber")
		client = &mock

		defaultVal := 14.78
		actual := GetNumber("MyKey", defaultVal)

		if actual != defaultVal {
			t.Errorf(
				"Should have returned the default value if the client returned a non number value, actualy returned %v",
				actual,
			)
		}
	})
	t.Run("Should return the found value if the client returns a valid number", func(t *testing.T) {
		client = &mock

		mock.setGetResult("2000.76")
		expected := 2000.76
		actual := GetNumber("MyKey", 100)
		if actual != expected {
			t.Errorf(
				"Should have returned the found value if the client returned a valid number, actualy returned %v",
				actual,
			)
		}

		mock.setGetResult("10")
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
	t.Run("Should return false if the client is empty", func(t *testing.T) {
		client = nil

		result := IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client was empty, actualy returned true",
			)
		}
	})
	t.Run("Should return false if the client returns an error", func(t *testing.T) {
		client = &mockErr

		result := IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client returns an error, actualy returned true",
			)
		}
	})
	t.Run("Should return false if the client returns an empty value", func(t *testing.T) {
		client = &mock

		mock.setGetResult("")
		result := IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client was empty, actualy returned true",
			)
		}

		mock.setGetResult("  ")
		result = IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client was empty, actualy returned true",
			)
		}
	})
	t.Run("Should return false if the client returns a non number value", func(t *testing.T) {
		client = &mock

		mock.setGetResult("thisIsNotANumber")
		result := IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client was empty, actualy returned true",
			)
		}
	})
	t.Run("Should return false if the client returns a non percentage value", func(t *testing.T) {
		client = &mock

		mock.setGetResult("101")
		result := IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client returned a non percentage value, actualy returned true",
			)
		}

		mock.setGetResult("-1")
		result = IsEnabledByPercent("MyKey")
		if result {
			t.Error(
				"Should have returned false if the client returned a non percentage value, actualy returned true",
			)
		}
	})
}
