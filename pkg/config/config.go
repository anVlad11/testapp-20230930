package config

type App struct {
	RootURL        string            `mapstructure:"root_url"`
	DataDirectory  string            `mapstructure:"data_directory"`
	DatabasePath   string            `mapstructure:"database_path"`
	WorkerCount    int               `mapstructure:"worker_count"`
	ContentTypes   map[string]string `mapstructure:"content_types"`
	RequestHeaders map[string]string `mapstructure:"request_headers"`
}
