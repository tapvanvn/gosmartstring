package gosmartstring

import (
	"github.com/google/uuid"
	"github.com/tapvanvn/gotokenize"
)

func BuildInstructionDo(name string, params []IObject, context *SSContext) gotokenize.Token {
	token := gotokenize.Token{
		Type:    TokenSSInstructionDo,
		Content: name,
	}
	for _, obj := range params {
		registerName := uuid.New().String()
		context.RegisterObject(registerName, obj)
		token.Children.AddToken(gotokenize.Token{
			Type:    TokenSSRegistry,
			Content: registerName,
		})
	}
	return token
}
