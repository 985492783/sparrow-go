package integration

const (
	StatusCode = "StatusCode"
)

type Metadata struct {
	ClientId  string `json:"clientId"`
	NameSpace string `json:"namespace"`
}

type Request interface {
}

type Response interface {
	Code() int
}
