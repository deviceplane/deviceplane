package validation

import "fmt"

func ValidateString(elem interface{}) error {
	switch elem.(type) {
	case string:
		return nil
	default:
		return fmt.Errorf("expected type string")
	}
}

func ValidateInteger(elem interface{}) error {
	switch elem.(type) {
	case int:
		return nil
	default:
		return fmt.Errorf("expected type integer")
	}
}

func ValidateBoolean(elem interface{}) error {
	switch elem.(type) {
	case bool:
		return nil
	default:
		return fmt.Errorf("expected type boolean")
	}
}

func ValidateStringOrInteger(elem interface{}) error {
	switch elem.(type) {
	case string, int:
		return nil
	default:
		return fmt.Errorf("expected type string or integer")
	}
}

func ValidateStringArray(elem interface{}) error {
	switch typedElem := elem.(type) {
	case []interface{}:
		return validateElementsAreStrings(typedElem)
	default:
		return fmt.Errorf("expected type array of strings")
	}
}

func ValidateStringIntegerArray(elem interface{}) error {
	switch typedElem := elem.(type) {
	case []interface{}:
		return validateElementsAreStringsOrIntegers(typedElem)
	default:
		return fmt.Errorf("expected type array of strings or integers")
	}
}

func ValidateStringOrStringArray(elem interface{}) error {
	switch typedElem := elem.(type) {
	case string:
		return nil
	case []interface{}:
		return validateElementsAreStrings(typedElem)
	default:
		return fmt.Errorf("expected type string or array of strings")
	}
}

func ValidateArrayOrObject(elem interface{}) error {
	switch typedElem := elem.(type) {
	case []interface{}:
		return validateElementsAreStrings(typedElem)
	case map[interface{}]interface{}:
		// TODO
		return nil
	default:
		return fmt.Errorf("expected type string or array of strings")
	}
}

func validateElementsAreStrings(elems []interface{}) error {
	for _, elem := range elems {
		switch elem.(type) {
		case string:
			continue
		default:
			return fmt.Errorf("expected type string")
		}
	}
	return nil
}

func validateElementsAreStringsOrIntegers(elems []interface{}) error {
	for _, elem := range elems {
		switch elem.(type) {
		case string, int:
			continue
		default:
			return fmt.Errorf("expected type string or integer")
		}
	}
	return nil
}
