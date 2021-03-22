package gosmartstring_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/gosmartstring"
)

func SSFuncTest(context *gosmartstring.SSContext, input gosmartstring.IObject, params []gosmartstring.IObject) gosmartstring.IObject {

	fmt.Println("call SSFuncTest", len(params))
	if len(params) == 1 {

		if sstring, ok := params[0].(gosmartstring.SSString); ok {

			id := sstring.Value

			fmt.Print("id", id)
		}
	}
	return nil
}

func createRuntime() *gosmartstring.SSRuntime {
	runtime := gosmartstring.CreateRuntime(nil)
	runtime.RegisterFunction("template", SSFuncTest)
	return runtime
}

func TestSSInstruction(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	instructionDo := gosmartstring.BuildInstructionDo("template",
		[]gosmartstring.IObject{
			gosmartstring.CreateString("test:html/index.html"),
		}, context)

	compiler := gosmartstring.SSCompiler{}
	compiler.Compile(&instructionDo, context)
}
