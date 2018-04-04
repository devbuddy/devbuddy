package tasks

import "errors"

func asString(value interface{}) (string, error) {
	result, ok := value.(string)
	if ok {
		return result, nil
	}

	_, ok = value.(bool)
	if ok {
		return "", errors.New("found a boolean, not a string")
	}

	return "", errors.New("not a string")
}
