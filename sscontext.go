package gosmartstring

type sscontext struct {
	Root    *sscontext
	Parent  *sscontext
	Level   int
	Runtime *ssruntime
	This    IObject
	HotLink bool
	//not public
	result     []byte
	registries map[string]ssregistry
}

func CreateContext(runtime *ssruntime) *sscontext {

	ctx := &sscontext{
		Level:   0,
		Parent:  nil,
		Runtime: runtime,
		This:    nil,
		HotLink: false,
	}
	ctx.Root = ctx
	return ctx
}

func (ctx *sscontext) CreateSubContext() *sscontext {

	subContext := &sscontext{

		Level:   ctx.Level + 1,
		Runtime: ctx.Runtime,
		Root:    ctx.Root,
		Parent:  ctx,
		This:    ctx.This,
		HotLink: ctx.HotLink,
	}
	return subContext
}

func (ctx *sscontext) RegisterObject(name string, object IObject) {

	ctx.registries[name] = CreateObjectRegistry(object)
}

func (ctx *sscontext) RegisterFunction(name string, sfunc IFunction) {

	ctx.registries[name] = CreateFunctionRegistry(sfunc)
}

func (ctx *sscontext) GetRegistry(name string) *ssregistry {

	if registry, ok := ctx.registries[name]; ok {

		return &registry

	} else if ctx.Parent != nil {

		return ctx.Parent.GetRegistry(name)
	}
	return ctx.Runtime.GetRegistry(name)
}

func (ctx *sscontext) StackResult(data []byte) {

	ctx.Root.result = append(ctx.Root.result, data...)
}
