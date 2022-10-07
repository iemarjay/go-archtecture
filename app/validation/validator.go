package validation

import "github.com/go-playground/validator/v10"

type Error struct {
	Namespace string `json:"namespace"`
	Field     string `json:"field"`
	Tag       string `json:"tag"`
	Value     string `json:"value"`
}

type ErrorBag struct {
	errors map[string]*Error
}

func NewErrorBag() *ErrorBag {
	return &ErrorBag{errors: map[string]*Error{}}
}

func (e *ErrorBag) Add(err *Error) {
	e.errors[err.Namespace] = err
}

func (e *ErrorBag) Error() string {
	errorString := ""

	for key, err := range e.errors {
		errorString += key + ": " + err.Value + "\n"
	}

	return errorString
}

func (e *ErrorBag) All() map[string]*Error {
	return e.errors
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return NewValidatorWithValidator(validator.New())
}

func NewValidatorWithValidator(validator *validator.Validate) *Validator {
	return &Validator{validator: validator}
}

func (av *Validator) Validate(i interface{}) (bool, error) {
	err := av.validator.Struct(i)
	if err == nil {
		return false, nil
	}

	errors := NewErrorBag()
	for _, err := range err.(validator.ValidationErrors) {
		var element Error
		element.Namespace = err.Namespace()
		element.Field = err.Field()
		element.Tag = err.Tag()
		element.Value = err.Param()

		errors.Add(&element)
	}

	return true, errors
}
