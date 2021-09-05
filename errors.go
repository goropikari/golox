package tlps

var (
	ParseError   = NewCustomError("ParseError")
	RuntimeError = NewCustomError("RuntimeError")
)

type CustomError struct {
	typ     string
	Token   *Token
	message string
}

func (e *CustomError) Error() string {
	return e.typ + ": " + e.message
}

func (e *CustomError) New(token *Token, message string) error {
	return &CustomError{typ: e.typ, Token: token, message: message}
}

func NewCustomError(typ string) *CustomError {
	return &CustomError{typ: typ}
}
