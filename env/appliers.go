package env

func overrideFile(file string) func(kv map[string]string) error {
	values, err := readFile(file)
	if err != nil {
		return func(kv map[string]string) error {
			return err
		}
	}

	return func(kv map[string]string) error {
		for k, v := range values {
			// cannot fail
			_ = overrideKeyValue(k, v)(kv)
		}

		return nil
	}
}

func overrideKeyValue(key string, value string) func(kv map[string]string) error {
	return func(kv map[string]string) error {
		kv[key] = value
		return nil
	}
}

func defaultFile(file string) func(kv map[string]string) error {
	values, err := readFile(file)
	if err != nil {
		return func(kv map[string]string) error {
			return err
		}
	}

	return func(kv map[string]string) error {
		for k, v := range values {
			// cannot fail
			_ = defaultKeyValue(k, v)(kv)
		}

		return nil
	}
}

func defaultKeyValue(key string, value string) func(kv map[string]string) error {
	return func(kv map[string]string) error {
		_, ok := kv[key]
		if ok {
			return nil
		}

		kv[key] = value
		return nil
	}
}
