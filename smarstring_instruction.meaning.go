package gosmartstring

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

type SmarstringInstructionMeaning struct {
	*gotokenize.AbstractMeaning
}

func CreateSSInstructionMeaning() *SmarstringInstructionMeaning {
	meaning := &SmarstringInstructionMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(CreateSSMeaning()),
	}
	return meaning
}

func (meaning *SmarstringInstructionMeaning) Prepare(proc *gotokenize.MeaningProcess) {

	context := proc.Context.BindingData.(*SSContext)

	meaning.AbstractMeaning.Prepare(proc)

	// fmt.Println("--1--")
	// proc.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
	// 	ExtendTypeSize: 6,
	// })
	// fmt.Println("--end 1--")

	tmpStream := gotokenize.CreateStream(0)

	for {
		token := meaning.AbstractMeaning.Next(proc)
		if token == nil {
			break
		}
		if token.Type == TokenSSLSmartstring {
			tmpStream.AddToken(meaning.buildSmarstring(token, context))
		} else {
			tmpStream.AddToken(*token)
		}
	}
	proc.SetStream(proc.Context.AncestorTokens, &tmpStream)
}

//process smartstring
func (meaning *SmarstringInstructionMeaning) buildSmarstring(token *gotokenize.Token, context *SSContext) gotokenize.Token {
	// fmt.Println("--build smartstring--")
	// token.Debug(0, SSNaming, &gotokenize.DebugOption{
	// 	ExtendTypeSize: 6,
	// })
	// fmt.Println("--end smartstring--")
	packToken := gotokenize.Token{
		Type: TokenSSInstructionPack,
	}
	iter := token.Children.Iterator()

	for {

		insToken := iter.Read()
		if insToken == nil {
			break
		}
		if insToken.Type == TokenSSLInstruction {
			meaning.packInstruction(insToken, &packToken, context)
		} else if insToken.Content == "+" {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionLink,
			})
		} else if insToken.Content == "." {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionReload, //reload the last returned
			})
		}
	}

	return packToken
}

//each instruction is a doToken, but we need determine the pre and post actions of the call
func (meaning *SmarstringInstructionMeaning) packInstruction(token *gotokenize.Token, packToken *gotokenize.Token, context *SSContext) {
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
		} else if childToken.Content == "." {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionReload, //reload the last returned
			})
		} else if childToken.Type == TokenSSLWord {
			//fmt.Println("--build word--")
			//childToken.Debug(0, SSNaming, &gotokenize.DebugOption{
			//	ExtendTypeSize: 6,
			//})
			//fmt.Println("--end word--")

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
					/*	} else if childToken2.Type == TokenSSLSmartstring {

						postToken := meaning.buildSmarstring(childToken2, context)
						postIter := postToken.Children.Iterator()
						for {
							paramToken := postIter.Read()
							if paramToken == nil {
								break
							}
							params = append(params, *paramToken)
						}*/

				} else if childToken2.Content != "," {
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
	if context.DebugLevel > 0 {
		fmt.Println("--do--")
		doToken.Debug(0, SSNaming, &gotokenize.DebugOption{
			ExtendTypeSize: 6,
		})
		context.PrintDebug(0)
		fmt.Println("--end do--")
	}
	//packToken.Children.AddToken(gotokenize.Token{
	//	Type: TokenSSInstructionRemember,
	//})
	packToken.Children.AddToken(doToken)

	return cmdAddress
}
func (meaning *SmarstringInstructionMeaning) GetMeaningLevel() int {

	return meaning.AbstractMeaning.GetMeaningLevel() + 1
}

func (meaning *SmarstringInstructionMeaning) Propagate(fn func(meaning gotokenize.IMeaning)) {

	fn(meaning)

	meaning.AbstractMeaning.Propagate(fn)
}

func (meaning *SmarstringInstructionMeaning) GetName() string {
	return "ss_instruction_meaning"
}
