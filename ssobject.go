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
	//Extend(functionName string, sfunc IFunction)
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

func ssFuncExportJson(context *SSContext, input IObject, params []IObject) IObject {
	if input.CanExport() {
		content := string(input.Export(context))
		return CreateSSJSON(content)
	}
	return nil
}

var SSObjectInterface = &SSInterface{
	functions: map[string]IFunction{
		"json": ssFuncExportJson,
	},
}

func CreateSSObject(baseObject IObject, objInterface *SSInterface) *SSObject {
	obj := &SSObject{
		baseObject:   baseObject,
		objInterface: objInterface,
	}
	if objInterface == nil {
		obj.objInterface = SSObjectInterface
	}
	return obj
}

//Object ssobject
type SSObject struct {
	baseObject   IObject
	objInterface *SSInterface
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

	return obj.objInterface.functions
}

func (obj *SSObject) Call(context *SSContext, name string, params []IObject) IObject {

	if name == "json" {
		if obj.CanExport() {
			content := string(obj.Export(context))
			return CreateSSJSON(content)
		}
	}
	if obj.objInterface != nil {

		if obj.baseObject != nil {

			return obj.objInterface.Call(context, obj.baseObject, name, params)
		}
		return obj.objInterface.Call(context, obj, name, params)
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
