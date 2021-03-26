package gosmartstring

import (
	"github.com/tapvanvn/gotokenize"
	"github.com/tapvanvn/gotokenize/json"
)

type SSJSON struct {
	IObject
	attributes map[string]IObject
}

func CreateSSJSON(jsonString string) *SSJSON {

	jsonObj := &SSJSON{
		IObject:    &SSObject{},
		attributes: map[string]IObject{},
	}

	if jsonString != "" {

		stream := gotokenize.CreateStream()
		stream.Tokenize(jsonString)
		meaning := json.CreateJSONMeaning()

		meaning.Prepare(&stream)

		for {
			token := meaning.Next()
			if token == nil {
				break
			}

			if token.Type == json.TokenJSONPair {

				parsePair(token, jsonObj)

			} else if token.Type == json.TokenJSONBlock {

				parseBlock(token, jsonObj)
				break
			}
		}
	}
	return jsonObj
}

func parseBlock(token *gotokenize.Token, object *SSJSON) {

	iter := token.Children.Iterator()

	for {

		child := iter.Read()
		if child == nil {
			break
		}

		if child.Type == json.TokenJSONPair {

			parsePair(child, object)
		}
	}
}
func parsePair(token *gotokenize.Token, object *SSJSON) {

	iter := token.Children.Iterator()

	key := iter.Get()
	value := iter.GetBy(1)
	if key == nil || value == nil {
		return
	}

	attribute := key.Children.ConcatStringContent()

	if value.Type == json.TokenJSONString || value.Type == json.TokenJSONNumberString {

		valueObject := CreateString(value.Children.ConcatStringContent())

		object.attributes[attribute] = valueObject

	} else if value.Type == json.TokenJSONSquare {

		valueArray := CreateSSArray()
		parseArray(value, valueArray)
		object.attributes[attribute] = valueArray

	} else if value.Type == json.TokenJSONBlock {

		valueObject := CreateSSJSON("")
		object.attributes[attribute] = valueObject
		parseBlock(value, valueObject)
	}

}

func parseArray(token *gotokenize.Token, object *SSArray) {

	iter := token.Children.Iterator()

	for {

		child := iter.Read()
		if child == nil {
			break
		}

		if child.Type == json.TokenJSONString || child.Type == json.TokenJSONNumberString {

			object.Stack = append(object.Stack, CreateString(child.Children.ConcatStringContent()))

		} else if child.Type == json.TokenJSONBlock {

			element := CreateSSJSON("")
			parseBlock(child, element)
			object.Stack = append(object.Stack, element)

		} else if child.Type == json.TokenJSONSquare {

			element := CreateSSArray()
			parseArray(child, element)
			object.Stack = append(object.Stack, element)
		}
	}
}

func (obj SSJSON) Call(context *SSContext, name string, params []IObject) IObject {

	if iobj, ok := obj.attributes[name]; ok {

		return iobj
	}

	return obj.IObject.Call(context, name, params)
}
