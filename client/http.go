package client

type HttpClient struct {
	url string
}

func (this *HttpClient) QueryService(name string) bool {
	return true
}
