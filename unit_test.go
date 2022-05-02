package gosmartstring_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tapvanvn/gosmartstring"
	ss "github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotokenize/v2"
)

const (
	contentSimple  = `{{dic("x"), dic("y")}}`
	contentSimple2 = `{{dic("x"), dic("y")+put("z")}}`
	contentSimple3 = `{{single+put("z")}}`
)

var (
	compiler = gosmartstring.SSCompiler{}
	runtime  = createRuntime()
)

func printUtf8(content string) {
	for _, c := range content {
		fmt.Printf("%c", c)
	}
}

func SSFuncTestDo(context *gosmartstring.SSContext, input gosmartstring.IObject, params []gosmartstring.IObject) gosmartstring.IObject {

	fmt.Println("call SSFuncTestDo", len(params))
	if len(params) == 1 {

		if sstring, ok := params[0].(*gosmartstring.SSString); ok {

			id := sstring.Value

			fmt.Println("id", id)
		}
	}
	return nil
}
func SSFuncTestEach(context *gosmartstring.SSContext, input gosmartstring.IObject, params []gosmartstring.IObject) gosmartstring.IObject {
	fmt.Println("call SSFuncTestEach", len(params))
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
func SSFPut(context *ss.SSContext, input ss.IObject, params []ss.IObject) ss.IObject {
	//fmt.Println("call put")
	if len(params) == 1 {
		//fmt.Println("call put param")
		if name, ok := params[0].(*ss.SSString); ok {
			//fmt.Println("call put param2")
			formatedName := strings.TrimSpace(name.Value)
			if formatedName != "" {
				hotObject := context.HotObject()
				//if hotObject != nil {
				//hotContent := "unknown"
				//if str, ok := hotObject.(*ss.SSString); ok {
				//	hotContent = str.Value
				//}
				//fmt.Printf("put %s to %s\n,", hotContent, formatedName)

				//} else {
				//fmt.Println("put nil to ", formatedName)
				//}
				context.RegisterObject(formatedName, hotObject)
			}
		}
	}
	return input
}

func createRuntime() *gosmartstring.SSRuntime {
	runtime := gosmartstring.CreateRuntime(nil)
	runtime.RegisterFunction("testDo", SSFuncTestDo)
	runtime.RegisterFunction("testEach", SSFuncTestEach)
	runtime.RegisterFunction("put", SSFPut)
	return runtime
}

func TestSSInstructionDo(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	stream := gotokenize.CreateStream(0)
	instructionDo := gosmartstring.BuildDo("testDo",
		[]gosmartstring.IObject{
			gosmartstring.CreateString("test:html/index.html"),
		}, context)

	stream.AddToken(instructionDo)
	compiler := gosmartstring.SSCompiler{}

	compiler.Compile(&stream, context)

	context.PrintDebug(0)
}

func TestSSInstructionEach(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	stream := gotokenize.CreateStream(0)

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
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	meaning.Prepare(proc)
	token := meaning.Next(proc)
	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, token.Content)
		if token.Type == gosmartstring.TokenSSLSmartstring {
			token.Children.Debug(0, gosmartstring.SSNaming, &gotokenize.DebugOption{
				ExtendTypeSize: 6,
			})
		}
		token = meaning.Next(proc)
	}
}
func TestSSLMeaning2(t *testing.T) {

	content := "{{testDo(testDo3(\"hey\"), \"hello\").test, testDo(\"bcd\").test+put(a)}}"
	meaning := gosmartstring.CreateSSMeaning()
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	meaning.Prepare(proc)
	token := meaning.Next(proc)
	for {
		if token == nil {
			break
		}
		fmt.Println(token.Type, token.Content)
		if token.Type == gosmartstring.TokenSSLSmartstring {
			token.Children.Debug(0, nil, nil)
		}
		token = meaning.Next(proc)
	}
}

func TestSSLInstruction(t *testing.T) {
	context := gosmartstring.CreateContext(createRuntime())
	content := `normal {{testDo(testDo3("hey"), "hello").test, testDo("bcd").test+put(a)}}`
	meaning := gosmartstring.CreateSSInstructionMeaning()
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	proc.Context.BindingData = context
	meaning.Prepare(proc)

	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		token.Children.Debug(0, gosmartstring.SSNaming, &gotokenize.DebugOption{
			ExtendTypeSize: 6,
		})
	}
}

func TestJSONMeaning(t *testing.T) {
	jsonString := `{
		"attr1":1,
		"attr2":"attr2value"
		}`
	obj := gosmartstring.CreateSSJSON(jsonString)
	value := obj.Call(nil, "attr1", nil)
	if value != nil {
		if sstring, ok := value.(*gosmartstring.SSString); ok {
			fmt.Println(sstring.Value)
		} else {
			fmt.Println("not string")
			t.Fail()
		}
	} else {
		fmt.Println("value nil")
		t.Fail()
	}
}

type testInter struct {
	Name string `json:"name"`
}

func TestParseInterface(t *testing.T) {
	name := &testInter{
		Name: "testname",
	}
	names := []*testInter{
		name,
	}
	context := gosmartstring.CreateContext(createRuntime())
	context.RegisterInterface("name", name)
	context.RegisterInterface("names", names)
	context.PrintDebug(0)
}
func debugCompiledStream(iter *gotokenize.Iterator, context *gosmartstring.SSContext) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		if token.Type == gosmartstring.TokenSSInstructionDo || token.Type == gosmartstring.TokenSSInstructionExport {
			debugInstruction(token, context)
		} else if token.Children.Length() > 0 {
			childIter := token.Children.Iterator()
			debugCompiledStream(childIter, context)
		}
	}
}
func debugInstruction(token *gotokenize.Token, context *gosmartstring.SSContext) {
	iter := token.Children.Iterator()
	addressToken := iter.Get()
	if addressToken.Type == gosmartstring.TokenSSRegistryIgnore {
		return
	}
	fmt.Println("do:", token.Content, "address:", addressToken.Content)

	obj := context.GetRegistry(addressToken.Content)

	if obj != nil && obj.Object != nil {

		fmt.Println("found", obj.Object.GetType())

		if obj.Object.CanExport() {

			content := string(obj.Object.Export(context))
			fmt.Println(content)
		}
	} else {
		fmt.Println("-----------------")
		fmt.Println("not found")
		context.PrintDebug(0)
		fmt.Println("-----------------")
	}
}
func TestSSLInstructionJSON(t *testing.T) {
	name := &testInter{
		Name: "testname",
	}
	names := []*testInter{
		name,
	}
	context := gosmartstring.CreateContext(createRuntime())
	context.RegisterInterface("name", name)
	context.RegisterInterface("names", names)
	context.RegisterObject("str", gosmartstring.CreateString("stringvalue"))
	content := "{{str}}"
	meaning := gosmartstring.CreateSSInstructionMeaning()
	stream := gotokenize.CreateStream(0)
	stream.Tokenize(content)
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	proc.Context.BindingData = context
	meaning.Prepare(proc)

	compileStream := gotokenize.CreateStream(0)
	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		compileStream.AddToken(*token)
	}
	compileStream.Debug(10, gosmartstring.SSNaming, nil)
	//compileStream.Debug(0, nil)

	compiler := gosmartstring.SSCompiler{}

	compiler.Compile(&compileStream, context)

	//ontext.PrintDebug(0)
	iter := compileStream.Iterator()
	debugCompiledStream(iter, context)
}

func TestCompileSimple(t *testing.T) {

	gosmartstring.SSInsructionMove(5000)
	context := gosmartstring.CreateContext(runtime)
	//context.DebugLevel = 1

	dic := gosmartstring.CreateSSStringMap()
	dic.Set("x", gosmartstring.CreateString("x_value"))
	dic.Set("y", gosmartstring.CreateString("y_value"))
	context.RegisterObject("dic", dic)
	context.RegisterObject("single", gosmartstring.CreateString("single_value"))

	stream := gotokenize.CreateStream(0)
	stream.Tokenize(contentSimple3)
	proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
	proc.Context.BindingData = context

	meaning := gosmartstring.CreateSSInstructionMeaning()
	meaning.Prepare(proc)
	compileStream := gotokenize.CreateStream(0)
	for {
		token := meaning.Next(proc)
		if token == nil {
			break
		}
		compileStream.AddToken(*token)
	}
	if err := compiler.Compile(&compileStream, context); err != nil {
		context.PrintDebug(0)
		t.Fatal(err)
	}
	fmt.Println("--final--")
	compileStream.Debug(0, ss.SSNaming, &gotokenize.DebugOption{
		ExtendTypeSize: 6,
	})
	gotokenize.DebugMeaning(meaning)
	fmt.Println("--end final--")
}
