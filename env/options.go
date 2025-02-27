package env

import (
	"os"
)

type Config struct {
	appliers []func(kv map[string]string) error
	setter   func(key, val string) error
}

func defaultOptions() *Config {
	return applyOptions(&Config{},
		WithEnvironment(os.Setenv),
	)
}

func applyOptions(options *Config, opts ...Option) *Config {
	for _, opt := range opts {
		opt(options)
	}

	return options
}

type Option func(cfg *Config)

// WithEnvironment is not used in most cases. It is exposed primarily for testing.
func WithEnvironment(fn func(key, value string) error) Option {
	return func(cfg *Config) {
		cfg.setter = fn
	}
}

func WithT(t TestingT) Option {
	return func(cfg *Config) {
		cfg.setter = func(key, val string) error {
			t.Setenv(key, val)
			return nil
		}
	}
}

func OverrideFile(file string) Option {
	return func(cfg *Config) {
		cfg.appliers = append(cfg.appliers, overrideFile(file))
	}
}
func OverrideEnvKeyFile(key string) Option {
	return func(cfg *Config) {
		cfg.appliers = append(cfg.appliers, overrideFile(os.Getenv(key)))
	}
}
func OverrideKeyValue(key, value string) Option {
	return func(cfg *Config) {
		cfg.appliers = append(cfg.appliers, overrideKeyValue(key, value))
	}
}
func DefaultFile(file string) Option {
	return func(cfg *Config) {
		cfg.appliers = append(cfg.appliers, defaultFile(file))
	}
}
func DefaultEnvKeyFile(key string) Option {
	return func(cfg *Config) {
		cfg.appliers = append(cfg.appliers, defaultFile(os.Getenv(key)))
	}
}
func DefaultKeyValue(key, value string) Option {
	return func(cfg *Config) {
		cfg.appliers = append(cfg.appliers, defaultKeyValue(key, value))
	}
}
