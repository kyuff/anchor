package decorate

func Make[T starter](name string, setup func() (T, error)) *Component {
	var component = &Component{}
	return makeComponent(component, name, setup, probeInner(component))
}
