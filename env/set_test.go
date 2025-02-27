package env_test

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/kyuff/anchor/env"
	"github.com/kyuff/anchor/internal/assert"
)

func TestSet(t *testing.T) {
	var (
		newKey = func() string {
			return fmt.Sprintf("KEY_%d", rand.IntN(100000))
		}
		newValue = func() string {
			return fmt.Sprintf("VALUE_%d", rand.IntN(100000))
		}
	)

	t.Run("override key value", func(t *testing.T) {
		// arrange
		var (
			key   = newKey()
			value = newValue()
		)

		assert.NoError(t, os.Setenv(key, newValue()))

		// act
		err := env.Set(
			env.OverrideKeyValue(key, value),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, os.Getenv(key), value)
	})

	t.Run("override key value on t", func(t *testing.T) {
		// arrange
		var (
			key   = newKey()
			value = newValue()
		)

		// act
		err := env.Set(
			env.OverrideKeyValue(key, value),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("default key value set", func(t *testing.T) {
		// arrange
		var (
			key   = newKey()
			value = newValue()
		)

		assert.NoError(t, os.Unsetenv(key))

		// act
		err := env.Set(
			env.DefaultKeyValue(key, value),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("default key value skip", func(t *testing.T) {
		// arrange
		var (
			key   = newKey()
			value = newValue()
		)

		assert.NoError(t, os.Setenv(key, value))

		// act
		err := env.Set(
			env.DefaultKeyValue(key, newValue()),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("override from file", func(t *testing.T) {
		// arrange
		var (
			key   = "ANOTHER_TEST_KEY"
			value = "true"
		)

		// act
		err := env.Set(
			env.OverrideFile("testdata/data.env"),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("override from file with whitespace key", func(t *testing.T) {
		// arrange
		var (
			key   = "TEST_KEY_412"
			value = "1251a"
		)

		// act
		err := env.Set(
			env.OverrideFile("testdata/data.env"),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("override from file fail on missing file", func(t *testing.T) {
		// act
		err := env.Set(
			env.OverrideFile("testdata/not_there.env"),
			env.WithT(t),
		)

		// assert
		assert.Error(t, err)
	})

	t.Run("override from file fail on malformed file", func(t *testing.T) {
		// act
		err := env.Set(
			env.OverrideFile("testdata/malformed.env"),
			env.WithT(t),
		)

		// assert
		assert.Error(t, err)
	})

	t.Run("default from file set value", func(t *testing.T) {
		// arrange
		var (
			key   = "ANOTHER_TEST_KEY"
			value = "true"
		)

		assert.NoError(t, os.Unsetenv(key))

		// act
		err := env.Set(
			env.DefaultFile("testdata/data.env"),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("default from file skip value", func(t *testing.T) {
		// arrange
		var (
			key   = "ANOTHER_TEST_KEY"
			value = newValue()
		)

		assert.NoError(t, os.Setenv(key, value))

		// act
		err := env.Set(
			env.DefaultFile("testdata/data.env"),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("override from env key file", func(t *testing.T) {
		// arrange
		var (
			key      = "ANOTHER_TEST_KEY"
			value    = "true"
			envKey   = newKey()
			envValue = "testdata/data.env"
		)

		assert.NoError(t, os.Setenv(envKey, envValue))

		// act
		err := env.Set(
			env.OverrideEnvKeyFile(envKey),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("default from env key file set value", func(t *testing.T) {
		// arrange
		var (
			key      = "ANOTHER_TEST_KEY"
			value    = "true"
			envKey   = newKey()
			envValue = "testdata/data.env"
		)

		assert.NoError(t, os.Setenv(envKey, envValue))
		assert.NoError(t, os.Unsetenv(key))

		// act
		err := env.Set(
			env.DefaultEnvKeyFile(envKey),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("default from env key file skip value", func(t *testing.T) {
		// arrange
		var (
			key      = "ANOTHER_TEST_KEY"
			value    = newValue()
			envKey   = newKey()
			envValue = "testdata/data.env"
		)

		assert.NoError(t, os.Setenv(envKey, envValue))
		assert.NoError(t, os.Setenv(key, value))

		// act
		err := env.Set(
			env.DefaultEnvKeyFile(envKey),
			env.WithT(t),
		)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, value, os.Getenv(key))
	})

	t.Run("default from env key fail on missing file", func(t *testing.T) {
		// arrange
		var (
			envKey   = newKey()
			envValue = "testdata/not_there.env"
		)

		assert.NoError(t, os.Setenv(envKey, envValue))

		// act
		err := env.Set(
			env.DefaultEnvKeyFile(envKey),
			env.WithT(t),
		)

		// assert
		assert.Error(t, err)
	})

	t.Run("panic on wrong configuration", func(t *testing.T) {
		assert.Panic(t, func() {
			env.MustSet(
				env.OverrideFile("testdata/not_there.env"),
				env.WithT(t),
			)
		})
	})

	t.Run("fail on setting environment", func(t *testing.T) {
		// act
		err := env.Set(
			env.OverrideFile("testdata/data.env"),
			env.WithEnvironment(func(key, value string) error {
				return errors.New("test")
			}),
		)

		// assert
		assert.Error(t, err)
	})
}
