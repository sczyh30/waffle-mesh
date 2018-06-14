package pkg

type BazResponse struct {
	Version string `json:"version"`
	Host string `json:"host"`
}

type FooResponse struct {
	Message string `json:"message"`
}
