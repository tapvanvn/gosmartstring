package gosmartstring

import "github.com/tapvanvn/gotokenize/v2"

func NewSSInstructionPack(token gotokenize.Token) *SSIntructionPack {
	return &SSIntructionPack{
		IObject: &SSObject{},
		Pack:    token,
	}
}

type SSIntructionPack struct {
	IObject
	Pack gotokenize.Token
}

func (pack *SSIntructionPack) GetType() string {
	return "ss_inspack"
}
