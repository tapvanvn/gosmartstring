package gosmartstring

type SSBool struct {
	IObject
	Value bool
}

func CreateBool(value bool) SSBool {
	return SSBool{
		IObject: &SSObject{},
		Value:   value,
	}
}

func (obj SSBool) CanExport() bool {
	return true
}

func (obj SSBool) Export() []byte {
	if obj.Value {
		return []byte("true")
	}
	return []byte("false")
}

func (obj SSBool) GetType() string {
	return "ssbool"
}
