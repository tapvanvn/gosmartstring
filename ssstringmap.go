package gosmartstring

import "sync"

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
