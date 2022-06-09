package gosmartstring

type SSInterface struct {
	parent    *SSInterface
	functions map[string]IFunction
}

func (inf *SSInterface) Call(context *SSContext, input IObject, name string, params []IObject) IObject {

	if sfunc, ok := inf.functions[name]; ok {

		return sfunc(context, input, params)

	} else if inf.parent != nil {

		return inf.parent.Call(context, input, name, params)
	}
	return nil
}
func (obj *SSInterface) Extend(funcs map[string]IFunction) *SSInterface {
	newInterface := &SSInterface{
		parent:    obj,
		functions: map[string]IFunction{},
	}
	for key, fn := range funcs {
		newInterface.functions[key] = fn
	}
	return newInterface
}
