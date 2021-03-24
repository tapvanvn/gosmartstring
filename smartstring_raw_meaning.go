package gosmartstring

import "github.com/tapvanvn/gotokenize"

const ()

type SmarstringRawMeaning struct {
	gotokenize.RawMeaning
	IsSmartstringOnly bool
}

func CreateSmarstringRawMeaning() SmarstringRawMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{
		".(){}[]+-*/\\,\"'": {TokenType: TokenSSLOperator, Separate: true},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)
	return SmarstringRawMeaning{
		RawMeaning:        meaning,
		IsSmartstringOnly: false,
	}

}

func (meaning *SmarstringRawMeaning) Prepare(stream *gotokenize.TokenStream) {

	meaning.RawMeaning.Prepare(stream)

	tmpStream := gotokenize.CreateStream()
	meaning.IsSmartstringOnly = true
	iter := meaning.GetIter()
	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()

		var reach = false

		if token.Content == "{" {
			token2 := iter.GetBy(1)
			if token2 != nil && token2.Content == "{" {
				iter.Read()
				iter.Read()
				newToken := meaning.meaningSmartString(meaning.continueSmartString(iter))
				tmpStream.AddToken(newToken)

				reach = true
			}

		} else if token.Content == "\\" {

			token2 := iter.GetBy(1)

			if token2 != nil && (token2.Content == "{" || token2.Content == "}" || token2.Content == "\"" || token2.Content == "'") {

				iter.Read()
				token.Content = token2.Content
			}
		}

		if !reach {

			meaning.IsSmartstringOnly = false

			_ = iter.Read()

			token.Type = TokenSSLNormalstring
			tmpStream.AddToken(*token)
		}
	}

}

func (meaning *SmarstringRawMeaning) continueSmartString(iter *gotokenize.Iterator) gotokenize.Token {

	rsToken := gotokenize.Token{
		Type: TokenSSLSmarstring,
	}

	for {
		if iter.EOS() {
			break
		}
		token := iter.Get()

		if token.Content == "}" {
			token2 := iter.GetBy(1)
			if token2 != nil && token2.Content == "}" {

				_ = iter.Read()
				_ = iter.Read()
				break
			}
		} else if token.Content == "\\" {
			token2 := iter.GetBy(1)
			if token2 != nil && (token2.Content == "{" || token2.Content == "}" || token2.Content == "\"" || token2.Content == "'") {

				_ = iter.Read()
				token.Content = token2.Content
			}
		}
		_ = iter.Read()
		rsToken.Children.AddToken(*token)

	}
	return rsToken
}

func (meaning *SmarstringRawMeaning) meaningSmartString(token gotokenize.Token) gotokenize.Token {

	iter := token.Children.Iterator()

	newToken := gotokenize.Token{
		Type: TokenSSLSmarstring,
	}

	meaningToken := meaning.getNextMeaningToken(&iter)

	for {
		if meaningToken == nil {
			break
		}
		newToken.Children.AddToken(*meaningToken)

		meaningToken = meaning.getNextMeaningToken(&iter)
	}
	return newToken
}

func (meaning *SmarstringRawMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()

		if token.Content == "{" {

			tmpToken := &gotokenize.Token{
				Content: "{",
				Type:    TokenSSLBlock,
			}

			meaning.continueUntil(iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "(" {

			tmpToken := &gotokenize.Token{
				Content: "(",
				Type:    TokenSSLParenthese,
			}

			meaning.continueUntil(iter, tmpToken, ")")

			return tmpToken

		} else if token.Content == "[" {

			tmpToken := &gotokenize.Token{
				Content: "[",
				Type:    TokenSSLSquare,
			}

			meaning.continueUntil(iter, tmpToken, "]")

			return tmpToken

		} else if token.Content == "\"" || token.Content == "'" {

			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenSSLString,
			}
			meaning.continueReadString(iter, tmpToken, token.Content)
			return tmpToken
		}

		return token

	}
	return nil
}

func (meaning *SmarstringRawMeaning) continueUntil(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

	var specialCharacter = false
	var lastSpecialToken *gotokenize.Token = nil

	for {

		token := meaning.getNextMeaningToken(iter)
		if token == nil {
			break
		}
		if token.Content == "\\" {

			specialCharacter = !specialCharacter
			lastSpecialToken = token

		} else if token.Content == reach {

			if specialCharacter {

				specialCharacter = false
				currentToken.Children.AddToken(*token)

			} else {

				break
			}
		} else {

			if specialCharacter {

				currentToken.Children.AddToken(*lastSpecialToken)
			}
			specialCharacter = false
			currentToken.Children.AddToken(*token)
		}
	}
}

func (meaning *SmarstringRawMeaning) continueReadString(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

	var specialCharacter = false
	var lastSpecialToken *gotokenize.Token

	for {

		token := iter.Read()
		if token == nil {
			break
		}

		if token.Content == "\\" {

			specialCharacter = !specialCharacter
			lastSpecialToken = token

		} else if token.Content == reach {

			if specialCharacter {

				specialCharacter = false
				currentToken.Children.AddToken(*token)

			} else {

				break
			}
		} else {

			if specialCharacter && token.Content != "{" && token.Content != "}" && token.Content != "\"" && token.Content != "'" {

				currentToken.Children.AddToken(*lastSpecialToken)
			}
			specialCharacter = false
			currentToken.Children.AddToken(*token)
		}
	}
}
