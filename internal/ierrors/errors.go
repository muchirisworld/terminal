package errors

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type InsufficientStockError struct {
	Message string
}

func (e *InsufficientStockError) Error() string {
	return e.Message
}
