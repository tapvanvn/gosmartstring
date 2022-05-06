package gosmartstring

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

type SmarstringMeaning struct {
	*gotokenize.AbstractMeaning
}

func CreateSSMeaning() *SmarstringMeaning {

	meaning := CreateSSRawMeaning()
	return &SmarstringMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(
			gotokenize.NewPatternMeaning(meaning,
				buildSSLPatterns(),
				SSLIgnores,
				getSSLGlobalNested(),
				gotokenize.NoTokens),
		),
	}
}

func (meaning *SmarstringMeaning) Prepare(process *gotokenize.MeaningProcess) {

	meaning.AbstractMeaning.Prepare(process)

	fmt.Println("--1.0--")
	process.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
		ExtendTypeSize: 6,
	})
	fmt.Println("--end 1.0--")

	tmpStream := gotokenize.CreateStream(0)

	for {
		var token = meaning.AbstractMeaning.Next(process)
		if token == nil {
			break
		}

		if token.Type == TokenSSLSmartstring {

			tmpStream.AddToken(meaning.parseInstruction(token))

		} else if gotokenize.IsContainToken(getSSLGlobalNested(), token.Type) {
			meaning.prepareStream(token)
			tmpStream.AddToken(*token)
		} else if token.Type == TokenSSLCommand {
			meaning.parseCommand(token)
			tmpStream.AddToken(*token)
		} else {

			tmpStream.AddToken(*token)
		}
	}
	process.SetStream(process.Context.AncestorTokens, &tmpStream)
	fmt.Println("--1.1--")
	process.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
		ExtendTypeSize: 6,
	})
	fmt.Println("--end 1.1--")
}
func (meaning *SmarstringMeaning) parseCommand(parentToken *gotokenize.Token) {
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	tmpStream.AddToken(*parentToken.Children.GetTokenAt(0))
	second := parentToken.Children.GetTokenAt(1)

	meaning.parseParentThese(second)
	childIter := second.Children.Iterator()
	for {
		childToken := childIter.Read()
		if childToken == nil {
			break
		}
		tmpStream.AddToken(*childToken)
	}
	parentToken.Children = tmpStream
}

func (meaning *SmarstringMeaning) prepareStream(parentToken *gotokenize.Token) {

	tmpStream := gotokenize.CreateStream(0)
	iter := parentToken.Children.Iterator()
	for {
		token := iter.Read()
		if token == nil {
			break
		}

		if token.Type == TokenSSLSmartstring {

			tmpStream.AddToken(meaning.parseInstruction(token))

		} else if gotokenize.IsContainToken(getSSLGlobalNested(), token.Type) {
			meaning.prepareStream(token)
			tmpStream.AddToken(*token)
		} else if token.Type == TokenSSLCommand {
			meaning.parseCommand(token)
			tmpStream.AddToken(*token)
		} else {
			tmpStream.AddToken(*token)
		}

	}
	parentToken.Children = tmpStream
}
func (meaning *SmarstringMeaning) parseParentThese(parentToken *gotokenize.Token) {

	parentToken.Debug(5, SSNaming, &gotokenize.DebugOption{ExtendTypeSize: 6})

	iter := parentToken.Children.Iterator()
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	pack := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLSmartstring, "")
	for {
		meaningToken := meaning.getNextInstruction(iter)
		if meaningToken == nil {
			break
		}
		if meaningToken.Content == "," {
			fmt.Println("detect comma")
			if pack.Children.Length() == 1 {
				tmpStream.AddToken(*pack.Children.GetTokenAt(0))
			} else if pack.Children.Length() > 1 {
				meaning.prepareStream(pack)
				tmpStream.AddToken(*pack)
			}
			pack = gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLSmartstring, "")
			continue
		}
		pack.Children.AddToken(*meaningToken)
	}
	if pack.Children.Length() == 1 {
		tmpStream.AddToken(*pack.Children.GetTokenAt(0))
	} else if pack.Children.Length() > 1 {
		meaning.prepareStream(pack)
		tmpStream.AddToken(*pack)
	}
	tmpStream.Debug(10, SSNaming, &gotokenize.DebugOption{ExtendTypeSize: 6})
	parentToken.Children = tmpStream
}

func (meaning *SmarstringMeaning) parseInstruction(token *gotokenize.Token) gotokenize.Token {

	if token.Type != TokenSSLSmartstring {

		return *token
	}

	newToken := gotokenize.Token{
		Type: token.Type,
	}

	iter := token.Children.Iterator()

	for {
		meaningToken := meaning.getNextInstruction(iter)
		if meaningToken == nil {
			break
		}
		newToken.Children.AddToken(*meaningToken)
	}

	return newToken
}

func (meaning *SmarstringMeaning) getNextInstruction(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		token := iter.Read()
		if token == nil {
			break
		}

		if token.Type == TokenSSLWord || token.Type == TokenSSLCommand {

			instructionToken := &gotokenize.Token{

				Type: TokenSSLInstruction,
			}

			if token.Type == TokenSSLCommand {

				//fmt.Println("try to reach nested intructions")
				meaning.parseCommand(token)
				instructionToken.Children.AddToken(*token)

			} else {

				instructionToken.Children.AddToken(*token)
			}
			meaning.reachUntilEndInstruction(iter, instructionToken)

			return instructionToken

		} else {

			if gotokenize.IsContainToken(getSSLGlobalNested(), token.Type) {

				meaning.prepareStream(token)
			}
			return token
		}

	}
	return nil
}

func (meaning *SmarstringMeaning) reachUntilEndInstruction(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	for {
		if iter.EOS() {
			break
		}

		token := iter.Get()
		//TODO: move to a array of allowed operator notation
		if token.Content == "," || token.Content == "+" || token.Content == "." {

			/*if token.Content == "," {

				_ = iter.Read()
			}*/

			break

		} else if token.Type == TokenSSLWord || token.Type == TokenSSLCommand {

			currentToken.Children.AddToken(meaning.parseInstruction(token))
		}
		_ = iter.Read()

	}
}

func (meaning *SmarstringMeaning) GetName() string {
	return "ss_meaning"
}
