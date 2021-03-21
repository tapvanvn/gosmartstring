package json

import (
	gotokenize "github.com/tapvanvn/gotokenize"
)

type JSONMeaning struct {
	parent *gotokenize.Meaning
}

func CreateJSONMeaning() JSONMeaning {
	jsonMeaning := JSONMeaning{
		parent: gotokenize.M
	}
}
