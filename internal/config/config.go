package config

type Config struct {
	Port                    uint
	CoinAPIUrl              string
	DefaultCurrency         string
	EmailStoragePath        string
	EmailToSendFrom         string
	EmailToSendFromPassword string
	EmailServiceUrl         string
	EmailServicePort        int
	EmailSubject            string
}
