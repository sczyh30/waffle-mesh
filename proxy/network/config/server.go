package config

type TlsConfig struct {
	CertFilePath string
	KeyFilePath string
}

type ServerConfig struct {
	Host string
	Port int

	TlsConfig TlsConfig
}