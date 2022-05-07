package gosmartstring

import (
	"github.com/tapvanvn/gotokenize/v2"
)

const ()

type SmarstringRawMeaning struct {
	*gotokenize.AbstractMeaning
	IsSmartstringOnly bool
}

func CreateSSRawMeaning() *SmarstringRawMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{
		"?:><!&^%@.(){}[]+-*/\\\"'": {TokenType: TokenSSLOperator, Separate: true},
		" ":                         {TokenType: gotokenize.TokenSpace, Separate: true},
		",\r\n":                     {TokenType: TokenSSLBreak, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)
	return &SmarstringRawMeaning{
		AbstractMeaning:   gotokenize.NewAbtractMeaning(meaning),
		IsSmartstringOnly: false,
	}
}

func (meaning *SmarstringRawMeaning) Prepare(proc *gotokenize.MeaningProcess) {

	meaning.AbstractMeaning.Prepare(proc)

	// fmt.Println("--0--")
	// proc.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
	// 	ExtendTypeSize: 6,
	// })
	// fmt.Println("--end 0--")

	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	meaning.IsSmartstringOnly = true

	iter := proc.Iter
	for {
		if iter.EOS() {
			break
		}
		token := iter.Get()

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
	proc.SetStream(proc.Context.AncestorTokens, &tmpStream)

	// fmt.Println("--0.1--")
	// proc.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
	// 	ExtendTypeSize: 6,
	// })
	// fmt.Println("--end 0.1--")
}

func (meaning *SmarstringRawMeaning) continueSmartString(iter *gotokenize.Iterator) gotokenize.Token {

	rsToken := gotokenize.Token{
		Type: TokenSSLSmartstring,
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
		if token.Type == 0 && token.Content != "" {
			token.Type = TokenSSLWord
		}
		_ = iter.Read()
		if token.Type != 0 || token.Content != "" {
			rsToken.Children.AddToken(*token)
		}
	}
	return rsToken
}

func (meaning *SmarstringRawMeaning) meaningSmartString(token gotokenize.Token) gotokenize.Token {

	iter := token.Children.Iterator()

	newToken := gotokenize.Token{
		Type: TokenSSLSmartstring,
	}

	meaningToken := meaning.getNextMeaningToken(iter)

	for {
		if meaningToken == nil {
			break
		}
		newToken.Children.AddToken(*meaningToken)

		meaningToken = meaning.getNextMeaningToken(iter)
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
func (meaning *SmarstringRawMeaning) GetName() string {
	return "ss_raw"
}
