package errors

type AftError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e AftError) Error() string {
	return e.Message
}
