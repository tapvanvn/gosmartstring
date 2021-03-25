package gosmartstring

import (
	"fmt"

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

func (ctx *SSContext) ID() string {
	return ctx.id.String()
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
	return uuid.NewString()
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

	if result != nil {

		fmt.Println(ctx.ID(), "stack result:", address, result.GetType())
	}
	if addressType == TokenSSRegistryGlobal {

		ctx.Root.RegisterObject(address, result)

	} else if addressType == TokenSSRegistry {

		ctx.RegisterObject(address, result)
	}
	if ctx.registryStack != nil {

		ctx.registryStack.Append(address)
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

func (ctx *SSContext) PrintDebug(level int) {

	for i := 0; i <= level; i++ {
		fmt.Print("|")
		if i == 0 {
			fmt.Printf("%s ", gotokenize.ColorName("ctx"))
		}
		fmt.Print(" ")
	}
	fmt.Println(ctx.id)

	for name, registry := range ctx.registries {

		for i := 0; i <= level; i++ {
			fmt.Print("|")
			if i == 0 {
				if registry.Function != nil {
					fmt.Printf("%s ", gotokenize.ColorName("fnc"))
				} else {
					fmt.Printf("%s ", gotokenize.ColorName("obj"))
				}
			}
			fmt.Print(" ")
		}
		if registry.Function != nil {
			fmt.Printf("%s", name)
		} else if registry.Object != nil {
			fmt.Printf("%s-%s", name, gotokenize.ColorContent(registry.Object.GetType()))
		} else {
			fmt.Printf("%s-%s", name, gotokenize.ColorContent("nil"))
		}

		fmt.Println()
	}
	if ctx.Parent != nil {

		ctx.Parent.PrintDebug(level + 1)
	}
}
