package gosmartstring

import (
	"errors"

	"github.com/tapvanvn/gotokenize"
)

type SSICompiler interface {
	Compile(stream *gotokenize.TokenStream, ctx *SSContext)
}

type SSCompiler struct {
}

func (compiler *SSCompiler) Compile(stream *gotokenize.TokenStream, context *SSContext) error {
	iter := stream.Iterator()
	for {
		if iter.EOS() {
			break
		}
		var err error = nil
		token := iter.Read()
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
		default:
			err = compiler.Compile(&token.Children, context)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
func (compiler *SSCompiler) reportResult(addressType int, name string, context *SSContext) error {
	if addressType == TokenSSRegistryIgnore {
		return nil
	}
	if addressType == TokenSSRegistryGlobal {
		context.Root.RegisterObject(name, context.This)
		return nil
	}
	context.RegisterObject(name, context.This)
	return nil
}

func (compiler *SSCompiler) compileLink(token *gotokenize.Token, context *SSContext) error {

	context.HotLink = true
	return nil
}

func (compiler *SSCompiler) compilePack(token *gotokenize.Token, context *SSContext) error {

	subContext := context.CreateSubContext()

	err := compiler.Compile(&token.Children, subContext)
	if err != nil {
		return err
	}
	//Todo: gather result
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
			}
		}
	}
	compiler.callRegistry(name, params, context)
	compiler.reportResult(output.Type, output.Content, context)

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
	return nil
}

func (compiler *SSCompiler) compileCount(token *gotokenize.Token, context *SSContext) error {
	return nil
}

func (compiler *SSCompiler) callRegistry(name string, params []IObject, context *SSContext) {

	var rs IObject = nil

	if context.This != nil {

		rs = context.This.Call(context, name, params)
	}
	registry := context.GetRegistry(name)

	if registry.Function != nil {

		rs = registry.Function(context, context.This, params)

	} else if registry.Object != nil && len(params) == 0 {

		rs = registry.Object
	}

	context.This = rs
}
