package gosmartstring

type ssregistry struct {
	Object   IObject
	Function IFunction
}

func CreateObjectRegistry(object IObject) ssregistry {
	return ssregistry{
		Object:   object,
		Function: nil,
	}
}

func CreateFunctionRegistry(sfunc IFunction) ssregistry {
	return ssregistry{
		Object:   nil,
		Function: sfunc,
	}
}
