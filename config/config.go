package config

type UserServiceConfig struct {
	Port                  string
	DecoderServiceAddress string
}

type DecoderServiceConfig struct {
	Port               string
	UserServiceAddress string
}
