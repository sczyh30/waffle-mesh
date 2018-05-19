package metrics

type MonitorServer interface {
	Start(stop chan struct{}) error
}

type simpleMetricsServer struct {

}

func (s *simpleMetricsServer) Start(stop chan struct{}) error {
	return nil
}

func NewMetricsServer() (*MonitorServer, error) {
	return nil, nil
}
