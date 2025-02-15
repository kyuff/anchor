package assert_test

import (
	"testing"

	"github.com/kyuff/anchor/internal/assert"
)

func TestAsserts(t *testing.T) {
	var testCases = []struct {
		name   string
		assert func(t *testing.T)
		failed bool
	}{
		{
			name: "Equal success",
			assert: func(t *testing.T) {
				assert.Equal(t, 1, 1)
			},
			failed: false,
		},
		{
			name: "Equal failed",
			assert: func(t *testing.T) {
				assert.Equal(t, 1, 2)
			},
			failed: true,
		},
		{
			name: "Equal success",
			assert: func(t *testing.T) {
				assert.Equalf(t, 1, 1, "hello test")
			},
			failed: false,
		},
		{
			name: "Equalf failed",
			assert: func(t *testing.T) {
				assert.Equalf(t, 1, 2, "hello test")
			},
			failed: true,
		},
		{
			name: "EqualSlice success",
			assert: func(t *testing.T) {
				assert.EqualSlice(t, []int{1, 2, 3}, []int{1, 2, 3})
			},
			failed: false,
		},
		{
			name: "EqualSlice failed",
			assert: func(t *testing.T) {
				assert.EqualSlice(t, []int{1, 2, 3}, []int{})
			},
			failed: true,
		},
		{
			name: "Truef success",
			assert: func(t *testing.T) {
				assert.Truef(t, true, "hello test")
			},
			failed: false,
		},
		{
			name: "Truef failed",
			assert: func(t *testing.T) {
				assert.Truef(t, false, "hello test")
			},
			failed: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// arrange
			var (
				x = &testing.T{}
			)

			// act
			testCase.assert(x)

			// arrange
			assert.Equal(t, testCase.failed, x.Failed())
		})
	}
}
