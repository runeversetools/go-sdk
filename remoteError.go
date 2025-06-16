package rvtools

type RemoteError struct {
	ErrorDescription string
}

func (e *RemoteError) Error() string {
	return "RemoteError: " + e.ErrorDescription
}
