package gosmartstring

type SSRuntime struct {
	registries map[string]ssregistry
	extend     interface{}
}

func CreateRuntime(extend interface{}) *SSRuntime {
	return &SSRuntime{
		registries: map[string]ssregistry{},
		extend:     extend,
	}
}

func (runtime *SSRuntime) RegisterObject(name string, object IObject) {

	runtime.registries[name] = CreateObjectRegistry(object)
}

func (runtime *SSRuntime) RegisterFunction(name string, sfunc IFunction) {

	runtime.registries[name] = CreateFunctionRegistry(sfunc)
}

func (runtime *SSRuntime) GetRegistry(name string) *ssregistry {

	if registry, ok := runtime.registries[name]; ok {

		return &registry
	}
	return nil
}
