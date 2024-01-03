# Shadowserver-API-go
Simple Golang program to access Shadowserver API

[![Go Report Card](https://goreportcard.com/badge/github.com/AM-CERT/Shadowserver-API-go)](https://goreportcard.com/report/github.com/AM-CERT/Shadowserver-API-go)
![GitHub last commit](https://img.shields.io/github/last-commit/AM-CERT/Shadowserver-API-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/AM-CERT/Shadowserver-API-go)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/AM-CERT/Shadowserver-API-go)

This software is based on [official API client](https://github.com/The-Shadowserver-Foundation/api_utils) and additional details can be found [here](https://github.com/The-Shadowserver-Foundation/api_utils/wiki)

## Description
The tool is developed and used by [AM-CERT](https://am-cert.am/) and include following capabilities:
- Query Shadowserver REST API
- Download Shadowserver reports in specified directory
- Schedule automated download of reports ([see usage with `systemd`](#usage-with-systemd))

## Usage
Download precompiled binaries from [latest release](https://github.com/AM-CERT/Shadowserver-API-go/releases) for your architecture or clone the repository and compile it with the following command
```
GOOS=linux GOARCH=amd64 go build -o shadowserver-api-go-linux-amd64 github.com/AM-CERT/Shadowserver-API-go/cmd/shadowserver-api-go
```
Download configuration `.env` file and edit with your API credentials
```
curl -o .env https://raw.githubusercontent.com/AM-CERT/Shadowserver-API-go/main/.env
```
Run the binary to see the usage
```
./shadowserver-api-go-darwin-arm64 -h
```
Run without any param to check the credentials are valid
```
./shadowserver-api-go-darwin-arm64
[::] Jan  3 15:22:13.150 [I] [app:./shadowserver-api-go-darwin-arm64] starting
{
 "pong": "2024-01-03 11:22:13Z"
}
```
Download reports as specified in `.env` file, report directory must exist
```
./shadowserver-api-go-darwin-arm64 -reports
```
You can pass parameter to the API call:
```
./shadowserver-api-go-darwin-arm64 -method reports/query -param '{"query":{"geo":"AM", "type":["|sinkhole","|honeypot"],"date":"2023-02-14"},"limit":1}'
```

### Usage with `systemd`
You can create a service to automatically download daily reports and keep an up-to-date directory structure with reports.

Create [systemd service](https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html) file:
```
# /etc/systemd/system/shadowserver-api-go.service
[Unit]
Description=Shadowserver-API-go service
After=network.target

[Service]
WorkingDirectory=/opt/Shadowserver-API-go
ExecStart=/opt/Shadowserver-API-go/shadowserver-api-go-linux-amd64 -reportsCron

[Install]
WantedBy=multi-user.target
```
Download configuration `.env` file and edit with your API credentials, download the binary for your OS/Arch
```
mkdir /opt/Shadowserver-API-go
cd /opt/Shadowserver-API-go
curl -o .env https://raw.githubusercontent.com/AM-CERT/Shadowserver-API-go/main/.env
wget https://github.com/AM-CERT/Shadowserver-API-go/releases/download/v0.1/shadowserver-api-go-linux-amd64
```
Reports directory must exist, create the directory mentioned in `.env` file:
```
mkdir /opt/Shadowserver-API-go/reports
```
Systemd command to reload, start, status and show logs:
```
systemctl daemon-reload
systemctl start shadowserver-api-go
systemctl status shadowserver-api-go
journalctl -f -u shadowserver-api-go
```
### Using as a golang library
If you want to extend the capabilities of this tool by developing your own app, all functions are safe to import and use.

For example to make an API call:
```go
// make a param struct
params := make(model.ShadowserverParam)

// load the query params
err = json.Unmarshal([]byte(`{"query":{"geo":"AM", "type":["|sinkhole","|honeypot"],"date":"2023-02-14"},"limit":1}`), &params)

// make the API call
data, _ := shadowserver.CallApi(method, params)

// print the result
shadowserver.PrintJson(data, true)
```

## Need Help
Please check the [issues](https://github.com/AM-CERT/Shadowserver-API-go/issues) first and open a new one if you can't find a solution.

## License
* GNU GENERAL PUBLIC LICENSE. [LICENSE](LICENSE) or https://www.gnu.org/licenses/gpl-3.0.en.html#license-text