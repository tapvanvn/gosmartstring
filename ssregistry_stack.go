package gosmartstring

type SSAddressStack struct {
	IObject
	Address  []map[string]string
	offset   int
	numStack int
}

func CreateAddressStack() SSAddressStack {

	stack := SSAddressStack{
		IObject:  &SSObject{},
		Address:  make([]map[string]string, 0),
		offset:   0,
		numStack: 1,
	}
	stack.Address = append(stack.Address, make(map[string]string, 0))
	return stack
}
func (stack *SSAddressStack) Inc() {
	stack.numStack++
	stack.Address = append(stack.Address, make(map[string]string, 0))

}

func (stack *SSAddressStack) Append(address string, translateAddress string) {

	stack.Address[stack.offset][address] = translateAddress
}

func (stack *SSAddressStack) Get(address string) (string, bool) {
	translate, ok := stack.Address[stack.offset][address]
	return translate, ok
}

func (stack *SSAddressStack) SetStack(offset int) {
	stack.offset = offset
}
func (stack *SSAddressStack) GetStackNum() int {

	return stack.numStack
}
