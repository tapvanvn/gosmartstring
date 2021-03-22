package gosmartstring

type SSJSON struct {
	IObject
	attributes map[string]IObject
}

func CreateSSJSON(jsonString string) (SSJSON, error) {
	json := SSJSON{
		IObject: &SSObject{},
	}
	return json, nil
}
