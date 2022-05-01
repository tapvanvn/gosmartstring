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

func (obj *SSStringMap) Set(key string, val IObject) {
	obj.Lock()
	defer obj.Unlock()
	obj.attributes[key] = val

}

func (obj *SSStringMap) Call(context *SSContext, name string, params []IObject) IObject {
	obj.Lock()
	defer obj.Unlock()
	if iobj, ok := obj.attributes[name]; ok {

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
