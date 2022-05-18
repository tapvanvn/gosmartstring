package gosmartstring

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

type SmarstringMeaning struct {
	*gotokenize.AbstractMeaning
}

func CreateSSPatternMeaning() *gotokenize.PatternMeaning {
	meaning := CreateSSRawMeaning()
	pattern := gotokenize.NewPatternMeaning(meaning,
		buildSSLPatterns(),
		SSLIgnores,
		getSSLGlobalNested(),
		gotokenize.NoTokens)
	return pattern
}

func CreateSSMeaning() *SmarstringMeaning {

	return &SmarstringMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(
			CreateSSPatternMeaning(),
		),
	}
}

func (meaning *SmarstringMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {
	token := meaning.getNextMeaningToken(process)
	// if token != nil {
	// 	fmt.Println("--ssmeaning--")
	// 	token.Debug(0, SSNaming, &gotokenize.DebugOption{ExtendTypeSize: 6})
	// 	fmt.Println("--end ssmeaning--")
	// }
	return token
}
func (meaning *SmarstringMeaning) getNextMeaningToken(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := process.Iter.Read()
	if token != nil {

		if token.Type == TokenSSLSmartstring {

			meaning.parseSmartstring(&process.Context, token)

		} else if token.Type == TokenSSLSquare {

			meaning.parseSquare(&process.Context, token)
		}
	}
	return token
}
func (meaning *SmarstringMeaning) parseSmartstring(context *gotokenize.MeaningContext, token *gotokenize.Token) {

	if token.Type != TokenSSLSmartstring {

		return
	}

	iter := token.Children.Iterator()

	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenSSLCommand {

			meaning.parseCommand(context, childToken)

		} else if token.Type == TokenSSLSquare {

			meaning.parseSquare(context, childToken)
		}
	}
}

func (meaning *SmarstringMeaning) parseCommand(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	second := parentToken.Children.GetTokenAt(1)

	meaning.parseParentThese(context, second)
}
func (meaning *SmarstringMeaning) parseSquare(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	fmt.Println(gotokenize.ColorRed("parsing Square"))
	iter := parentToken.Children.Iterator()
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

	for {
		token := iter.Read()
		if token == nil || token.Type != TokenSSLString {
			break
		}
		next := iter.Get()
		val := iter.GetBy(1)
		if next == nil || val == nil || next.Content != ":" {
			break
		}

		pair := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLPair, "")

		pair.Content = token.Children.ConcatStringContent()
		iter.Read()
		if val.Type == TokenSSLString {
			valToken := iter.Read()
			pair.AddChild(*valToken)
		} else {
			pack := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLSmartstring, "")
			for {
				valToken := iter.Read()
				if valToken == nil || valToken.Content == "," {

					break
				}
				pack.AddChild(*valToken)
			}
			meaning.processChild(context, pack)
			pair.AddChild(*pack)
		}

		tmpStream.AddToken(*pair)
	}

	parentToken.Children = tmpStream
}
func (meaning *SmarstringMeaning) processChild(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {
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
func (meaning *SmarstringMeaning) parseParentThese(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	iter := parentToken.Children.Iterator()
	tmpStream := gotokenize.CreateStream(meaning.GetMeaningLevel())
	pack := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLSmartstring, "")
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Content == "," {

			if pack.Children.Length() == 1 {

				first := pack.Children.GetTokenAt(0)
				if first.Type == TokenSSLSquare {
					meaning.parseSquare(context, first)
				} else if first.Type == TokenSSLCommand {
					meaning.parseCommand(context, first)
				}
				tmpStream.AddToken(*first)

			} else if pack.Children.Length() > 1 {

				meaning.processChild(context, pack)

				tmpStream.AddToken(*pack)
			}
			pack = gotokenize.NewToken(meaning.GetMeaningLevel(), TokenSSLSmartstring, "")
			continue
		}
		pack.Children.AddToken(*childToken)
	}
	if pack.Children.Length() == 1 {
		first := pack.Children.GetTokenAt(0)

		if first.Type == TokenSSLSquare {
			meaning.parseSquare(context, first)
		} else if first.Type == TokenSSLCommand {
			meaning.parseCommand(context, first)
		}
		tmpStream.AddToken(*first)
	} else if pack.Children.Length() > 1 {
		meaning.processChild(context, pack)
		tmpStream.AddToken(*pack)
	}

	parentToken.Children = tmpStream
	// fmt.Println("**parenthese**")
	// parentToken.Debug(0, SSNaming, &gotokenize.DebugOption{ExtendTypeSize: 6})
	// fmt.Println("**end parenthese**")
}

func (meaning *SmarstringMeaning) GetName() string {
	return "ss_meaning"
}
