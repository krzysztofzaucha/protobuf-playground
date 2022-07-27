package internal

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Config struct {
	Server Server `json:"server"`
}
