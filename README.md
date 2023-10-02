# Testapp 20230930

### Prerequisites

* Go 1.21
* [optional] Make

### How to configure

Application is configured via yaml config file. Consult with [config.yaml](config.yaml) to see the example.

| Parameter       | Description                                                                                             | Default config                                                                                     |
|-----------------|---------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------|
| root_url        | Initial page                                                                                            | https://ru.wikipedia.org                                                                           |
| data_directory  | Directory to store downloaded content                                                                   | "./data"                                                                                           |
| database_path   | Path to SQLite database. Will be created if not exists                                                  | "./data/sqlite.db"                                                                                 |
| worker_count    | Amount of workers running in parallel. Values less than 1 are treated as 1                              | 3                                                                                                  |
| content_types   | Valid content type list. Also defines extensions for saving files without extensions provided by server | { text/html: ".html", text/css: ".css", application/javascript: ".js" application/json: ".json"}   |
| request_headers | Request header overrides for downloading                                                                | {user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0} |

To use a configuration file, pass a path to it via `--config-path` flag.

### How to run

You could use `make run-local` in the project directory to build and run the application with default config path.

Alternatively, you could run the application with `go run ./cmd/app --config-path {path_to_config.yaml}` to use
different config file.
