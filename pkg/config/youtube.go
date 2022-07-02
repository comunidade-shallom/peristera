package config

type Channel struct {
	Name string `fig:"name" yaml:"name"`
	ID   string `fig:"id" yaml:"id"`
}

type YouTube struct {
	Token    string    `fig:"token" yaml:"token"`
	Channels []Channel `fig:"channels" yaml:"channels"`
}
