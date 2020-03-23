package config

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port struct {
			Http int `mapstructure:"http"`
		}
	}
	Database struct {
		Type   string `mapstructure:"type"`
		Port   int    `mapstructure:"port"`
		Schema string `mapstructure:"schema"`
		Host   string `mapstructure:"host"`
		User   string `mapstructure:"user"`
		Pass   string `mapstructure:"pass"`
	}
}
