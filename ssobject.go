package gosmartstring

//func(context, input, param) output
type IFunction func(context *sscontext, input IObject, params []IObject) IObject

//IObject interface for ssobject
type IObject interface {
	Parent() IObject
	CanExport() bool
	Export() []byte
	GetType() string
	GetExtendFunc() map[string]IFunction
	Call(context *sscontext, name string, params []IObject) IObject
	Extend(functionName string, sfunc IFunction)
}

//Object ssobject
type SSObject struct {
	parent          IObject
	extendFunctions map[string]IFunction
}

var EmptyObject = &SSObject{
	parent:          nil,
	extendFunctions: nil,
}

//MARK: implement IObject
func (obj *SSObject) Parent() IObject {
	return obj.parent
}

func (obj *SSObject) CanExport() bool {
	return false
}

func (obj *SSObject) Export() []byte {
	return nil
}

func (obj *SSObject) GetType() string {
	return "ssobject"
}

func (obj *SSObject) GetExtendFunc() map[string]IFunction {

	return obj.extendFunctions
}

func (obj *SSObject) Extend(functionName string, sfunc IFunction) {

	obj.extendFunctions[functionName] = sfunc
}

func (obj *SSObject) Call(context *sscontext, name string, params []IObject) IObject {

	if name == "json" {

	}
	if sfunc, ok := obj.extendFunctions[name]; ok {

		return sfunc(context, obj, params)
	}
	return nil
}
