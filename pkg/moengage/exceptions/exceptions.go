package exceptions

type CreateClientException struct {
	Msg string // description of error
}

func (e *CreateClientException) Error() string { return e.Msg }
