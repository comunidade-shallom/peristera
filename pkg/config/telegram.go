package config

type SetOfCommand struct {
	Text        string `fig:"text" yaml:"text"`
	Description string `fig:"description" yaml:"description"`
}

type CommandMapper struct {
	Endpoints []string `fig:"endpoints" yaml:"endpoints"`
	Handler   string   `fig:"handler" yaml:"handler"`
}

type MenuMapper struct {
	Text    string `fig:"text" yaml:"text"`
	Handler string `fig:"handler" yaml:"handler"`
}

type TelegramCommands struct {
	SetOf   []SetOfCommand  `fig:"set_of" yaml:"set_of"`
	Mappers []CommandMapper `fig:"mappers" yaml:"mappers"`
	Menu    []MenuMapper    `fig:"menu" yaml:"menu"`
}

type Telegram struct {
	Token     string           `fig:"token" yaml:"token"`
	Commands  TelegramCommands `fig:"commands" yaml:"commands"`
	Broadcast []int64          `fig:"broadcast" yaml:"broadcast"`
	Admins    []int64          `fig:"admins" yaml:"admins"`
	Roots     []int64          `fig:"roots" yaml:"roots"`
}
