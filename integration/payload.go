package integration

const (
	StatusCode = "StatusCode"
)

type Request interface {
}

type Response interface {
	Code() int
}

type SuccessRequest struct {
	Name string
}

type SuccessResponse struct {
}

func (success *SuccessResponse) Code() int {
	return 200
}
