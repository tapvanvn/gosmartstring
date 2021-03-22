package gosmartstring

type SSNumber struct {
	IObject
}

func CreateSSNumber(value int) SSNumber {

	return SSNumber{
		IObject: &SSObject{},
	}
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
