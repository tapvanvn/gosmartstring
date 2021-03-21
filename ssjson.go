package gosmartstring

type SSJSON struct {
	parent     *SSObject
	attributes map[string]IObject
}

func CreateSSJSON(jsonString string) (SSJSON, error) {
	json := SSJSON{}
	return json, nil
}
