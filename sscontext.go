package gosmartstring

import (
	"fmt"
	"strconv"

	"github.com/tapvanvn/gotokenize/v2"
)

var _context_id = 0
var _address_id = int64(0)

//NOTE: translate address is effected on the current context not the base
type SSContext struct {
	id      int
	Root    *SSContext
	Parent  *SSContext
	Level   int
	Runtime *SSRuntime
	This    IObject

	//not public
	hotLink    bool
	remember   bool
	result     []IObject
	registries map[string]ssregistry
	//registryCount int
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
		id:         _context_id,
		Level:      0,
		Parent:     nil,
		Runtime:    runtime,
		This:       nil,
		hotLink:    false,
		remember:   false,
		registries: map[string]ssregistry{},
		//registryCount: 0,
		registryStack: nil,
	}

	ctx.Root = ctx
	return ctx
}

func (ctx *SSContext) Register(name string, resitry ssregistry) {
	finalAddress := name

	if ctx.registryStack != nil {
		finalAddress = ctx.IssueAddress()
		ctx.registryStack.Append(name, finalAddress)
	}
	//fmt.Printf("res obj: %s name:%s final:%s\n", ctx.ID(), name, finalAddress)
	ctx.registries[finalAddress] = resitry
}

func (ctx *SSContext) RegisterObject(name string, object IObject) {

	ctx.Register(name, CreateObjectRegistry(object))
}

func (ctx *SSContext) RegisterInterface(name string, object interface{}) {

	parseObject := ParseInterface(object)

	ctx.Register(name, CreateObjectRegistry(parseObject))

}

func (ctx *SSContext) RegisterFunction(name string, sfunc IFunction) {

	ctx.Register(name, CreateFunctionRegistry(sfunc))

}

func (ctx *SSContext) IssueAddress() string {
	_address_id++
	return fmt.Sprintf("%d", _address_id)
}

func (ctx *SSContext) GetRegistry(name string) *ssregistry {
	var address = name
	if ctx.registryStack != nil {
		if translate, ok := ctx.registryStack.Get(name); ok {
			address = translate
		}
	}

	if registry, ok := ctx.registries[address]; ok {
		//fmt.Printf("get %s name:%s address:%s\n", ctx.ID(), name, address)
		return &registry

	} else if ctx.Parent != nil {
		//fmt.Printf("getparent %s name:%s address:%s\n", ctx.ID(), name, address)
		return ctx.Parent.GetRegistry(name)
	}

	return ctx.Runtime.GetRegistry(address)
}

func (ctx *SSContext) StackResult(addressType int, address string, result IObject) {

	if addressType == TokenSSRegistryGlobal {

		ctx.Root.RegisterObject(address, result)

	} else if addressType == TokenSSRegistry {

		ctx.RegisterObject(address, result)

	} else if addressType != TokenSSRegistryIgnore {

		fmt.Println("unknown address type", addressType)
	}
}

func (ctx *SSContext) SetStackRegistry(stack *SSAddressStack) {

	ctx.registryStack = stack
}

//TODO: should we apply push, pop to ctx, so it can auto change to last environment.
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
func (ctx *SSContext) DebugCurrentStack(level int) {

	for i := 0; i <= level; i++ {
		fmt.Print("|")
		if i == 0 {
			fmt.Printf("%s ", gotokenize.ColorContent("stk"))
		}
		fmt.Print(" ")
	}
	fmt.Println()
	if ctx.registryStack != nil {

		for address, translate := range ctx.registryStack.Address[ctx.registryStack.offset] {
			for i := 0; i <= level; i++ {
				fmt.Print("|")
				if i == 0 {

					fmt.Printf("%s ", gotokenize.ColorContent(strconv.Itoa(ctx.registryStack.offset)))
				}
				fmt.Print(" ")
			}
			fmt.Printf("stack ori:%s-> trans:%s\n", address, translate)
		}
	} else {

		for i := 0; i <= level; i++ {
			fmt.Print("|")
			fmt.Print(" ")
		}
		fmt.Println("no stack binding ")
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
	ctx.DebugCurrentStack(level)
	if ctx.Parent != nil {

		ctx.Parent.PrintDebug(level + 1)
	}
}
