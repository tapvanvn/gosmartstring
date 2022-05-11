package gosmartstring

import (
	"fmt"
	"sync"

	"github.com/tapvanvn/gotokenize/v2"
)

func CreateSSStringMap() *SSStringMap {

	mapObj := &SSStringMap{
		IObject:    &SSObject{},
		attributes: map[string]IObject{},
	}

	return mapObj
}

type SSStringMap struct {
	IObject
	sync.Mutex
	attributes map[string]IObject
}

func (obj *SSStringMap) Get(key string) IObject {
	obj.Lock()
	defer obj.Unlock()
	if iobj, ok := obj.attributes[key]; ok {

		return iobj
	}
	return nil
}

func (obj *SSStringMap) Set(key string, val IObject) {
	fmt.Println("set:", gotokenize.ColorRed(key), "value:", getObjectPresent(val))
	obj.Lock()
	defer obj.Unlock()
	obj.attributes[key] = val
}

func (obj *SSStringMap) Call(context *SSContext, name string, params []IObject) IObject {
	attr := name
	if name == "get" {
		if len(params) == 1 {
			if str, ok := params[0].(*SSString); ok {
				attr = str.Value
			}
		}
	}
	obj.Lock()
	defer obj.Unlock()

	if iobj, ok := obj.attributes[attr]; ok {

		if attr == "category_uuid" {
			if iobj != nil {
				fmt.Println("category _uuid", getObjectPresent(iobj))
			} else {
				fmt.Println("category _uuid nil")
			}
		}
		return iobj
	}

	return obj.IObject.Call(context, name, params)
}

func (obj *SSStringMap) GetType() string {
	return "ssstringmap"
}

func (obj *SSStringMap) CanExport() bool {
	return false
}
func getObjectPresent(obj IObject) string {
	if obj != nil {
		if str, ok := obj.(*SSString); ok {
			return "[str]" + str.Value
		} else {
			return obj.GetType()
		}
	} else {
		return "nil"
	}
}
func (obj *SSStringMap) PrintDebug(level int) {

	indent := "|"
	for i := 1; i < level; i++ {
		indent += " |"
	}
	fmt.Printf("%s%s\n", indent, gotokenize.ColorTeal("sstringmap"))
	for key, val := range obj.attributes {

		fmt.Printf("%s%s:%s\n", indent, gotokenize.ColorPurple(key), getObjectPresent(val))

	}
}
