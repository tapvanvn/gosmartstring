package gosmartstring

//func(context, input, param) output
type IFunction func(context *SSContext, input IObject, params []IObject) IObject
type IFunctionIterate func(context *SSContext, key IObject, val IObject, data interface{}) error

//IObject interface for ssobject
type IObject interface {
	CanExport() bool
	Export(context *SSContext) []byte
	GetType() string
	GetExtendFunc() map[string]IFunction
	Call(context *SSContext, name string, params []IObject) IObject
	Extend(functionName string, sfunc IFunction)
	IsTrue() bool
	ToString() string
	PrintDebug(level int)
}
type IIterator interface {
	IsEnd() bool
}

type IIterable interface {
	Iterator() IIterator //create iterator
	Iterate(context *SSContext, iterFunction IFunctionIterate, iterator IIterator, data interface{}) error
}

func CreateSSObject(baseObject IObject) *SSObject {
	return &SSObject{
		baseObject:      baseObject,
		extendFunctions: map[string]IFunction{},
	}
}

//Object ssobject
type SSObject struct {
	baseObject      IObject
	extendFunctions map[string]IFunction
}

//MARK: implement IObject

func (obj *SSObject) CanExport() bool {
	return false
}

func (obj *SSObject) Export(context *SSContext) []byte {
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

func (obj *SSObject) Call(context *SSContext, name string, params []IObject) IObject {

	if name == "json" {
		if obj.CanExport() {
			content := string(obj.Export(context))
			return CreateSSJSON(content)
		}
	}
	if sfunc, ok := obj.extendFunctions[name]; ok {

		if obj.baseObject != nil {

			return sfunc(context, obj.baseObject, params)
		}
		return sfunc(context, obj, params)
	}
	return nil
}

func (obj *SSObject) PrintDebug(level int) {

}

func (obj *SSObject) ToString() string {

	return ""
}
func (obj *SSObject) IsTrue() bool {
	return obj != nil
}
