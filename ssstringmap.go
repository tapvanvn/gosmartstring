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
	//fmt.Println("set:", gotokenize.ColorRed(key), "value:", getObjectPresent(val))
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

type SSStringMapIterator struct {
	keys    []string
	current int
}

func (iter *SSStringMapIterator) IsEnd() bool {
	return iter.current >= len(iter.keys)
}
func (obj *SSStringMap) Iterator() IIterator {
	iter := &SSStringMapIterator{
		current: 0,
		keys:    make([]string, 0),
	}
	for key, _ := range obj.attributes {
		iter.keys = append(iter.keys, key)
	}
	return iter
}

func (obj *SSStringMap) Iterate(context *SSContext, iter IFunctionIterate, iterator IIterator, data interface{}) error {
	if iterator.IsEnd() {
		return nil
	}
	if ssIter, ok := iterator.(*SSStringMapIterator); ok {
		attName := ssIter.keys[ssIter.current]
		if err := iter(context, CreateString(attName), obj.attributes[attName], data); err != nil {
			return err
		}
		ssIter.current++
	}
	return nil
}
