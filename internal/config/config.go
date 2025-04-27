package config

type Config struct {
	DBURL         string
	JWTSecret     []byte
	EncryptionKey string
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPass      string
	EmailFrom     string
	Port          string
}

func Load() *Config {
	return &Config{
		DBURL:         "postgres://javaconnect:4-test@192.168.1.105:5432/gobank?sslmode=disable",
		JWTSecret:     []byte("c852470a8904dbb453d34d2451a3fa85ded41d5d34838e931d9bc7124f893cf99430757a8017f2373da9c8bd6ef0f70cde9b1bdf87013d6908f9825ab17b0151895ee3695a0e3df4c6d8047140c9de08c532900f7e7f82d7c663da311978694939a1d2e9268ae3459334b20a25ca97b715a6b8d99f92150b4afb503d362bb43630384bc36d3ad321a533e29e27862cb35b4bf01303342f059fa02cc22f8aed39a9f2a1be4326a865f969381e60c53b6065f52c35d22e48d633eefc5db057e05dc84fa3c6f8e7d0249534c1b52ff26ce9de9fbabc7c05677514170bff7cef49d767e3f70b9a36de1d4f2d21dc77440f5f6bc2719b76cf3d5c8861cc6854ba7c66"),
		EncryptionKey: "9SzkstvA9ESaq1B/+Hvq+R+u+3OKlBYXIMB32moWmnVPF9IKj5tQQ4Vnj2TeX05+",
		Port:          "8080",
	}
}
