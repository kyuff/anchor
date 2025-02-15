package anchor

import (
	"errors"
	"fmt"
	"sync"
)

// Singleton creates a func that returns a value T created by fn.
// The value is created only once. If fn returns an error, Singleton panics.
//
// This function is meant to be used as part of application wiring, where
// a panic seems may seem acceptable.
//
// It is similar to sync.OnceValue, but differs in the way that the creator
// function fn can return an error.
func Singleton[T any](fn func() (T, error)) func() T {
	var (
		once  sync.Once
		value T
		err   error
		valid bool
		init  = func() {
			defer func() {
				msg := recover()
				if !valid {
					err = errors.Join(err, fmt.Errorf("panic: %v", msg))
					panic(err)
				}
			}()
			value, err = fn()
			if err != nil {
				valid = false
			} else {
				fn = nil
				valid = true
			}
		}
	)

	return func() T {
		once.Do(init)
		if !valid {
			panic(err)
		}

		return value
	}
}
