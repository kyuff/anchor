package assert

import (
	"fmt"
	"testing"
	"time"
)

func Equal[T comparable](t *testing.T, expected, got T) bool {
	t.Helper()
	return Equalf(t, expected, got, "Items was not equal")
}

func Equalf[T comparable](t *testing.T, expected, got T, format string, args ...any) bool {
	t.Helper()
	if expected != got {
		t.Logf(`
%s
Expected: %v
     Got: %v`, fmt.Sprintf(format, args...), expected, got)
		t.Fail()
		return false
	}
	return true
}

func Truef(t *testing.T, got bool, format string, args ...any) bool {
	t.Helper()
	if !got {
		t.Logf(format, args...)
		t.Fail()
		return false
	}

	return true
}

func EqualSlice[T comparable](t *testing.T, expected, got []T) bool {
	t.Helper()
	if len(expected) != len(got) {
		t.Errorf(`Expected %d elements, but got %d`, len(expected), len(got))
		t.Logf("\nExpected: %v\nGot:      %v", expected, got)
		return false
	}

	for i := range len(expected) {
		if !Equal(t, expected[i], got[i]) {
			return false
		}
	}

	return true
}

func NoError(t *testing.T, got error) bool {
	t.Helper()
	if got != nil {
		t.Logf("Unexpected error: %s", got)
		t.Fail()
		return false
	}

	return true
}

func Error(t *testing.T, got error) bool {
	t.Helper()
	if got == nil {
		t.Logf("Expected error: %s", got)
		t.Fail()
		return false
	}

	return true
}

func Panic(t *testing.T, assert func()) {
	t.Helper()
	defer func() {
		if m := recover(); m != nil {
			return
		}
		t.Logf(`Expected panic, but it did not happen!`)
		t.Fail()
	}()

	assert()
}

func NoPanic(t *testing.T, assert func()) {
	t.Helper()
	defer func() {
		if m := recover(); m != nil {
			t.Logf("Unexpected panic: %v", m)
			t.Fail()
		}
	}()

	assert()
}

// NoErrorEventually retries fn up to 10 times over the given duration.
// It returns true if any call to fn returns nil.
// If all attempts fail, it logs the last error and fails the test.
func NoErrorEventually(t *testing.T, duration time.Duration, fn func() error) bool {
	t.Helper()

	const attempts = 10
	interval := duration / attempts
	var lastErr error

	for i := 0; i < attempts; i++ {
		if err := fn(); err == nil {
			return true
		} else {
			lastErr = err
		}

		if i < attempts-1 {
			time.Sleep(interval)
		}
	}

	t.Logf("Unexpected error: %s", lastErr)
	t.Fail()
	return false
}

func NotNil(t *testing.T, got any) bool {
	t.Helper()
	if got == nil {
		t.Logf("Expected non-nil value, but got nil")
		t.Fail()
		return false
	}
	return true
}
