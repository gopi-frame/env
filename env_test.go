package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExpand(t *testing.T) {
	t.Run("without variable", func(t *testing.T) {
		value := Expand("test_value")
		assert.Equal(t, "test_value", value)
	})

	t.Run("without default value", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "test_value")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value := Expand("${TEST_ENV_KEY}")
		assert.Equal(t, "test_value", value)
	})

	t.Run("with default value", func(t *testing.T) {
		value := Expand("${TEST_ENV_KEY|default_value}")
		assert.Equal(t, "default_value", value)
	})

	t.Run("with default value and escaped", func(t *testing.T) {
		value := Expand("${TEST_ENV_KEY|default_value\\|escaped}")
		assert.Equal(t, "default_value|escaped", value)
	})

	t.Run("with fallback env and without default value", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "test_value")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value := Expand("${NOT_EXIST_ENV|TEST_ENV_KEY|}")
		assert.Equal(t, "test_value", value)
	})

	t.Run("with fallback env and with default value", func(t *testing.T) {
		value := Expand("${NOT_EXIST_ENV|TEST_ENV_KEY|default_value}")
		assert.Equal(t, "default_value", value)
	})

	t.Run("start with vertical", func(t *testing.T) {
		value := Expand("${|default_value}")
		assert.Equal(t, "default_value", value)
	})
}

func TestSet(t *testing.T) {
	err := Set("TEST_ENV_KEY", "test_value")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	value := Get("TEST_ENV_KEY")
	assert.Equal(t, "test_value", value)
}

func TestSetInt(t *testing.T) {
	err := SetInt("TEST_ENV_KEY", 123)
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "123", Get("TEST_ENV_KEY"))
	assert.Equal(t, 123, MustGetInt("TEST_ENV_KEY"))
}

func TestSetInt64(t *testing.T) {
	err := SetInt64("TEST_ENV_KEY", 123)
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "123", Get("TEST_ENV_KEY"))
	assert.Equal(t, int64(123), MustGetInt64("TEST_ENV_KEY"))
}

func TestSetUint64(t *testing.T) {
	err := SetUint64("TEST_ENV_KEY", 123)
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "123", Get("TEST_ENV_KEY"))
	assert.Equal(t, uint64(123), MustGetUint64("TEST_ENV_KEY"))
}

func TestSetFloat64(t *testing.T) {
	err := SetFloat64("TEST_ENV_KEY", 123.456)
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "123.456", Get("TEST_ENV_KEY"))
	assert.Equal(t, 123.456, MustGetFloat64("TEST_ENV_KEY"))
}

func TestSetBool(t *testing.T) {
	err := SetBool("TEST_ENV_KEY", true)
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "true", Get("TEST_ENV_KEY"))
	assert.Equal(t, true, MustGetBool("TEST_ENV_KEY"))
}

func TestSetStrings(t *testing.T) {
	err := SetStrings("TEST_ENV_KEY", []string{"a", "b", "c"})
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "a,b,c", Get("TEST_ENV_KEY"))
	assert.Equal(t, []string{"a", "b", "c"}, MustGetStrings("TEST_ENV_KEY"))
}

func TestSetJSON(t *testing.T) {
	err := SetJSON("TEST_ENV_KEY", map[string]string{"a": "b"})
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_KEY")
	}()
	assert.Equal(t, "{\"a\":\"b\"}", Get("TEST_ENV_KEY"))
	assert.Equal(t, map[string]string{"a": "b"}, MustGetJSON[map[string]string]("TEST_ENV_KEY"))
}

func TestGetT(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetInt("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, 123, value)
	})
	t.Run("int64", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetInt64("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, int64(123), value)
	})
	t.Run("uint64", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetUint64("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), value)
	})
	t.Run("float64", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123.456")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetFloat64("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, 123.456, value)
	})
	t.Run("bool", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "true")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetBool("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, true, value)
	})
	t.Run("string slice", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "a,b,c")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetStrings("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, value)
	})
	t.Run("json", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "{\"a\":\"b\"}")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
		}()
		value, err := GetJSON[map[string]string]("TEST_ENV_KEY")
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"a": "b"}, value)
	})
}

func TestGetOr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "test_value")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value := GetOr("TEST_ENV_KEY", "default_value")
			assert.Equal(t, "default_value", value)
		}()
		value := GetOr("TEST_ENV_KEY", "default_value")
		assert.Equal(t, "test_value", value)
	})
	t.Run("int", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetIntOr("TEST_ENV_KEY", 456)
			assert.NoError(t, err)
			assert.Equal(t, 456, value)
		}()
		value, err := GetIntOr("TEST_ENV_KEY", 456)
		assert.NoError(t, err)
		assert.Equal(t, 123, value)
	})
	t.Run("int64", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetInt64Or("TEST_ENV_KEY", 456)
			assert.NoError(t, err)
			assert.Equal(t, int64(456), value)
		}()
		value, err := GetInt64Or("TEST_ENV_KEY", 456)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), value)
	})
	t.Run("uint64", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetUint64Or("TEST_ENV_KEY", 456)
			assert.NoError(t, err)
			assert.Equal(t, uint64(456), value)
		}()
		value, err := GetUint64Or("TEST_ENV_KEY", 456)
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), value)
	})
	t.Run("float64", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "123.456")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetFloat64Or("TEST_ENV_KEY", 456.789)
			assert.NoError(t, err)
			assert.Equal(t, 456.789, value)
		}()
		value, err := GetFloat64Or("TEST_ENV_KEY", 456.789)
		assert.NoError(t, err)
		assert.Equal(t, 123.456, value)
	})
	t.Run("bool", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "true")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetBoolOr("TEST_ENV_KEY", false)
			assert.NoError(t, err)
			assert.Equal(t, false, value)
		}()
		value, err := GetBoolOr("TEST_ENV_KEY", false)
		assert.NoError(t, err)
		assert.Equal(t, true, value)
	})
	t.Run("string slice", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "a,b,c")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetStringsOr("TEST_ENV_KEY", []string{"d", "e", "f"})
			assert.NoError(t, err)
			assert.Equal(t, []string{"d", "e", "f"}, value)
		}()
		value, err := GetStringsOr("TEST_ENV_KEY", []string{"d", "e", "f"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, value)
	})
	t.Run("json", func(t *testing.T) {
		_ = os.Setenv("TEST_ENV_KEY", "{\"a\":\"b\"}")
		defer func() {
			_ = os.Unsetenv("TEST_ENV_KEY")
			value, err := GetJSONOr[map[string]string]("TEST_ENV_KEY", map[string]string{"c": "d"})
			assert.NoError(t, err)
			assert.Equal(t, map[string]string{"c": "d"}, value)
		}()
		value, err := GetJSONOr[map[string]string]("TEST_ENV_KEY", map[string]string{"c": "d"})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"a": "b"}, value)
	})
}

func TestLoad(t *testing.T) {
	err := Load("testdata/.env")
	if !assert.NoError(t, err) {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, "localhost", Get("DB_HOST"))
	assert.Equal(t, 3306, MustGetInt("DB_PORT"))
	assert.Equal(t, "root", Get("DB_USER"))
	assert.Equal(t, "root", Get("DB_PASSWORD"))
	assert.Equal(t, "gopi", Get("DB_DATABASE"))
	assert.Equal(t, "root:root@tcp(localhost:3306)/gopi", Get("DB_DSN"))
}
