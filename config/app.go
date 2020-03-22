package config

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port struct {
			Http int `mapstructure:"http"`
		}
	}
	Database struct {
		Type          string `mapstructure:"type"`
		Port          int    `mapstructure:"port"`
		Schema        string `mapstructure:"schema"`
		Host          string `mapstructure:"host"`
		User          string `mapstructure:"user"`
		Pass          string `mapstructure:"pass"`
		MigrationPath string `mapstructure:"migration_path"`
	}
	RequestPath struct {
		Ping        string `mapstructure:"ping"`
		GetArticle  string `mapstructure:"get_articles"`
		PostArticle string `mapstructure:"post_articles"`
		GetTags     string `mapstructure:"get_tags"`
	}
}
