package gosmartstring

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

var _context_id = 0

type SSContext struct {
	id      int
	Root    *SSContext
	Parent  *SSContext
	Level   int
	Runtime *SSRuntime
	This    IObject

	//not public
	hotLink       bool
	remember      bool
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
	return fmt.Sprint(ctx.id)
}
func CreateContext(runtime *SSRuntime) *SSContext {
	_context_id++
	ctx := &SSContext{
		id:            _context_id,
		Level:         0,
		Parent:        nil,
		Runtime:       runtime,
		This:          nil,
		hotLink:       false,
		remember:      false,
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
	fmt.Printf("[%d] registerObject %s type:%s", ctx.id, name, object.GetType())
	ctx.registries[finalAddress] = CreateObjectRegistry(object)
}

func (ctx *SSContext) RegisterInterface(name string, object interface{}) {

	parseObject := ParseInterface(object)

	finalAddress := name
	if ctx.registryStack != nil {
		finalAddress = ctx.IssueAddress()
		ctx.registryStack.Append(name, finalAddress)
	}

	ctx.registries[finalAddress] = CreateObjectRegistry(parseObject)
}

func (ctx *SSContext) DebugCurrentStack() {
	if ctx.registryStack != nil {
		for address, translate := range ctx.registryStack.Address[ctx.registryStack.offset] {

			fmt.Println("stack", ctx.registryStack.offset, address, "->", translate)
		}
	} else {
		fmt.Println("no stack binding")
	}
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
	return fmt.Sprintf("%d", ctx.registryCount)
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
			ctx.registryStack = context.registryStack
		} else {
			ctx.Root = nil
			ctx.Parent = nil
			ctx.registryStack = nil
		}
	}
}

func (ctx *SSContext) PrintDebug(level int) {

	for i := 0; i <= level; i++ {
		fmt.Print("|")
		if i == 0 {
			fmt.Printf("%s ", gotokenize.ColorContent("ctx"))
		}
		fmt.Print(" ")
	}
	fmt.Println(gotokenize.ColorName(fmt.Sprint(ctx.id)))

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
