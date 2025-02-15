package nats

type Config struct {
	// Name is an optional name label which will be sent to the server
	// on CONNECT to identify the client.
	Name string `mapstructure:"name"`

	User string `mapstructure:"user"`

	Password string `mapstructure:"password"`

	Servers []string `mapstructure:"servers"`
}
