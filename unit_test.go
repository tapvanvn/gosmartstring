package gosmartstring_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize"
)

func SSFuncTestDo(context *gosmartstring.SSContext, input gosmartstring.IObject, params []gosmartstring.IObject) gosmartstring.IObject {

	fmt.Println("call SSFuncTestDo", len(params))
	if len(params) == 1 {

		if sstring, ok := params[0].(*gosmartstring.SSString); ok {

			id := sstring.Value

			fmt.Print("id", id)
		}
	}
	return nil
}
func SSFuncTestEach(context *gosmartstring.SSContext, input gosmartstring.IObject, params []gosmartstring.IObject) gosmartstring.IObject {
	fmt.Println("call SSFuncTestDo", len(params))
	if element := context.GetRegistry("element"); element != nil && element.Object != nil {
		if sstring, ok := element.Object.(*gosmartstring.SSString); ok {
			fmt.Println("element", sstring.Value)
			return element.Object
		} else {
			fmt.Println("element is not a string")
		}
	} else {
		//TODO report error
		fmt.Println("element not found")
	}
	return nil
}

func createRuntime() *gosmartstring.SSRuntime {
	runtime := gosmartstring.CreateRuntime(nil)
	runtime.RegisterFunction("testDo", SSFuncTestDo)
	runtime.RegisterFunction("testEach", SSFuncTestEach)
	return runtime
}

func TestSSInstructionDo(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	stream := gotokenize.CreateStream()
	instructionDo := gosmartstring.BuildDo("testDo",
		[]gosmartstring.IObject{
			gosmartstring.CreateString("test:html/index.html"),
		}, context)

	stream.AddToken(instructionDo)
	compiler := gosmartstring.SSCompiler{}

	compiler.Compile(&stream, context)

	context.PrintDebug()
}

func TestSSInstructionEach(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	stream := gotokenize.CreateStream()

	ssArray := gosmartstring.CreateSSArray()
	ssArray.Stack = append(ssArray.Stack, gosmartstring.CreateString("e1"))
	ssArray.Stack = append(ssArray.Stack, gosmartstring.CreateString("e2"))

	context.RegisterObject("testArray", ssArray)

	instructionDo := gosmartstring.BuildDo("testEach", []gosmartstring.IObject{}, context)

	instructionEach := gosmartstring.BuildEach("testArray", "element", []gotokenize.Token{instructionDo}, context)

	stream.AddToken(instructionEach)

	compiler := gosmartstring.SSCompiler{}

	compiler.Compile(&stream, context)

	//context.PrintDebug()
}

func TestSSLMeaning(t *testing.T) {

	content := "{{testDo(\"abc\").test, testDo(\"bcd\").test+put(a)}}"
	meaning := gosmartstring.CreateSSMeaning()
	stream := gotokenize.CreateStream()
	stream.Tokenize(content)
	meaning.Prepare(&stream)
	token := meaning.Next()
	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, token.Content)
		if token.Type == gosmartstring.TokenSSLSmarstring {
			token.Children.Debug(0, nil)
		}
		token = meaning.Next()
	}
}
func TestSSLMeaning2(t *testing.T) {

	content := "{{testDo(testDo3(\"hey\"), \"hello\").test, testDo(\"bcd\").test+put(a)}}"
	meaning := gosmartstring.CreateSSMeaning()
	stream := gotokenize.CreateStream()
	stream.Tokenize(content)
	meaning.Prepare(&stream)
	token := meaning.Next()
	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, token.Content)
		if token.Type == gosmartstring.TokenSSLSmarstring {
			token.Children.Debug(0, nil)
		}
		token = meaning.Next()
	}
}

func TestSSLInstruction(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	content := "{{testDo(testDo3(\"hey\"), \"hello\").test, testDo(\"bcd\").test+put(a)}}"
	meaning := gosmartstring.CreateSSInstructionMeaning(context)
	stream := gotokenize.CreateStream()
	stream.Tokenize(content)
	meaning.Prepare(&stream)
	token := meaning.Next()
	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, token.Content)

		token.Children.Debug(0, nil)

		token = meaning.Next()
	}
}
