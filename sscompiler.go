package gosmartstring

import (
	"errors"
	"fmt"

	"github.com/tapvanvn/gotokenize"
)

type SSICompiler interface {
	Compile(stream *gotokenize.TokenStream, ctx *SSContext) error
}

type SSCompiler struct {
}

func (compiler *SSCompiler) Compile(stream *gotokenize.TokenStream, context *SSContext) error {
	iter := stream.Iterator()
	for {
		if iter.EOS() {
			break
		}

		token := iter.Read()
		if err := compiler.CompileToken(token, context); err != nil {
			return err
		}
	}
	return nil
}

func (compiler *SSCompiler) CompileToken(token *gotokenize.Token, context *SSContext) error {
	var err error = nil

	switch token.Type {
	case TokenSSInstructionLink:

	case TokenSSInstructionDo:
		err = compiler.compileDo(token, context)
	case TokenSSInstructionPack:
		err = compiler.compilePack(token, context)
	case TokenSSInstructionExport:
		err = compiler.compileExport(token, context)
	case TokenSSInstructionIf:
		err = compiler.compileIf(token, context)
	case TokenSSInstructionCase:
		err = compiler.compileCase(token, context)
	case TokenSSInstructionCount:
		err = compiler.compileCount(token, context)
	case TokenSSInstructionEach:
		err = compiler.compileEach(token, context)
	default:
		err = compiler.Compile(&token.Children, context)
	}
	if err != nil {
		return err
	}
	return nil
}

func (compiler *SSCompiler) compileLink(token *gotokenize.Token, context *SSContext) error {

	context.HotLink = true
	return nil
}

func (compiler *SSCompiler) compilePack(token *gotokenize.Token, context *SSContext) error {

	//subContext := context.CreateSubContext()

	err := compiler.Compile(&token.Children, context)
	if err != nil {
		return err
	}
	return nil
}

func (compiler *SSCompiler) compileDo(token *gotokenize.Token, context *SSContext) error {

	name := token.Content
	iter := token.Children.Iterator()
	output := iter.Read()

	if len(name) == 0 || output == nil {

		return errors.New("instruction do syntax error")
	}

	params := []IObject{}
	for {
		if iter.EOS() {
			break
		}
		childToken := iter.Read()

		if childToken.Type == TokenSSRegistry {

			param := context.GetRegistry(childToken.Content)

			if param != nil && param.Object != nil {

				params = append(params, param.Object)
			} else {
				fmt.Println("registry not found " + childToken.Content)
				context.PrintDebug(0)
				return errors.New("registry not found " + childToken.Content)
			}
		}
	}
	compiler.callRegistry(name, params, context)
	context.StackResult(output.Type, output.Content, context.This)

	return nil
}

func (compiler *SSCompiler) compileExport(token *gotokenize.Token, context *SSContext) error {
	return nil
}

func (compiler *SSCompiler) compileIf(token *gotokenize.Token, context *SSContext) error {
	return nil
}

func (compiler *SSCompiler) compileCase(token *gotokenize.Token, context *SSContext) error {

	return nil
}

func (compiler *SSCompiler) compileEach(token *gotokenize.Token, context *SSContext) error {

	arrayName := token.Content
	iter := token.Children.Iterator()
	output := iter.Read()
	elementNameToken := iter.Read()

	if len(arrayName) == 0 || output == nil || elementNameToken == nil || elementNameToken.Content == "" {

		return errors.New("instruction each syntax error")
	}
	elementName := elementNameToken.Content

	fmt.Println("arrayName:", arrayName, "elementName:", elementName)

	arrayObject := context.GetRegistry(arrayName)
	array, ok := arrayObject.Object.(*SSArray)

	if !ok {
		fmt.Println("instruction each error " + arrayName + " is not an array")
		return errors.New("instruction each error " + arrayName + " is not an array")
	}
	addressStack := CreateAddressStack()
	context.SetStackRegistry(&addressStack)
	defer context.SetStackRegistry(nil)

	fmt.Println("array elements num:", len(array.Stack))
	offset := iter.Offset

	for _, element := range array.Stack {

		context.RegisterObject(elementName, element)
		fmt.Println("ins-each set ", elementName, element.GetType())
		iter.Seek(offset)
		for {
			childToken := iter.Read()
			if childToken == nil {
				break
			}
			err := compiler.CompileToken(childToken, context)
			if err != nil {
				return err
			}
		}
		addressStack.Inc()
	}

	context.StackResult(output.Type, output.Content, &addressStack)
	return nil
}

func (compiler *SSCompiler) compileCount(token *gotokenize.Token, context *SSContext) error {
	return nil
}

func (compiler *SSCompiler) callRegistry(name string, params []IObject, context *SSContext) {

	fmt.Println("registry call:", gotokenize.ColorName(name), len(params))
	var rs IObject = nil

	if context.This != nil {

		rs = context.This.Call(context, name, params)
	}
	registry := context.GetRegistry(name)

	if registry == nil {
		//TODO: report registry nil
		fmt.Println("cannot reach registry " + name)
	} else if registry.Function != nil {

		rs = registry.Function(context, context.This, params)

	} else if registry.Object != nil && len(params) == 0 {

		rs = registry.Object

		fmt.Println("rs", rs.GetType())

	} else {

		fmt.Println("registry call fail")
	}

	context.This = rs
}
