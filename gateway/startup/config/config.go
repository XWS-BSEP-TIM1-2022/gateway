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
	JobServiceHost        string
	JobServicePort        string
	RolePermissions       map[string][]string
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
		JobServiceHost:        getEnv("JOB_SERVICE_HOST", "localhost"),
		JobServicePort:        getEnv("JOB_SERVICE_PORT", "8088"),
		CertificatePath:       getEnv("CERTIFICATE_PATH", "certificates/dislinkt.cer"),
		CertificateKeyPath:    getEnv("CERTIFICATE_KEY_PATH", "certificates/dislinkt_private_key.key"),
		RolePermissions: map[string][]string{
			"ADMIN": []string{"user_getAll", "user_read", "user_write", "user_delete", "post_read", "post_write", "post_delete", "post_getAll", "job_read", "job_write", "job_delete", "connection_read", "connection_write", "connection_delete"},
			"USER":  []string{"post_read", "user_read", "user_write", "post_write", "post_delete", "job_read", "job_write", "job_delete", "connection_read", "connection_write", "connection_delete"},
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
