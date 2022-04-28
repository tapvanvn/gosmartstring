package gosmartstring

import (
	"github.com/tapvanvn/gotokenize/v2"
)

type SmarstringInstructionMeaning struct {
	*gotokenize.AbstractMeaning
	//SmarstringMeaning
}

func CreateSSInstructionMeaning() SmarstringInstructionMeaning {
	return SmarstringInstructionMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(CreateSSMeaning()),
	}
}

func (meaning *SmarstringInstructionMeaning) Prepare(proc *gotokenize.MeaningProcess, context *SSContext) {

	meaning.AbstractMeaning.Prepare(proc)

	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		token := meaning.AbstractMeaning.Next(proc)
		if token == nil {
			break
		}
		if token.Type == TokenSSLSmarstring {
			tmpStream.AddToken(meaning.buildSmarstring(token, context))
		} else {
			tmpStream.AddToken(*token)
		}
	}
	proc.SetStream(proc.Context.AncestorTokens, &tmpStream)
}

func (meaning *SmarstringInstructionMeaning) buildSmarstring(token *gotokenize.Token, context *SSContext) gotokenize.Token {
	packToken := gotokenize.Token{
		Type: TokenSSInstructionPack,
	}
	iter := token.Children.Iterator()
	for {
		insToken := iter.Read()
		if insToken == nil {
			break
		}

		meaning.buildInstruction(insToken, &packToken, context)
	}
	return packToken
}

func (meaning *SmarstringInstructionMeaning) buildInstruction(token *gotokenize.Token, packToken *gotokenize.Token, context *SSContext) {
	iter := token.Children.Iterator()
	lastInstructionNum := packToken.Children.Length()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}

		if childToken.Type == TokenSSLCommand {
			meaning.buildCommand(childToken, false, packToken, context)
		} else if childToken.Content == "+" {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionLink,
			})
		} else if childToken.Content != "" {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionRemember,
			})
			doToken := gotokenize.Token{
				Type:    TokenSSInstructionDo,
				Content: childToken.Content,
			}
			doToken.Children.AddToken(gotokenize.Token{
				Type:    TokenSSRegistryIgnore,
				Content: context.IssueAddress(),
			})
			packToken.Children.AddToken(doToken)
		}
	}
	if packToken.Children.Length() > lastInstructionNum {
		//add and export here
		exportToken := gotokenize.Token{
			Type:    TokenSSInstructionExport,
			Content: "",
		}
		exportToken.Children.AddToken(gotokenize.Token{
			Type:    TokenSSRegistry,
			Content: context.IssueAddress(),
		})
		packToken.Children.AddToken(exportToken)
	}
}

func (meaning *SmarstringInstructionMeaning) buildCommand(token *gotokenize.Token, isParam bool, packToken *gotokenize.Token, context *SSContext) string {
	iter := token.Children.Iterator()

	nameToken := iter.Read()
	if nameToken == nil || nameToken.Content == "" {
		return ""
	}
	params := []gotokenize.Token{}
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenSSLParenthese {

			childIter := childToken.Children.Iterator()

			for {
				childToken2 := childIter.Read()
				if childToken2 == nil {
					break
				}
				if childToken2.Type == TokenSSLCommand {

					address := meaning.buildCommand(childToken2, true, packToken, context)
					if address != "" {
						params = append(params, gotokenize.Token{
							Type:    TokenSSRegistry,
							Content: address,
						})
					}
				} else {
					address := context.IssueAddress()
					paramToken := gotokenize.Token{
						Type:    TokenSSRegistry,
						Content: address,
					}
					value := ""
					if childToken2.Type == 0 {
						value = childToken2.Content
					} else if childToken2.Type == TokenSSLString {
						value = childToken2.Children.ConcatStringContent()
					}

					context.RegisterObject(address, CreateString(value))

					params = append(params, paramToken)
				}
			}
		}
	}

	addressType := TokenSSRegistryIgnore
	if isParam {

		addressType = TokenSSRegistry
	}

	cmdAddress := context.IssueAddress()
	doToken := gotokenize.Token{
		Type:    TokenSSInstructionDo,
		Content: nameToken.Content,
	}
	doToken.Children.AddToken(gotokenize.Token{
		Type:    addressType,
		Content: cmdAddress,
	})
	for _, param := range params {

		doToken.Children.AddToken(param)
	}
	packToken.Children.AddToken(gotokenize.Token{
		Type: TokenSSInstructionRemember,
	})
	packToken.Children.AddToken(doToken)

	return cmdAddress
}
