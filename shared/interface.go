package shared

type Result struct {
	Header       map[string][]string
	URL          string
	Status       int32
	ErrorMessage string
}

type Transport interface {
	Transport(header map[string][]string, URL string, config map[string]string) Result
}

type ResponseResult struct {
	Header map[string][]string
	Body   []byte
}

type Response interface {
	Modify(URL string, method string, header map[string][]string, body []byte) ResponseResult
}
