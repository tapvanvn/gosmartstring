package gosmartstring

type ssruntime struct {
	registries map[string]ssregistry
	extend     interface{}
}

func CreateRuntime(extend interface{}) *ssruntime {
	return &ssruntime{
		registries: map[string]ssregistry{},
		extend:     extend,
	}
}

func (runtime *ssruntime) RegisterObject(name string, object IObject) {

	runtime.registries[name] = CreateObjectRegistry(object)
}

func (runtime *ssruntime) RegisterFunction(name string, sfunc IFunction) {

	runtime.registries[name] = CreateFunctionRegistry(sfunc)
}

func (runtime *ssruntime) GetRegistry(name string) *ssregistry {

	if registry, ok := runtime.registries[name]; ok {

		return &registry
	}
	return nil
}
