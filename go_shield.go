package goshield

type (
	DataMap           map[string]interface{}
	RulesMap          map[string][]string
	CustomMessagesMap map[string]string
	ValidationErrors  map[string][]string
	ValidationFunc    func(interface{}, ...string) bool
	ValidationOptions struct {
		function ValidationFunc
		message  string
	}

	validator struct {
		data           DataMap
		customMessages CustomMessagesMap
		rules          RulesMap
		errors         ValidationErrors
		allRules       []string
	}

	ValidatorInterface interface {
		IsValid() bool
		IsFailed() bool
		Errors() ValidationErrors
	}
)

var validationHandler = map[string]ValidationOptions{
	"string": {
		function: isString,
		message:  "The :attribute field must be a string",
	},
	"numeric": {
		function: isNumeric,
		message:  "The :attribute field must be a number",
	},
	"alpha": {
		function: isAlpha,
		message:  "The :attribute field must contain only letters",
	},
	"email": {
		function: isEmail,
		message:  "The :attribute field must be a valid email",
	},
	"min": {
		function: isMin,
		message:  "The :attribute field must be a minimum of :min characters/count",
	},
	"max": {
		function: isMax,
		message:  "The :attribute field must be a maximum of :max characters/count",
	},
	"required": {
		function: isRequired,
		message:  "The :attribute field is required",
	},
	"in": {
		function: inEnum,
		message:  "The :attribute field must be in :in",
	},
}

func AddValidationFunc(name string, fn ValidationFunc, msg string) {
	validationHandler[name] = ValidationOptions{
		function: fn,
		message:  msg,
	}
}

func GetAllValidationRules() []string {
	//get only keys from map
	keys := make([]string, 0, len(validationHandler))
	for k := range validationHandler {
		keys = append(keys, k)
	}
	return keys
}

func Validator(data DataMap, rules RulesMap, customMessage ...CustomMessagesMap) (ValidatorInterface, error) {
	cm := make(CustomMessagesMap)
	if len(customMessage) > 0 {
		cm = customMessage[0]
	}
	v := &validator{
		data:           data,
		rules:          rules,
		errors:         make(ValidationErrors),
		allRules:       GetAllValidationRules(),
		customMessages: cm,
	}

	err := v.doValidate()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *validator) IsFailed() bool {
	return !v.IsValid()
}

func (v *validator) Errors() ValidationErrors {
	return v.errors
}
