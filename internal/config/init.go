package config

var (
	Values *Config
)

func init() {
	Values = Read()
}
