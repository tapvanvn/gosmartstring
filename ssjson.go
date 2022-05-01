package gosmartstring

import (
	"sync"

	"github.com/tapvanvn/gotokenize/v2"
	"github.com/tapvanvn/gotokenize/v2/json"

	jsonEnc "encoding/json"
)

type SSJSON struct {
	IObject
	sync.Mutex
	attributes map[string]IObject
}

func CreateSSJSON(jsonString string) *SSJSON {

	jsonObj := &SSJSON{
		IObject:    &SSObject{},
		attributes: map[string]IObject{},
	}

	if jsonString != "" {

		stream := gotokenize.CreateStream(0)
		stream.Tokenize(jsonString)
		meaning := json.CreateJSONMeaning()

		proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
		meaning.Prepare(proc)

		for {
			token := meaning.Next(proc)
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

func ParseInterface(object interface{}) IObject {
	if data, err := jsonEnc.Marshal(object); err == nil {
		return ParseJSONString(string(data))
	}
	return nil
}

func ParseJSONString(jsonString string) IObject {
	stack := []IObject{}
	if jsonString != "" {

		stream := gotokenize.CreateStream(0)
		stream.Tokenize(jsonString)
		meaning := json.CreateJSONMeaning()

		proc := gotokenize.NewMeaningProcessFromStream(gotokenize.NoTokens, &stream)
		meaning.Prepare(proc)

		for {
			token := meaning.Next(proc)
			if token == nil {
				break
			}

			if token.Type == json.TokenJSONPair {

				jsonObj := CreateSSJSON("")
				parsePair(token, jsonObj)
				stack = append(stack, jsonObj)

			} else if token.Type == json.TokenJSONBlock {
				jsonObj := CreateSSJSON("")
				parseBlock(token, jsonObj)
				stack = append(stack, jsonObj)
			} else if token.Type == json.TokenJSONSquare {
				array := CreateSSArray()
				parseArray(token, array)
				stack = append(stack, array)
			}
		}
		numElement := len(stack)
		if numElement == 1 {
			return stack[0]
		} else if numElement > 1 {
			array := CreateSSArray()
			for _, element := range stack {
				array.Stack = append(array.Stack, element)
			}
			return array
		}
	}
	return nil
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
		object.Lock()
		object.attributes[attribute] = valueObject
		object.Unlock()

	} else if value.Type == json.TokenJSONSquare {

		valueArray := CreateSSArray()
		parseArray(value, valueArray)
		object.Lock()
		object.attributes[attribute] = valueArray
		object.Unlock()

	} else if value.Type == json.TokenJSONBlock {

		valueObject := CreateSSJSON("")
		object.Lock()
		object.attributes[attribute] = valueObject
		object.Unlock()
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

func (obj *SSJSON) Call(context *SSContext, name string, params []IObject) IObject {
	obj.Lock()
	defer obj.Unlock()
	if iobj, ok := obj.attributes[name]; ok {

		return iobj
	}

	return obj.IObject.Call(context, name, params)
}

func (obj *SSJSON) GetType() string {
	return "ssjson"
}

func (obj *SSJSON) CanExport() bool {
	return false
}
