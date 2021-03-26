package gosmartstring

import (
	"github.com/tapvanvn/gotokenize"
)

func BuildDo(name string, params []IObject, context *SSContext) gotokenize.Token {

	token := gotokenize.Token{
		Type:    TokenSSInstructionDo,
		Content: name,
	}
	addressToken := gotokenize.Token{
		Type: TokenSSRegistry,
	}
	if context == context.Root {
		addressToken.Type = TokenSSRegistryGlobal
	}
	addressToken.Content = context.IssueAddress()

	token.Children.AddToken(addressToken)

	for _, obj := range params {

		registerName := context.IssueAddress()

		context.RegisterObject(registerName, obj)

		token.Children.AddToken(gotokenize.Token{

			Type:    TokenSSRegistry,
			Content: registerName,
		})
	}
	return token
}

func BuildEach(arrayName string, varName string, actionsTokens []gotokenize.Token, context *SSContext) gotokenize.Token {
	token := gotokenize.Token{
		Type:    TokenSSInstructionEach,
		Content: arrayName,
	}
	addressToken := gotokenize.Token{
		Type: TokenSSRegistry,
	}
	if context == context.Root {
		addressToken.Type = TokenSSRegistryGlobal
	}
	addressToken.Content = context.IssueAddress()

	token.Children.AddToken(addressToken)
	//var address
	token.Children.AddToken(gotokenize.Token{
		Type:    TokenSSRegistry,
		Content: varName,
	})
	//instruction
	for _, insToken := range actionsTokens {
		token.Children.AddToken(insToken)
	}
	return token
}
