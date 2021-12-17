package httperror

type Error interface {
	// ST represents chain of actions string
	// until the error occurs (kind of stack trace).
	ST() string

	// Status determines what status should have
	// been returned to client.
	Status() int

	// Message is occurred error short description.
	Message() string

	// Err is a occured error.
	Err() error
}
