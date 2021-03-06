package gosmartstring

//func(context, input, param) output
type IFunction func(context *SSContext, input IObject, params []IObject) IObject

//IObject interface for ssobject
type IObject interface {
	CanExport() bool
	Export(context *SSContext) []byte
	GetType() string
	GetExtendFunc() map[string]IFunction
	Call(context *SSContext, name string, params []IObject) IObject
	Extend(functionName string, sfunc IFunction)
	PrintDebug(level int)
}

//Object ssobject
type SSObject struct {
	extendFunctions map[string]IFunction
}

//MARK: implement IObject

func (obj SSObject) CanExport() bool {
	return false
}

func (obj SSObject) Export(context *SSContext) []byte {
	return nil
}

func (obj SSObject) GetType() string {
	return "ssobject"
}

func (obj SSObject) GetExtendFunc() map[string]IFunction {

	return obj.extendFunctions
}

func (obj SSObject) Extend(functionName string, sfunc IFunction) {

	obj.extendFunctions[functionName] = sfunc
}

func (obj SSObject) Call(context *SSContext, name string, params []IObject) IObject {

	if name == "json" {
		if obj.CanExport() {
			content := string(obj.Export(context))
			return CreateSSJSON(content)
		}
	}
	if sfunc, ok := obj.extendFunctions[name]; ok {

		return sfunc(context, &obj, params)
	}
	return nil
}

func (obj SSObject) PrintDebug(level int) {

}
