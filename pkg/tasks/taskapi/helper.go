package taskapi

import "fmt"

func asString(value interface{}) (string, error) {
	result, ok := value.(string)
	if ok {
		return result, nil
	}

	return "", fmt.Errorf("not a string: %T (%+v)", value, value)
}

func asListOfStrings(value interface{}) ([]string, error) {
	if v, ok := value.([]string); ok {
		return v, nil
	}

	elements, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a list of strings: type %T (%+v)", value, value)
	}

	listOfStrings := []string{}

	for _, element := range elements {
		str, ok := element.(string)
		if !ok {
			return nil, fmt.Errorf("not a list of strings: invalid element: type %T (%+v)", element, element)
		}
		listOfStrings = append(listOfStrings, str)
	}

	return listOfStrings, nil
}
