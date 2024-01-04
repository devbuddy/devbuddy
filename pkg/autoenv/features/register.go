package features

type Register map[string]Feature

func (e Register) Register(feature Feature) {
	if _, ok := e[feature.Name()]; ok {
		panic("Can't re-register a definition: " + feature.Name())
	}
	e[feature.Name()] = feature
}

func (e Register) Get(name string) Feature {
	return e[name]
}

func (e Register) Names() (names []string) {
	for name := range e {
		names = append(names, name)
	}
	return
}
