package goshield

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func _requiredAmountOfParams(amount int, params []string) {
	if len(params) != amount {
		panic("The rule must have " + strconv.Itoa(amount) + " parameter(s)")
	}
}

func isString(value interface{}, options ...string) bool {
	_, ok := value.(string)
	return ok
}

func isNumeric(value interface{}, options ...string) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func inEnum(value interface{}, options ...string) bool {
	str, ok := value.(string)
	if !ok {
		return false // Value is not a string
	}
	enum := strings.Split(str, ",")
	return slices.Contains(enum, "test")

}

func isAlpha(value interface{}, option ...string) bool {
	str, ok := value.(string)
	if !ok {
		return false // Value is not a string
	}
	alphaRegex := regexp.MustCompile("^[a-zA-Z]+$")
	return alphaRegex.MatchString(str)
}

func isEmail(value interface{}, option ...string) bool {
	str, ok := value.(string)
	if !ok {
		return false // Value is not a string
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(str)
}

func isMin(s interface{}, option ...string) bool {
	_requiredAmountOfParams(1, option)

	i, err := strconv.Atoi(option[0])
	if err != nil {
		panic("Min rule must have at least one parameter")
	}

	switch s.(type) {
	case int:
		return s.(int) > i
	case float64:
		return s.(float64) > float64(i)
	case float32:
		return s.(float32) > float32(i)
	default:
		return len(fmt.Sprint(s)) > i
	}
}

func isMax(s interface{}, option ...string) bool {
	_requiredAmountOfParams(1, option)

	i, err := strconv.Atoi(option[0])
	if err != nil {
		panic("Min rule must have at least one parameter")
	}
	switch s.(type) {
	case int:
		return s.(int) < i
	case float64:
		return s.(float64) < float64(i)
	case float32:
		return s.(float32) < float32(i)
	default:
		return len(fmt.Sprint(s)) < i
	}
}

func isRequired(s interface{}, option ...string) bool {
	return len(fmt.Sprint(s)) > 0
}
