package config

type Store struct {
	Path string `fig:"path" yaml:"path"  default:"/peristera.d/store"`
}
