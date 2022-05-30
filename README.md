# go-sparkplug-primary
Sparkplug B primary application implemented in Go

## Configuration

The application can be configured using the following environment variables (see [.env.example](./.env.example) for an example):

| Variable            | Default                  | Description                                                                   |
| ------------------- | ------------------------ | ----------------------------------------------------------------------------- |
| `LOG_FORMAT`        | `"text"`                 | Log format. Can be `text` or `json`.                                          |
| `LOG_FILE`          | `""`                     | Log file for the application (empty string for stdout)                        |
| `LOG_LEVEL`         | `"info"`                 | Log level for the application (panic, fatal, error, warn, info, debug, trace) |
| `MQTT_ENDPOINT`     | `"tcp://localhost:1883"` | Endpoint of MQTT broker                                                       |
| `MQTT_CLIENT_ID`    | `"go-primary"`           | Client ID for MQTT connection                                                 |
| `MQTT_USERNAME`     | `""`                     | Username for MQTT connection                                                  |
| `MQTT_PASSWORD`     | `""`                     | Password for MQTT connection                                                  |
| `SPARKPLUG_HOST_ID` | `"go-primary"`           | Host ID for `STATE` messages                                                  |