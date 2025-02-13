package anchor

import "context"

type Config struct {
	logger Logger
	ctx    context.Context
}

func defaultOptions() *Config {
	return applyOptions(&Config{},
		// add default options here
		WithNoopLogger(),
		withBackgroundContext(),
	)

}

func applyOptions(options *Config, opts ...Option) *Config {
	for _, opt := range opts {
		opt(options)
	}

	return options
}
