package gosmartstring

type SSBool struct {
	parent *SSObject
	Value  bool
}

func CreateBool(value bool) SSBool {
	return SSBool{
		parent: &SSObject{},
		Value:  value,
	}
}

//MARK: implement IObject
func (obj *SSBool) Parent() IObject {
	return obj.parent
}

func (obj *SSBool) CanExport() bool {
	return true
}

func (obj *SSBool) Export() []byte {
	if obj.Value {
		return []byte("true")
	}
	return []byte("false")
}

func (obj *SSBool) GetType() string {
	return "ssbool"
}

func (obj *SSBool) GetExtendFunc() map[string]IFunction {

	return obj.parent.GetExtendFunc()
}

func (obj *SSBool) Extend(functionName string, sfunc IFunction) {

	obj.parent.Extend(functionName, sfunc)
}

func (obj *SSBool) Call(context *SSContext, name string, params []IObject) IObject {

	return obj.parent.Call(context, name, params)
}
