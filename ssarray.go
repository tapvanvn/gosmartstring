package gosmartstring

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type SSArray struct {
	parent *SSObject
	stack  []IObject
}

func CreateSSArray() SSArray {
	return SSArray{
		parent: &SSObject{},
		stack:  []IObject{},
	}
}

//MARK: implement IObject
func (obj *SSArray) Parent() IObject {
	return obj.parent
}

func (obj *SSArray) CanExport() bool {
	return false
}

func (obj *SSArray) Export() []byte {
	return nil
}

func (obj *SSArray) GetType() string {
	return "ssarray"
}

func (obj *SSArray) GetExtendFunc() map[string]IFunction {

	return obj.parent.GetExtendFunc()
}

func (obj *SSArray) Extend(functionName string, sfunc IFunction) {

	obj.parent.Extend(functionName, sfunc)
}

func (obj *SSArray) Call(context *SSContext, name string, params []IObject) IObject {

	if name == "add" {

		obj.add(params)
		return obj
	} else if name == "random" {

		return obj.random()
	}
	return obj.parent.Call(context, name, params)
}
func (obj *SSArray) add(params []IObject) {

	obj.stack = append(obj.stack, params...)
}

func (obj *SSArray) random() IObject {
	size := len(obj.stack)
	if size == 0 {
		return EmptyObject
	}
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	offset := rand.Intn(size)
	return obj.stack[offset]
}
