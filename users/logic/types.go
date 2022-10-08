package logic

type validator interface {
	Validate(i interface{}) (bool, error)
}
