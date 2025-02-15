package logger

import "log/slog"

//go:generate go tool moq -pkg logger_test -rm -out mocks_test.go . Handler

type Handler slog.Handler
