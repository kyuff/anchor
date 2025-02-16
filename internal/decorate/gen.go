package decorate

//go:generate go tool moq -skip-ensure -pkg decorate_test -rm -out mocks_test.go . fullComponent starter setupper closer namer
