package gosmartstring

type SSNumber struct {
	SSObject
}

func CreateSSNumber(value int) SSNumber {

	return SSNumber{
		SSObject: SSObject{},
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
