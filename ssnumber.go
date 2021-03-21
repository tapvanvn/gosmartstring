package gosmartstring

type SSNumber struct {
	parent *SSObject
}

func CreateSSNumber(value int) SSNumber {

	return SSNumber{
		parent: &SSObject{},
	}
}

//MARK: implement IObject
func (obj *SSNumber) Parent() IObject {
	return obj.parent
}

func (obj *SSNumber) CanExport() bool {
	return true
}

func (obj *SSNumber) Export() []byte {
	/*if obj.Value {
		return []byte("true")
	}*/
	return []byte("false")
}

func (obj *SSNumber) GetType() string {
	return "ssnumber"
}

func (obj *SSNumber) GetExtendFunc() map[string]IFunction {

	return obj.parent.GetExtendFunc()
}

func (obj *SSNumber) Extend(functionName string, sfunc IFunction) {

	obj.parent.Extend(functionName, sfunc)
}

func (obj *SSNumber) Call(context Context, name string, params []IObject) IObject {

	return obj.parent.Call(context, name, params)
}
