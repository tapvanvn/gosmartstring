package gosmartstring

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/tapvanvn/gotokenize"
)

type SSContext struct {
	id      uuid.UUID
	Root    *SSContext
	Parent  *SSContext
	Level   int
	Runtime *SSRuntime
	This    IObject
	HotLink bool

	//not public
	result        []IObject
	registries    map[string]ssregistry
	registryCount int
	registryStack *SSAddressStack
}

type ssResultInfo struct {
	stream      *gotokenize.TokenStream
	beginOffset int
	endOffset   int
}

func CreateContext(runtime *SSRuntime) *SSContext {

	ctx := &SSContext{
		id:            uuid.New(),
		Level:         0,
		Parent:        nil,
		Runtime:       runtime,
		This:          nil,
		HotLink:       false,
		registries:    map[string]ssregistry{},
		registryCount: 0,
	}
	ctx.Root = ctx
	return ctx
}

func (ctx *SSContext) CreateSubContext() *SSContext {

	subContext := &SSContext{
		id:         uuid.New(),
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
	ctx.registryCount++
}

func (ctx *SSContext) RegisterFunction(name string, sfunc IFunction) {

	ctx.registries[name] = CreateFunctionRegistry(sfunc)
	ctx.registryCount++
}

func (ctx *SSContext) IssueAddress() string {
	ctx.registryCount++
	return strconv.Itoa(ctx.registryCount)
}

func (ctx *SSContext) GetRegistry(name string) *ssregistry {

	if registry, ok := ctx.registries[name]; ok {

		return &registry

	} else if ctx.Parent != nil {

		return ctx.Parent.GetRegistry(name)
	}
	return ctx.Runtime.GetRegistry(name)
}

func (ctx *SSContext) StackResult(addressType int, address string, result IObject) {

	if addressType == TokenSSRegistryGlobal {

		ctx.Root.RegisterObject(address, result)

		if ctx.registryStack != nil {

			ctx.registryStack.Append(address)
		}
	} else if addressType == TokenSSRegistry {

		ctx.RegisterObject(address, result)
	}
}

func (ctx *SSContext) SetStackRegistry(stack *SSAddressStack) {

	ctx.registryStack = stack
}

func (ctx *SSContext) BindingTo(context *SSContext) {
	if context != ctx {
		if context != nil {
			ctx.Root = context.Root
			ctx.Parent = context
		} else {
			ctx.Root = nil
			ctx.Parent = nil
		}
	}
}

func (ctx *SSContext) PrintDebug() {

	if ctx.Parent != nil {

		ctx.Parent.PrintDebug()
	}
	fmt.Println("context id:", ctx.id)
	for name, registry := range ctx.registries {

		if registry.Function != nil {
			fmt.Println(name, "function")
		} else if registry.Object != nil {
			fmt.Println(name, "object", registry.Object.GetType())
		} else {
			fmt.Println(name, "object nil")
		}
	}
}
