package integration

const (
	StatusCode = "StatusCode"
)

// 请求元数据
type Metadata struct {
	ClientId  string `json:"clientId"`
	NameSpace string `json:"namespace"`
	HeaderMap map[string]string
}

func (meta *Metadata) Headers(mp map[string]string) {
	meta.HeaderMap = mp
}

func (meta *Metadata) mustImplementMetadata() {
}

func (meta *Metadata) GetHeader(key, defaultValue string) string {
	if value, ok := meta.HeaderMap[key]; ok {
		return value
	}
	return defaultValue
}

type Request interface {
	Headers(map[string]string)
	mustImplementMetadata()
	GetHeader(key, defaultValue string) string
}

type Response interface {
	Code() int
	SetCode(int)
}
type ResponseData struct {
	StatusCode int    `json:"statusCode"`
	Resp       string `json:"resp"`
}

func (response *ResponseData) Code() int {
	return response.StatusCode
}

func (response *ResponseData) SetCode(code int) {
	response.StatusCode = code
}
