package gosmartstring

type Registry struct {
	Object   IObject
	Function IFunction
}

func CreateObjectRegistry(object IObject) Registry {
	return Registry{
		Object:   object,
		Function: nil,
	}
}

func CreateFunctionRegistry(sfunc IFunction) Registry {
	return Registry{
		Object:   nil,
		Function: sfunc,
	}
}
