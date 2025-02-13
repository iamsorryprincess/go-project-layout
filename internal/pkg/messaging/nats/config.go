package nats

type Config struct {
	// URL for connection to server
	URL string `mapstructure:"url"`

	// Name is an optional name label which will be sent to the server
	// on CONNECT to identify the client.
	Name string `mapstructure:"name"`
}
