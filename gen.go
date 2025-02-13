package anchor

//go:generate go tool moq -skip-ensure -pkg anchor_test -rm -out mocks_test.go . Component fullComponent
