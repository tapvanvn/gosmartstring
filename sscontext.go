package gosmartstring

type SSContext struct {
	Root    *SSContext
	Parent  *SSContext
	Level   int
	Runtime *SSRuntime
	This    IObject
	HotLink bool
	//not public
	result     []byte
	registries map[string]ssregistry
}

func CreateContext(runtime *SSRuntime) *SSContext {

	ctx := &SSContext{
		Level:      0,
		Parent:     nil,
		Runtime:    runtime,
		This:       nil,
		HotLink:    false,
		result:     make([]byte, 0),
		registries: map[string]ssregistry{},
	}
	ctx.Root = ctx
	return ctx
}

func (ctx *SSContext) CreateSubContext() *SSContext {

	subContext := &SSContext{

		Level:      ctx.Level + 1,
		Runtime:    ctx.Runtime,
		Root:       ctx.Root,
		Parent:     ctx,
		This:       ctx.This,
		HotLink:    ctx.HotLink,
		registries: map[string]ssregistry{},
	}
	return subContext
}

func (ctx *SSContext) RegisterObject(name string, object IObject) {

	ctx.registries[name] = CreateObjectRegistry(object)
}

func (ctx *SSContext) RegisterFunction(name string, sfunc IFunction) {

	ctx.registries[name] = CreateFunctionRegistry(sfunc)
}

func (ctx *SSContext) GetRegistry(name string) *ssregistry {

	if registry, ok := ctx.registries[name]; ok {

		return &registry

	} else if ctx.Parent != nil {

		return ctx.Parent.GetRegistry(name)
	}
	return ctx.Runtime.GetRegistry(name)
}

func (ctx *SSContext) StackResult(data []byte) {

	ctx.Root.result = append(ctx.Root.result, data...)
}
