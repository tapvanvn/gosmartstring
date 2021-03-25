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
		registryStack: nil,
	}
	ctx.Root = ctx
	return ctx
}

/*
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
*/

func (ctx *SSContext) RegisterObject(name string, object IObject) {

	finalAddress := name
	if ctx.registryStack != nil {
		finalAddress = ctx.IssueAddress()
		ctx.registryStack.Append(name, finalAddress)
	}

	if object != nil {
		content := ""
		if sstring, ok := object.(*SSString); ok {
			content = sstring.Value
			if len(content) > 30 {
				content = content[:30]
			}
		}
		if ctx.registryStack != nil {
			fmt.Println(ctx.ID(), "stack result:", name, "->", finalAddress, object.GetType(), content)
		} else {
			fmt.Println(ctx.ID(), "stack result:", finalAddress, object.GetType(), content)
		}
	} else {
		if ctx.registryStack != nil {
			fmt.Println(ctx.ID(), "stack result nil:", name, "->", finalAddress)
		} else {
			fmt.Println(ctx.ID(), "stack result nil:", finalAddress, object.GetType())
		}
	}

	ctx.registries[finalAddress] = CreateObjectRegistry(object)
}

func (ctx *SSContext) RegisterFunction(name string, sfunc IFunction) {

	finalAddress := name
	if ctx.registryStack != nil {
		finalAddress = ctx.IssueAddress()
		ctx.registryStack.Append(name, finalAddress)
	}

	ctx.registries[finalAddress] = CreateFunctionRegistry(sfunc)
}

func (ctx *SSContext) IssueAddress() string {
	ctx.registryCount++
	return strconv.Itoa(ctx.registryCount)
}

func (ctx *SSContext) GetRegistry(name string) *ssregistry {
	var address = name
	if ctx.registryStack != nil {
		if translate, ok := ctx.registryStack.Get(name); ok {
			address = translate
		}
	}

	if registry, ok := ctx.registries[address]; ok {

		return &registry

	} else if ctx.Parent != nil {

		return ctx.Parent.GetRegistry(address)
	}

	return ctx.Runtime.GetRegistry(address)
}

func (ctx *SSContext) StackResult(addressType int, address string, result IObject) {

	if addressType == TokenSSRegistryGlobal {

		ctx.Root.RegisterObject(address, result)

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
			content := ""
			if sstring, ok := registry.Object.(*SSString); ok {
				content = sstring.Value
				if len(sstring.Value) > 30 {
					content = sstring.Value[:30]
				}
			}
			fmt.Printf("%s-%s : %s", name, gotokenize.ColorContent(registry.Object.GetType()), content)
		} else {
			fmt.Printf("%s-%s", name, gotokenize.ColorContent("nil"))
		}

		fmt.Println()
	}
	if ctx.Parent != nil {

		ctx.Parent.PrintDebug(level + 1)
	}
}
