package gosmartstring

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

type SSInt struct {
	IObject
	Value int64
}

func CreateSSInt(value int64) SSInt {

	ssint := SSInt{
		IObject: &SSObject{},
		Value:   value,
	}
	ssint.Extend("random", ssintFuncRandom)
	return ssint
}

func (obj SSInt) CanExport() bool {
	return true
}

func (obj SSInt) Export(context *SSContext) []byte {

	return []byte(strconv.FormatInt(obj.Value, 10))
}

func (obj SSInt) GetType() string {

	return "ssint"
}

func ssintFuncRandom(context *SSContext, input IObject, params []IObject) IObject {
	rand.Seed(int64(time.Now().Nanosecond()))
	var max int64 = math.MaxInt64
	if len(params) > 0 {
		switch params[0].(type) {
		case *SSString:
			if test, err := strconv.ParseInt(params[0].(*SSString).Value, 10, 64); err == nil {
				max = test
			}
		case *SSInt:
			max = params[0].(*SSInt).Value
		}
	}
	return CreateSSInt(rand.Int63n(max))
}
