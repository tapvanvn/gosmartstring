package gosmartstring

import (
	"fmt"
	"strconv"

	"github.com/tapvanvn/gotokenize"
)

type SSContext struct {
	Root    *SSContext
	Parent  *SSContext
	Level   int
	Runtime *SSRuntime
	This    IObject
	HotLink bool
	//not public
	result        []ssResultInfo
	registries    map[string]ssregistry
	registryCount int
}

type ssResultInfo struct {
	stream      *gotokenize.TokenStream
	beginOffset int
	endOffset   int
}

func CreateContext(runtime *SSRuntime) *SSContext {

	ctx := &SSContext{
		Level:         0,
		Parent:        nil,
		Runtime:       runtime,
		This:          nil,
		HotLink:       false,
		result:        make([]ssResultInfo, 0),
		registries:    map[string]ssregistry{},
		registryCount: 0,
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

func (ctx *SSContext) StackResult(stream *gotokenize.TokenStream, beginOffset int, endOffset int) {
	ctx.Root.result = append(ctx.Root.result, ssResultInfo{
		stream:      stream,
		beginOffset: beginOffset,
		endOffset:   endOffset,
	})
}

func (ctx *SSContext) PrintDebug() {
	if ctx.Parent != nil {
		ctx.Parent.PrintDebug()
	}
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
