package config

import "os"

type Config struct {
	GrpcPort              string
	HttpPort              string
	UserServiceHost       string
	UserServicePort       string
	PostServiceHost       string
	PostServicePort       string
	CertificatePath       string
	CertificateKeyPath    string
	ConnectionServiceHost string
	ConnectionServicePort string
}

func NewConfig() *Config {
	return &Config{
		GrpcPort:              getEnv("GATEWAY_GRPC_PORT", "8080"),
		HttpPort:              getEnv("GATEWAY_HTTP_PORT", "8090"),
		UserServiceHost:       getEnv("USER_SERVICE_HOST", "localhost"),
		UserServicePort:       getEnv("USER_SERVICE_PORT", "8085"),
		PostServiceHost:       getEnv("POST_SERVICE_HOST", "localhost"),
		PostServicePort:       getEnv("POST_SERVICE_PORT", "8086"),
		ConnectionServiceHost: getEnv("CONNECTION_SERVICE_HOST", "localhost"),
		ConnectionServicePort: getEnv("CONNECTION_SERVICE_PORT", "8087"),
		CertificatePath:       getEnv("CERTIFICATE_PATH", "certificates/dislinkt.cer"),
		CertificateKeyPath:    getEnv("CERTIFICATE_KEY_PATH", "certificates/dislinkt_private_key.key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
