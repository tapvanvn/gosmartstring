package gosmartstring

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type SSInt struct {
	IObject
	Value int64
}

func CreateSSInt(value int64) *SSInt {

	ssint := &SSInt{
		//IObject: CreateSSObject(),
		Value: value,
	}
	ssint.IObject = CreateSSObject(ssint)
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
	fmt.Println("call random int")
	if ssint, ok := input.(*SSInt); ok {
		var b [8]byte
		crypto_rand.Read(b[:])
		rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

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
		ssint.Value = rand.Int63n(max)
		fmt.Println("randomed", ssint.Value)
		return ssint
	}
	fmt.Println("error input is not int", input.GetType())
	return nil //TODO: return error
}
