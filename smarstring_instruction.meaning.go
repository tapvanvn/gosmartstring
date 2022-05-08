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

/*
func (meaning *SmarstringInstructionMeaning) Prepare(proc *gotokenize.MeaningProcess) {

	context := proc.Context.BindingData.(*SSContext)

	meaning.AbstractMeaning.Prepare(proc)

	fmt.Println("--2--")
	proc.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{

		ExtendTypeSize: 6,
	})
	fmt.Println("--end 2--")

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
	fmt.Println("--2.1--")
	proc.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
		ExtendTypeSize: 6,
	})
	fmt.Println("--end 2.1--")
}*/
func (meaning *SmarstringInstructionMeaning) Next(proc *gotokenize.MeaningProcess) *gotokenize.Token {
	token := meaning.getNextMeaningToken(proc)
	return token
}
func (meaning *SmarstringInstructionMeaning) processChild(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {
	proc := gotokenize.NewMeaningProcessFromStream(append(context.AncestorTokens, parentToken.Type), &parentToken.Children)

	newStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		token := meaning.Next(proc)

		if token == nil {

			break
		}
		newStream.AddToken(*token)
	}
	parentToken.Children = newStream
}

func (meaning *SmarstringInstructionMeaning) getNextMeaningToken(proc *gotokenize.MeaningProcess) *gotokenize.Token {
	token := proc.Iter.Read()
	if token != nil {
		sscontext := proc.Context.BindingData.(*SSContext)
		if token.Type == TokenSSLSmartstring {
			meaning.buildSmarstring(token, sscontext)
		}
	}
	return token
}

//process smartstring
func (meaning *SmarstringInstructionMeaning) buildSmarstring(token *gotokenize.Token, sscontext *SSContext) {
	fmt.Println("--build smartstring--")
	token.Debug(0, SSNaming, &gotokenize.DebugOption{
		ExtendTypeSize: 6,
	})
	fmt.Println("--end smartstring--")

	packToken := gotokenize.Token{

		Type: TokenSSInstructionPack,
	}
	iter := token.Children.Iterator()

	for {

		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenSSLWord {

			meaning.buildDoInstruction(childToken, "", sscontext)
			packToken.Children.AddToken(*childToken)

		} else if childToken.Type == TokenSSLCommand {

			meaning.buildCommand(childToken, "", sscontext)
			packToken.Children.AddToken(*childToken)

		} else if childToken.Content == "+" {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionLink,
			})
		} else if childToken.Content == "." {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionReload, //reload the last returned
			})
		} else if childToken.Content == "*" {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionReload,
			})
			exportToken := gotokenize.Token{
				Type:    TokenSSInstructionExport,
				Content: "",
			}

			exportToken.Children.AddToken(gotokenize.Token{
				Type:    TokenSSRegistry,
				Content: sscontext.IssueAddress(),
			})

			packToken.Children.AddToken(exportToken)
		} else if childToken.Type == TokenSSLBreak {
			packToken.Children.AddToken(gotokenize.Token{
				Type: TokenSSInstructionReset,
			})
		} else {
			// fmt.Println("--err token")
			// childToken.Debug(0, SSNaming, &gotokenize.DebugOption{ExtendTypeSize: 6})
			// fmt.Println("--end err token")
		}
	}

	*token = packToken

	// fmt.Println("--after build smartstring--")
	// packToken.Debug(0, SSNaming, &gotokenize.DebugOption{
	// 	ExtendTypeSize: 6,
	// })
	// fmt.Println("--end after  smartstring--")

}

//each instruction is a doToken, but we need determine the pre and post actions of the call
func (meaning *SmarstringInstructionMeaning) buildDoInstruction(wordToken *gotokenize.Token, outputAddress string, sscontext *SSContext) {

	wordToken.Type = TokenSSInstructionDo

	wordToken.Children.AddToken(gotokenize.Token{
		Type:    TokenSSRegistryIgnore,
		Content: sscontext.IssueAddress(),
	})

	/*if outputAddress != "" {
		//add and export here
		exportToken := gotokenize.Token{
			Type:    TokenSSInstructionExport,
			Content: "",
		}
		outputAddress = sscontext.IssueAddress()
		exportToken.Children.AddToken(gotokenize.Token{
			Type:    TokenSSRegistry,
			Content: outputAddress,
		})
		wordToken.Children.AddToken(exportToken)
	}*/
}

func (meaning *SmarstringInstructionMeaning) buildCommand(token *gotokenize.Token, outputAddress string, sscontext *SSContext) {

	iter := token.Children.Iterator()

	nameToken := iter.Read()
	if nameToken == nil || nameToken.Content == "" {
		return
	}

	params := []gotokenize.Token{}
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}

		if childToken.Type == TokenSSLParenthese {

			//pack := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLSmartstring, "")
			childIter := childToken.Children.Iterator()

			for {
				childToken2 := childIter.Read()
				if childToken2 == nil {
					break
				}

				if childToken2.Type == TokenSSLCommand {

					address := sscontext.IssueAddress()
					meaning.buildCommand(childToken2, address, sscontext)

					/*params = append(params, gotokenize.Token{
						Type:    TokenSSRegistry,
						Content: address,
					})*/

					params = append(params, *childToken2)
				} else if childToken2.Type == TokenSSLSmartstring {

					meaning.buildSmarstring(childToken2, sscontext)
					params = append(params, *childToken2)

				} else if childToken2.Type == TokenSSLWord {

					address := sscontext.IssueAddress()
					meaning.buildDoInstruction(childToken2, address, sscontext)
					params = append(params, *childToken2)

				} else if childToken2.Type == TokenSSLString {

					address := sscontext.IssueAddress()
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

					sscontext.RegisterObject(address, CreateString(value))

					params = append(params, paramToken)
				}
			}
		}
	}

	addressType := TokenSSRegistryIgnore
	cmdAddress := ""
	/*if outputAddress != "" {

		addressType = TokenSSRegistry
		cmdAddress = sscontext.IssueAddress()
	}*/
	newToken := gotokenize.Token{
		Type:    TokenSSInstructionDo,
		Content: nameToken.Content,
	}
	newToken.Children.AddToken(gotokenize.Token{
		Type:    addressType,
		Content: cmdAddress,
	})
	for _, param := range params {

		newToken.Children.AddToken(param)
	}
	if sscontext.DebugLevel > 0 {
		fmt.Println("--do--")
		newToken.Debug(0, SSNaming, &gotokenize.DebugOption{
			ExtendTypeSize: 6,
		})
		sscontext.PrintDebug(0)
		fmt.Println("--end do--")
	}
	*token = newToken
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
