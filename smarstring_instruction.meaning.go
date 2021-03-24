package gosmartstring

import "github.com/tapvanvn/gotokenize"

type SmarstringInstructionMeaning struct {
	SmarstringMeaning
}

func CreateSSInstructionMeaning() SmarstringInstructionMeaning {
	return SmarstringInstructionMeaning{
		SmarstringMeaning: CreateSSMeaning(),
	}
}

func (meaning *SmarstringInstructionMeaning) Prepare(stream *gotokenize.TokenStream, context *SSContext) {

	meaning.SmarstringMeaning.Prepare(stream)
	tmpStream := gotokenize.CreateStream()
	token := meaning.SmarstringMeaning.Next()
	for {
		if token == nil {
			break
		}
		if token.Type == TokenSSLSmarstring {
			tmpStream.AddToken(meaning.buildSmarstring(token, context))
		} else {
			tmpStream.AddToken(*token)
		}
		token = meaning.SmarstringMeaning.Next()
	}
	meaning.SetStream(tmpStream)
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
		meaning.buildIntruction(insToken, &packToken, context)
	}
	return packToken
}

func (meaning *SmarstringInstructionMeaning) buildIntruction(token *gotokenize.Token, packToken *gotokenize.Token, context *SSContext) {
	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenSSLCommand {
			meaning.buildCommand(childToken, packToken, context)
		} else if childToken.Content == "+" {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionLink,
			})
		} else if childToken.Content != "" {
			doToken := gotokenize.Token{
				Type:    TokenSSInstructionDo,
				Content: childToken.Content,
			}
			doToken.Children.AddToken(gotokenize.Token{
				Type:    TokenSSRegistry,
				Content: context.IssueAddress(),
			})
			packToken.Children.AddToken(doToken)
		}
	}
}

func (meaning *SmarstringInstructionMeaning) buildCommand(token *gotokenize.Token, packToken *gotokenize.Token, context *SSContext) string {
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

					address := meaning.buildCommand(childToken2, packToken, context)
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
	cmdAddress := context.IssueAddress()
	doToken := gotokenize.Token{
		Type:    TokenSSInstructionDo,
		Content: nameToken.Content,
	}
	doToken.Children.AddToken(gotokenize.Token{
		Type:    TokenSSRegistry,
		Content: cmdAddress,
	})
	for _, param := range params {

		doToken.Children.AddToken(param)
	}
	packToken.Children.AddToken(doToken)

	return cmdAddress
}
