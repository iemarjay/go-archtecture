package utils

type ErrorBag struct {
	errors map[interface{}]error
}

func NewErrorBag(errors ...map[interface{}]error) *ErrorBag {
	if len(errors) > 0 {
		return &ErrorBag{errors: errors[0]}
	}

	return &ErrorBag{errors: map[interface{}]error{}}
}

func (e *ErrorBag) Add(key interface{}, err error) {
	e.errors[key] = err
}

func (e *ErrorBag) Error() string {
	errorString := ""

	for _, err := range e.errors {
		errorString += err.Error() + "\n"
	}

	return errorString
}
