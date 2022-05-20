package gosmartstring

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
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

func (obj *SSArray) GetType() string {
	return "ssarray"
}

func (obj *SSArray) Call(context *SSContext, name string, params []IObject) IObject {

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

func (obj *SSArray) ToString() string {

	return fmt.Sprintf("%v", obj.Stack)
}

type SSArrayIterator struct {
	current int
	size    int
}

func (iter *SSArrayIterator) IsEnd() bool {
	return iter.current >= iter.size
}
func (obj *SSArray) Iterator() IIterator {
	return &SSArrayIterator{
		current: 0,
		size:    len(obj.Stack),
	}
}

func (obj *SSArray) Iterate(context *SSContext, iter IFunctionIterate, iterator IIterator, data interface{}) error {
	if iterator.IsEnd() {
		return nil
	}
	if ssIter, ok := iterator.(*SSArrayIterator); ok && ssIter.current < len(obj.Stack) {

		if err := iter(context, CreateSSInt(int64(ssIter.current)), obj.Stack[ssIter.current], data); err != nil {
			return err
		}
		ssIter.current++
	}
	return nil
}
