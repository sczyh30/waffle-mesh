package bootstrap

const (
	DEFAULT_GRPC_SERVER_PORT = ":24242"
)

type BrainServer struct {

}

func NewServer() (*BrainServer, error) {
	server := &BrainServer{}
	return server, nil
}

func (s *BrainServer) Start(stop chan struct{}) error {
	return nil
}