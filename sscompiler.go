package gosmartstring

import (
	"github.com/tapvanvn/gotokenize"
)

type SSICompiler interface {
	Compile(token *gotokenize.Token, ctx *sscontext)
}
type SSCompiler struct {
}

func (compiler *SSCompiler) Compile(token *gotokenize.Token, context *sscontext) {
	switch token.Type {
	case TokenSSInstructionLink:

	case TokenSSInstructionDo:
		compiler.compileDo(token, context)
	case TokenSSInstructionPack:
		compiler.compilePack(token, context)
	case TokenSSInstructionExport:
		compiler.compileExport(token, context)
	case TokenSSInstructionIf:
		compiler.compileIf(token, context)
	case TokenSSInstructionCase:
		compiler.compileCase(token, context)
	case TokenSSInstructionCount:
		compiler.compileCount(token, context)
	default:
		//not support
		break
	}
}
func (compiler *SSCompiler) compileLink(token *gotokenize.Token, context *sscontext) {

	context.HotLink = true
}

func (compiler *SSCompiler) compilePack(token *gotokenize.Token, context *sscontext) {
	subContext := context.CreateSubContext()
	iter := token.Children.Iterator()
	for {
		if iter.EOS() {
			break
		}
		childToken := iter.Read()
		compiler.Compile(childToken, subContext)
	}
}

func (compiler *SSCompiler) compileDo(token *gotokenize.Token, context *sscontext) {

	name := token.Content
	if len(name) > 0 {
		iter := token.Children.Iterator()
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
	}
}

func (compiler *SSCompiler) compileExport(token *gotokenize.Token, context *sscontext) {

}

func (compiler *SSCompiler) compileIf(token *gotokenize.Token, context *sscontext) {

}

func (compiler *SSCompiler) compileCase(token *gotokenize.Token, context *sscontext) {

}

func (compiler *SSCompiler) compileEach(token *gotokenize.Token, context *sscontext) {

}

func (compiler *SSCompiler) compileCount(token *gotokenize.Token, context *sscontext) {

}

func (compiler *SSCompiler) callRegistry(name string, params []IObject, context *sscontext) {

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
