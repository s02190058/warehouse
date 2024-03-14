package response

const (
	statusOK    = "ok"
	statusError = "error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func OK() Response {
	return Response{
		Status: statusOK,
	}
}

func Error(err error) Response {
	return Response{
		Status:  statusError,
		Message: err.Error(),
	}
}
