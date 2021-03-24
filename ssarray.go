package gosmartstring

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type SSArray struct {
	IObject
	Stack []IObject
}

func CreateSSArray() *SSArray {
	return &SSArray{
		IObject: &SSObject{},
		Stack:   []IObject{},
	}
}

//MARK: implement IObject

func (obj SSArray) GetType() string {
	return "ssarray"
}

func (obj SSArray) Call(context *SSContext, name string, params []IObject) IObject {

	if name == "add" {

		obj.add(params)
		return obj
	} else if name == "random" {

		return obj.random()
	}
	return obj.IObject.Call(context, name, params)
}
func (obj *SSArray) add(params []IObject) {

	obj.Stack = append(obj.Stack, params...)
}

func (obj *SSArray) random() IObject {
	size := len(obj.Stack)
	if size == 0 {
		return nil
	}
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	offset := rand.Intn(size)
	return obj.Stack[offset]
}
