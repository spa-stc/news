package web

// Web error.
type Error struct {
	// Human readable, public error message (defaults to "something went wrong").
	Message string
	// Http response code.
	Code int

	// Root error (debug mode + logging only).
	Root error
}

func (e Error) Error() string {
	return e.Message
}

func RespondError(message string, code int, root error) error {
	return Error{
		Message: message,
		Root:    root,
		Code:    code,
	}
}
