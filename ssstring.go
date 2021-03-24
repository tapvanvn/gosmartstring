package gosmartstring

type SSString struct {
	IObject
	Value string
}

func CreateString(value string) *SSString {
	return &SSString{
		IObject: &SSObject{},
		Value:   value,
	}
}

func (obj SSString) CanExport() bool {
	return true
}

func (obj SSString) Export(context *SSContext) []byte {

	return []byte(obj.Value)
}

func (obj SSString) GetType() string {
	return "ssstring"
}
