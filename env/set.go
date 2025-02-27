package env

type TestingT interface {
	Setenv(key, val string)
}

func Set(options ...Option) error {
	cfg := applyOptions(defaultOptions(), options...)

	var kv = readOS()
	for _, apply := range cfg.appliers {
		err := apply(kv)
		if err != nil {
			return err
		}
	}

	for key, value := range kv {
		if err := cfg.setter(key, value); err != nil {
			return err
		}
	}

	return nil
}

func MustSet(options ...Option) {
	if err := Set(options...); err != nil {
		panic(err)
	}
}
