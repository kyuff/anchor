package assert

import (
	"fmt"
	"reflect"
	"testing"
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

func NotEqual[T comparable](t *testing.T, unexpected, got T) bool {
	t.Helper()
	if unexpected == got {
		t.Logf(`
Items was equal
Expected: %v
     Got: %v`, unexpected, got)
		t.Fail()
		return false
	}
	return true
}

func NotNil(t *testing.T, got any) bool {
	t.Helper()
	if reflect.ValueOf(got).IsNil() {
		t.Logf("Expected a value, but got nil")
		t.Fail()
		return false
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
		return false
	}

	for i := range len(expected) {
		if !Equal(t, expected[i], got[i]) {
			return false
		}
	}

	return true
}
