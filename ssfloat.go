package gosmartstring

import (
	"fmt"
)

type SSFloat struct {
	IObject
	Value float64
}

func CreateSSFloat(value float64) SSFloat {

	ssfloat := SSFloat{
		IObject: &SSObject{},
		Value:   value,
	}
	return ssfloat
}

func (obj SSFloat) CanExport() bool {
	return true
}

func (obj SSFloat) Export(context *SSContext) []byte {

	return []byte(fmt.Sprintf("%f", obj.Value))
}

func (obj SSFloat) GetType() string {

	return "ssfloat"
}
