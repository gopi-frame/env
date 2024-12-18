package env

import (
	"encoding/json"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"io"
	"os"
	"strconv"
	"strings"
)

// Load loads the named file(s) into the environment.
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}
	current := os.Environ()
	currentKeys := make(map[string]bool)
	for _, v := range current {
		kv := strings.SplitN(v, "=", 2)
		currentKeys[kv[0]] = true
	}
	for _, filename := range filenames {
		envMap, err := godotenv.Read(filename)
		if err != nil {
			return err
		}
		for k, v := range envMap {
			if _, ok := currentKeys[k]; !ok {
				v := Expand(v)
				if err := Set(k, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Override loads the named file(s) into the environment, overriding any existing values.
func Override(filenames ...string) error {
	for _, filename := range filenames {
		envMap, err := godotenv.Read(filename)
		if err != nil {
			return err
		}
		for k, v := range envMap {
			v := os.Expand(v, Expand)
			if err := Set(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func FromReader(r io.Reader) (map[string]string, error) {
	return godotenv.Parse(r)
}

func Unmarshal[T any]() (T, error) {
	return env.ParseAs[T]()
}

func MustUnmarshal[T any]() T {
	return env.Must(env.ParseAs[T]())
}

// Set sets the value of the environment variable named by the key.
func Set(key string, value string) error {
	return os.Setenv(key, value)
}

// SetInt is a shorthand for Set(key, strconv.Itoa(value)).
func SetInt(key string, value int) error {
	return os.Setenv(key, strconv.Itoa(value))
}

// SetInt64 is a shorthand for Set(key, strconv.FormatInt(value, 10)).
func SetInt64(key string, value int64) error {
	return os.Setenv(key, strconv.FormatInt(value, 10))
}

// SetUint64 is a shorthand for Set(key, strconv.FormatUint(value, 10)).
func SetUint64(key string, value uint64) error {
	return os.Setenv(key, strconv.FormatUint(value, 10))
}

// SetFloat64 is a shorthand for Set(key, strconv.FormatFloat(value, 'f', -1, 64)).
func SetFloat64(key string, value float64) error {
	return os.Setenv(key, strconv.FormatFloat(value, 'f', -1, 64))
}

// SetBool is a shorthand for Set(key, strconv.FormatBool(value)).
func SetBool(key string, value bool) error {
	return os.Setenv(key, strconv.FormatBool(value))
}

// SetStrings is a shorthand for Set(key, strings.Join(value, ",")).
func SetStrings(key string, value []string) error {
	return os.Setenv(key, strings.Join(value, ","))
}

// SetJSON is a shorthand for Set(key, json.Marshal(value)).
func SetJSON(key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return os.Setenv(key, string(b))
}

// Get returns the value of the environment variable named by the key.
func Get(key string) string {
	return os.Getenv(key)
}

// GetOr returns the value of the environment variable named by the key.
// If the variable is not present, it returns the default value.
func GetOr(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// GetT returns the value of the environment variable named by the key.
// It converts the value to the specified type using the provided conversion function.
func GetT[T any](key string, convert func(s string) (T, error)) (T, error) {
	return convert(os.Getenv(key))
}

// MustGetT returns the value of the environment variable named by the key.
// It converts the value to the specified type using the provided conversion function.
// If the value cannot be converted, it panics.
func MustGetT[T any](key string, convert func(s string) (T, error)) T {
	value, err := GetT[T](key, convert)
	if err == nil {
		return value
	}
	panic(err)
}

// GetTOr returns the value of the environment variable named by the key.
// It converts the value to the specified type using the provided conversion function.
// If the variable is not present, it returns the default value.
// If the value cannot be converted, it returns the default value and an error.
func GetTOr[T any](key string, convert func(s string) (T, error), defaultValue T) (T, error) {
	if value, ok := os.LookupEnv(key); ok {
		if value, err := convert(value); err != nil {
			return defaultValue, err
		} else {
			return value, nil
		}
	} else {
		return defaultValue, nil
	}
}

// GetInt is a shorthand for GetT[int](key, strconv.Atoi)
func GetInt(key string) (int, error) {
	return GetT(key, strconv.Atoi)
}

// MustGetInt is a shorthand for MustGetT[int](key, strconv.Atoi)
func MustGetInt(key string) int {
	return MustGetT(key, strconv.Atoi)
}

// GetIntOr is a shorthand for GetTOr[int](key, strconv.Atoi, defaultValue)
func GetIntOr(key string, defaultValue int) (int, error) {
	return GetTOr(key, strconv.Atoi, defaultValue)
}

// GetInt64 is a shorthand for GetT[int64](key, strconv.ParseInt)
func GetInt64(key string) (int64, error) {
	return GetT(key, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	})
}

// MustGetInt64 is a shorthand for MustGetT[int64](key, strconv.ParseInt)
func MustGetInt64(key string) int64 {
	return MustGetT(key, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	})
}

// GetInt64Or is a shorthand for GetTOr[int64](key, strconv.ParseInt, defaultValue)
func GetInt64Or(key string, defaultValue int64) (int64, error) {
	return GetTOr(key, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}, defaultValue)
}

// GetUint64 is a shorthand for GetT[uint64](key, strconv.ParseUint)
func GetUint64(key string) (uint64, error) {
	return GetT(key, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	})
}

// MustGetUint64 is a shorthand for MustGetT[uint64](key, strconv.ParseUint)
func MustGetUint64(key string) uint64 {
	return MustGetT(key, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	})
}

// GetUint64Or is a shorthand for GetTOr[uint64](key, strconv.ParseUint, defaultValue)
func GetUint64Or(key string, defaultValue uint64) (uint64, error) {
	return GetTOr(key, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	}, defaultValue)
}

// GetFloat64 is a shorthand for GetT[float64](key, strconv.ParseFloat)
func GetFloat64(key string) (float64, error) {
	return GetT(key, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

// MustGetFloat64 is a shorthand for MustGetT[float64](key, strconv.ParseFloat)
func MustGetFloat64(key string) float64 {
	return MustGetT(key, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

// GetFloat64Or is a shorthand for GetTOr[float64](key, strconv.ParseFloat, defaultValue)
func GetFloat64Or(key string, defaultValue float64) (float64, error) {
	return GetTOr(key, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	}, defaultValue)
}

// GetBool is a shorthand for GetT[bool](key, strconv.ParseBool)
func GetBool(key string) (bool, error) {
	return GetT(key, func(s string) (bool, error) {
		return strconv.ParseBool(s)
	})
}

// MustGetBool is a shorthand for MustGetT[bool](key, strconv.ParseBool)
func MustGetBool(key string) bool {
	return MustGetT(key, func(s string) (bool, error) {
		return strconv.ParseBool(s)
	})
}

// GetBoolOr is a shorthand for GetTOr[bool](key, strconv.ParseBool, defaultValue)
func GetBoolOr(key string, defaultValue bool) (bool, error) {
	return GetTOr(key, func(s string) (bool, error) {
		return strconv.ParseBool(s)
	}, defaultValue)
}

// GetStrings is a shorthand for GetT[[]string]
func GetStrings(key string) ([]string, error) {
	return GetT(key, func(s string) ([]string, error) {
		return strings.Split(s, ","), nil
	})
}

func MustGetStrings(key string) []string {
	return MustGetT(key, func(s string) ([]string, error) {
		return strings.Split(s, ","), nil
	})
}

// GetStringsOr is a shorthand for GetTOr[[]string](key, func(s string) ([]string, error), defaultValue)
func GetStringsOr(key string, defaultValue []string) ([]string, error) {
	return GetTOr(key, func(s string) ([]string, error) {
		if s == "" {
			if _, ok := os.LookupEnv(key); !ok {
				return nil, errors.New("env not found")
			}
		}
		return strings.Split(s, ","), nil
	}, defaultValue)
}

// GetJSON is a shorthand for GetT[T](key, json.Unmarshal)
func GetJSON[T any](key string) (T, error) {
	return GetT[T](key, func(s string) (T, error) {
		var value T
		if err := json.Unmarshal([]byte(s), &value); err != nil {
			return value, err
		}
		return value, nil
	})
}

func MustGetJSON[T any](key string) T {
	return MustGetT(key, func(s string) (T, error) {
		var value T
		if err := json.Unmarshal([]byte(s), &value); err != nil {
			return value, err
		}
		return value, nil
	})
}

// GetJSONOr is a shorthand for GetTOr[T](key, json.Unmarshal, defaultValue)
func GetJSONOr[T any](key string, defaultValue T) (T, error) {
	return GetTOr(key, func(s string) (T, error) {
		var value T
		if err := json.Unmarshal([]byte(s), &value); err != nil {
			return value, err
		}
		return value, nil
	}, defaultValue)
}

// Expand expands the environment variables in the given string.
// The string can contain environment variables in the form of the following syntax:
//   - ${ENV_KEY}: expands to the value of the environment variable ENV_KEY
//   - ${ENV_KEY|default}: expands to the value of the environment variable ENV_KEY, or the default value if the environment variable is not set
//   - ${ENV_KEY|FALLBACK_ENV_KEY_1|...|default}: expands to the value of the environment variable ENV_KEY, or the first fallback environment variable which is set, or the default value if all the fallback environment variables are not set
//
// Some spacial cases:
//   - ${ENV_KEY|}: expands to the value of the environment variable ENV_KEY, or the empty string if the environment variable is not set
//   - ${|default}: equals to default value
//   - ${|}: equals to the empty string
//   - ${ENV_KEY|FALLBACK_ENV_KEY|...|}: expands to the value of the environment variable ENV_KEY, or the first fallback environment variable which is set, or the empty string if all the fallback environment variables are not set
//   - ${ENV_KEY\\|_USERNAME}: env key with `|` character
//   - ${ENV_KEY|default_\\|value}: default value with `|` character
//
// Example:
//
//	os.Setenv("ENV_KEY", "value")
//	os.Setenv("FALLBACK_ENV_KEY_1", "fallback_value_1")
//	value := Expand("${ENV_KEY}") // value = "value"
//	value2 := Expand("${ENV_KEY|default}") // value2 = "value"
//	value3 := Expand("${NOT_FOUND_ENV_KEY|default_value}") // value3 = "default_value"
//	value4 := Expand("${NOT_FOUND_ENV_KEY|FALLBACK_ENV_KEY_1|}") // value4 = "fallback_value_1"
//	value5 := Expand("${NOT_FOUND_ENV_KEY|NOT_FOUND_FALLBACK_ENV_KEY_1|default_value}") // value5 = "fallback_value_1"
func Expand(key string) string {
	return os.Expand(key, func(key string) string {
		var envKeys []string
		var defaultValue string
		var startIndex int
		var withDefault bool
		key = strings.TrimSpace(key)
		if key[0] == '|' {
			withDefault = true
			key = key[1:]
		}
		for i, b := range key {
			if b == '|' {
				if key[i-1] != '\\' {
					envKeys = append(envKeys, key[startIndex:i])
					startIndex = i + 1
					withDefault = true
				}
			}
		}
		if !withDefault {
			envKeys = append(envKeys, key)
		} else {
			defaultValue = strings.ReplaceAll(key[startIndex:], "\\|", "|")
		}
		for _, envKey := range envKeys {
			if lookupEnv, ok := os.LookupEnv(envKey); ok {
				return lookupEnv
			}
		}
		return defaultValue
	})
}
