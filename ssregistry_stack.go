package gosmartstring

type SSAddressStack struct {
	IObject
	Address []string
}

func CreateAddressStack() SSAddressStack {
	return SSAddressStack{
		IObject: &SSObject{},
		Address: make([]string, 0),
	}
}

func (stack *SSAddressStack) Append(address string) {

	stack.Address = append(stack.Address, address)
}
