package integration

const (
	StatusCode = "StatusCode"
)

type Request interface {
}

type Response interface {
	Code() int
}
