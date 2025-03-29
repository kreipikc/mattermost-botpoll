package utils

import "fmt"

func ConvertAllIntToInt(tuple interface{}) (int, error) {
	var result int
	switch v := tuple.(type) {
	case int8:
		result = int(v)
	case int16:
		result = int(v)
	case int32:
		result = int(v)
	case int64:
		result = int(v)
	case uint8:
		result = int(v)
	case uint16:
		result = int(v)
	case uint32:
		result = int(v)
	case uint64:
		result = int(v)
	case float64:
		result = int(v)
	default:
		return 0, fmt.Errorf("неожиданный тип для Id: %T", v)
	}
	return result, nil
}
