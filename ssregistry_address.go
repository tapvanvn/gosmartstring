package gosmartstring

func NewSSRegistryAddress(address string) *SSRegistryAddress {
	return &SSRegistryAddress{
		IObject: &SSObject{},
		Address: address,
	}
}

type SSRegistryAddress struct {
	IObject
	Address string
}
