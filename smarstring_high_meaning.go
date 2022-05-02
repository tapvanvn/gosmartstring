package gosmartstring

import (
	"github.com/tapvanvn/gotokenize/v2"
)

type SmarstringMeaning struct {
	*gotokenize.AbstractMeaning
	//gotokenize.PatternMeaning
}

func CreateSSMeaning() *SmarstringMeaning {
	meaning := CreateSSRawMeaning()
	return &SmarstringMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(
			gotokenize.NewPatternMeaning(meaning, SSLPatterns, SSLIgnores, SSLGlobalNested, gotokenize.NoTokens),
		),
	}
}

func (meaning *SmarstringMeaning) Prepare(process *gotokenize.MeaningProcess) {

	meaning.AbstractMeaning.Prepare(process)

	//fmt.Println("--1.0--")
	//process.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
	//	ExtendTypeSize: 6,
	//})
	//fmt.Println("--end 1.0--")

	tmpStream := gotokenize.CreateStream(0)

	var token = meaning.AbstractMeaning.Next(process)

	for {
		if token == nil {
			break
		}

		if token.Type == TokenSSLSmartstring {

			tmpStream.AddToken(meaning.parseInstruction(token))

		} else {

			tmpStream.AddToken(*token)
		}

		token = meaning.AbstractMeaning.Next(process)
	}
	process.SetStream(process.Context.AncestorTokens, &tmpStream)
}

func (meaning *SmarstringMeaning) parseInstruction(token *gotokenize.Token) gotokenize.Token {

	if token.Type != TokenSSLSmartstring && token.Type != TokenSSLParenthese {

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

				//debugPrint("try to reach nested intructions")
				instructionToken.Children.AddToken(meaning.parseInstruction(token))

			} else {

				instructionToken.Children.AddToken(*token)
			}
			meaning.reachUntilEndInstruction(iter, instructionToken)

			return instructionToken

		} else {

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

		if token.Content == "," || token.Content == "+" {

			if token.Content == "," {

				_ = iter.Read()
			}

			break

		} else if token.Type == TokenSSLWord || token.Type == TokenSSLCommand {

			currentToken.Children.AddToken(meaning.parseInstruction(token))
		}
		_ = iter.Read()

	}
}
