package goshield

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

func (v *validator) doValidate() error {
	for key, rulesArr := range v.rules {
		err := v.validateKey(key, rulesArr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *validator) validateKey(key string, rulesArr []string) error {

	if slices.Contains(rulesArr, "required") {
		if _, ok := v.data[key]; !ok {

			customKey := fmt.Sprintf("%s.%s", key, "required")
			cMsg, ok := v.customMessages[customKey]
			if ok {
				v.errors[key] = append(v.errors[key], cMsg)
				return nil
			}
			r := strings.Replace(validationHandler["required"].message, ":attribute", key, -1)
			v.errors[key] = append(v.errors[key], r)

			return nil
		}
	}

	for _, ruleStr := range rulesArr {
		if ruleStr == "required" {
			continue
		}

		if ruleStr == "" {
			return fmt.Errorf("rule cannot be empty for key %s", key)
		}

		r := strings.Split(ruleStr, ":")
		rule := strings.ToLower(r[0])
		opts := make([]string, 0)
		if len(r) > 1 {
			opts = strings.Split(r[1], ",")
		}
		//
		//if !slices.Contains(v.allRules, rule) {
		//	return fmt.Errorf("rule %s does not exist", rule)
		//}
		validator, ok := validationHandler[rule]

		if !ok {
			return fmt.Errorf("validator for rule %s does not exist", rule)
		}

		if !validator.function(v.data[key], opts...) {

			customKey := fmt.Sprintf("%s.%s", key, rule)
			cMsg, ok := v.customMessages[customKey]
			if ok {
				v.errors[key] = append(v.errors[key], cMsg)
			} else {
				r := formatError(validator.message, key, opts...)
				v.errors[key] = append(v.errors[key], r)
			}
		}
	}
	return nil
}

func formatError(message string, key string, opts ...string) string {
	replaceMap := make(map[string]string)
	matches := findAllMatches(message)
	index := 0
	for _, match := range matches {
		if match == ":attribute" {
			replaceMap[":attribute"] = key
			continue
		}
		replaceMap[match] = opts[index]
		index++
	}
	return replaceMultipleRegex(message, replaceMap)
}

func replaceMultipleRegex(message string, replaceMap map[string]string) string {
	for k, v := range replaceMap {
		message = strings.Replace(message, k, v, -1)
	}
	return message
}

func findAllMatches(message string) []string {
	regex := regexp.MustCompile(`:[a-zA-Z]+`)
	return regex.FindAllString(message, -1)
}
