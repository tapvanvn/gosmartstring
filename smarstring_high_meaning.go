package gosmartstring

import "github.com/tapvanvn/gotokenize"

type SmarstringMeaning struct {
	gotokenize.PatternMeaning
}

func CreateSmarstringMeaning() SmarstringMeaning {
	meaning := CreateSmarstringRawMeaning()
	return SmarstringMeaning{
		PatternMeaning: gotokenize.CreatePatternMeaning(&meaning, SSLPatterns, SSLIgnores, SSLGlobalNested),
	}
}

func (meaning *SmarstringMeaning) Prepare(stream *gotokenize.TokenStream) {

	meaning.PatternMeaning.Prepare(stream)

	tmpStream := gotokenize.CreateStream()
	//SmartStringOnly := true

	var token = meaning.PatternMeaning.Next()

	for {
		if token == nil {
			break
		}

		if token.Type == TokenSSLSmarstring {

			tmpStream.AddToken(meaning.parseInstruction(token))

		} else {

			//SmartStringOnly = false

			tmpStream.AddToken(*token)
		}
		token = meaning.PatternMeaning.Next()
	}

	meaning.SetStream(tmpStream)
}

func (meaning *SmarstringMeaning) parseInstruction(token *gotokenize.Token) gotokenize.Token {

	newToken := gotokenize.Token{
		Type: token.Type,
	}

	iter := token.Children.Iterator()

	meaningToken := meaning.getNextInstruction(&iter)

	for {
		if meaningToken == nil {
			break
		}

		newToken.Children.AddToken(*meaningToken)

		meaningToken = meaning.getNextInstruction(&iter)
	}

	return newToken

}

func (meaning *SmarstringMeaning) getNextInstruction(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		token := iter.Read()
		if token == nil {
			break
		}

		if token.Type == 0 || token.Type == TokenSSLCommand {

			instructionToken := &gotokenize.Token{
				Type: TokenSSInstructionDo,
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

		} else if token.Type == 0 || token.Type == TokenSSLCommand {

			//debugPrint("try to reach nested intructions 2")
			currentToken.Children.AddToken(meaning.parseInstruction(token))
		}
		_ = iter.Read()

	}

}
