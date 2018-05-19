package server

const(
	DEFAULT_GRPC_PORT = 22332
)

type ProxyArgs struct {
	GrpcPort uint32
	ListenerPort uint32
	MetricsPort uint32
}

type ProxyServer struct {

}

func NewProxy(args *ProxyArgs) (*ProxyServer, error) {
	return &ProxyServer{}, nil
}

func (s *ProxyServer) StartProxy() error {
	return nil
}
