package gosmartstring

import (
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

/*
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
			meaning.parseSmartstring(token)
		}
		tmpStream.AddToken(*token)
	}
	process.SetStream(process.Context.AncestorTokens, &tmpStream)
	fmt.Println("--1.1--")
	process.Stream.Debug(0, SSNaming, &gotokenize.DebugOption{
		ExtendTypeSize: 6,
	})
	fmt.Println("--end 1.1--")
}*/
func (meaning *SmarstringMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {
	token := meaning.getNextMeaningToken(process)
	return token
}
func (meaning *SmarstringMeaning) getNextMeaningToken(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := process.Iter.Read()
	if token != nil {

		if token.Type == TokenSSLSmartstring {

			meaning.parseSmartstring(&process.Context, token)
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
		}
	}
}

func (meaning *SmarstringMeaning) parseCommand(context *gotokenize.MeaningContext, parentToken *gotokenize.Token) {

	second := parentToken.Children.GetTokenAt(1)

	meaning.parseParentThese(context, second)
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
				tmpStream.AddToken(*pack.Children.GetTokenAt(0))
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
		tmpStream.AddToken(*pack.Children.GetTokenAt(0))
	} else if pack.Children.Length() > 1 {
		meaning.processChild(context, pack)
		tmpStream.AddToken(*pack)
	}

	parentToken.Children = tmpStream

}

func (meaning *SmarstringMeaning) GetName() string {
	return "ss_meaning"
}
