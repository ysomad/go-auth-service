package httperror

type httpError struct {
	st     string
	status int
	msg    string
	err    error
}

func New(st string, status int, msg string, err error) (httpError, error) {
	return httpError{
		st:     st,
		status: status,
		msg:    msg,
		err:    err,
	}, nil
}

func (e httpError) ST() string      { return e.st }
func (e httpError) Status() int     { return e.status }
func (e httpError) Message() string { return e.msg }
func (e httpError) Err() error      { return e.err }
