package server

type ProxyServer struct {

}

func NewProxy() (*ProxyServer, error) {
	return &ProxyServer{}, nil
}

func (s *ProxyServer) StartProxy() error {
	return nil
}
