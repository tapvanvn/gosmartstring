package gosmartstring

type SSJSON struct {
	SSObject
	attributes map[string]IObject
}

func CreateSSJSON(jsonString string) (SSJSON, error) {
	json := SSJSON{
		SSObject: SSObject{},
	}
	return json, nil
}
