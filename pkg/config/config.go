package config

type App struct {
	RootURL         string            `mapstructure:"root_url"`
	DataDirectory   string            `mapstructure:"data_directory"`
	DatabasePath    string            `mapstructure:"database_path"`
	DownloaderCount int               `mapstructure:"downloader_count"`
	ExtractorCount  int               `mapstructure:"extractor_count"`
	ContentTypes    map[string]string `mapstructure:"content_types"`
}
