package config

type Channel struct {
	Name string `fig:"name" yaml:"name"`
	ID   string `fig:"id" yaml:"id"`
	URL  string `fig:"url" yaml:"url"`
}

type YouTube struct {
	Token    string    `fig:"token" yaml:"token"`
	Channels []Channel `fig:"channels" yaml:"channels"`
}

func (c Channel) GetURL() string {
	if c.URL != "" {
		return c.URL
	}

	return "https://www.youtube.com/channel/" + c.ID
}
