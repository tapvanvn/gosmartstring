package gosmartstring

import (
	"errors"
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

type SSICompiler interface {
	Compile(stream *gotokenize.TokenStream, ctx *SSContext) error
}

type SSCompiler struct {
}

func (compiler *SSCompiler) Compile(stream *gotokenize.TokenStream, context *SSContext) error {

	iter := stream.Iterator()

	for {

		token := iter.Read()

		if token == nil {

			break
		}
		if err := compiler.CompileToken(token, context); err != nil {
			fmt.Println("compile error :", err.Error())
			return err
		}
	}
	return nil
}

func (compiler *SSCompiler) CompileToken(token *gotokenize.Token, context *SSContext) error {

	switch token.Type {
	case TokenSSInstructionLink:
		return compiler.compileLink(token, context)
	case TokenSSInstructionRemember:
		return compiler.compileRemember(token, context)
	case TokenSSInstructionDo:
		return compiler.compileDo(token, context)
	case TokenSSInstructionPack:
		return compiler.compilePack(token, context)
	case TokenSSInstructionExport:
		return compiler.compileExport(token, context)
	case TokenSSInstructionIf:
		return compiler.compileIf(token, context)
	case TokenSSInstructionCase:
		return compiler.compileCase(token, context)
	case TokenSSInstructionCount:
		return compiler.compileCount(token, context)
	case TokenSSInstructionEach:
		return compiler.compileEach(token, context)
	default:
		return compiler.Compile(&token.Children, context)
	}
}

func (compiler *SSCompiler) compileLink(token *gotokenize.Token, context *SSContext) error {

	context.hotLink = true
	return nil
}
func (compiler *SSCompiler) compileRemember(token *gotokenize.Token, context *SSContext) error {

	context.remember = true
	return nil
}

func (compiler *SSCompiler) compilePack(token *gotokenize.Token, context *SSContext) error {

	//subContext := context.CreateSubContext()

	if err := compiler.Compile(&token.Children, context); err != nil {

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

		childToken := iter.Read()
		if childToken == nil {
			break
		}

		if childToken.Type == TokenSSRegistry {

			param := context.GetRegistry(childToken.Content)

			if param != nil && param.Object != nil {

				params = append(params, param.Object)
			} else {

				return errors.New("registry not found " + childToken.Content)
			}
		}
	}

	if err := compiler.callRegistry(name, params, context); err != nil {

		return err
	}

	context.StackResult(output.Type, output.Content, context.This)
	if !context.remember {

		context.This = nil
	}
	context.remember = false
	return nil
}

func (compiler *SSCompiler) compileExport(token *gotokenize.Token, context *SSContext) error {

	iter := token.Children.Iterator()
	output := iter.Read()
	if output == nil {

		return errors.New("instruction do syntax error")
	}

	context.StackResult(output.Type, output.Content, context.This)
	if !context.remember {

		context.This = nil
	}
	context.remember = false
	return nil
}

func (compiler *SSCompiler) compileIf(token *gotokenize.Token, context *SSContext) error {
	return nil
}

func (compiler *SSCompiler) compileCase(token *gotokenize.Token, context *SSContext) error {

	return nil
}

func (compiler *SSCompiler) compileEach(token *gotokenize.Token, context *SSContext) error {
	//#0 - array name
	//#1 - output address
	//#2 - element name
	//#... - content
	arrayName := token.Content
	iter := token.Children.Iterator()
	output := iter.Read()
	elementNameToken := iter.Read()

	if len(arrayName) == 0 || output == nil || elementNameToken == nil || elementNameToken.Content == "" {

		return errors.New("instruction each syntax error")
	}
	elementName := elementNameToken.Content

	arrayObject := context.GetRegistry(arrayName)
	array, ok := arrayObject.Object.(*SSArray)

	if !ok {

		return errors.New("instruction each error " + arrayName + " is not an array")
	}
	addressStack := CreateAddressStack()
	context.SetStackRegistry(&addressStack)

	offset := iter.Offset
	for _, element := range array.Stack {

		context.RegisterObject(elementName, element)

		iter.Seek(offset)

		for {
			childToken := iter.Read()
			if childToken == nil {
				break
			}
			if err := compiler.CompileToken(childToken, context); err != nil {
				//TODO: should we issue an error token instead ?
				return err
			}
		}
		addressStack.Inc()
	}
	context.SetStackRegistry(nil)

	context.StackResult(output.Type, output.Content, &addressStack)
	return nil
}

func (compiler *SSCompiler) compileCount(token *gotokenize.Token, context *SSContext) error {
	return nil
}

func (compiler *SSCompiler) callRegistry(name string, params []IObject, context *SSContext) error {

	var rs IObject = nil

	if !context.hotLink && context.This != nil {

		rs = context.This.Call(context, name, params)

	} else {

		registry := context.GetRegistry(name)

		if registry == nil {
			//TODO: report registry nil
			return errors.New("cannot reach registry " + name)

		} else if registry.Function != nil {

			rs = registry.Function(context, context.This, params)

		} else if registry.Object != nil && len(params) == 0 {

			rs = registry.Object

		} else {

			return errors.New("registry unknown: " + name)
		}
	}
	context.This = rs
	return nil
}
