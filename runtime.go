package gosmartstring

type Runtime struct {
	root       *Runtime
	parent     *Runtime
	registries map[string]Registry
	extend     interface{}
}

func CreateRuntime(extend interface{}) Runtime {
	return Runtime{
		root:       nil,
		parent:     nil,
		registries: map[string]Registry{},
		extend:     extend,
	}
}

func (runtime *Runtime) RegisterObject(name string, object IObject) {

	runtime.registries[name] = CreateObjectRegistry(object)
}

func (runtime *Runtime) RegisterFunction(name string, sfunc IFunction) {

	runtime.registries[name] = CreateFunctionRegistry(sfunc)
}

func (runtime *Runtime) CreateSubRuntime() Runtime {
	subRuntime := &Runtime{
		parent:     runtime.parent,
		registries: map[string]Registry{},
		extend:     runtime.extend,
	}
	if runtime.root != nil {
		subRuntime.root = runtime.root
	} else {
		subRuntime.root = runtime
	}
	return *subRuntime
}

func (runtime *Runtime) GetRegistry(name string) *Registry {

	if registry, ok := runtime.registries[name]; ok {

		return &registry
	}
	return nil
}

//TODO: complete compile
func (runtime *Runtime) Compile(content string) {

}
