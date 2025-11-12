## Link Checker HTTP server

### About service
HTTP server server written by **golang** The service checks resource link availability. It assigns a unique number to a link group and generates a PDF report based on the group number, showing the status of all resources. The service provides backups of all data.

### Technologies and libraries used
- [viper](https://github.com/spf13/viper) - Go configuration with fangs!
- [gofpdf](https://github.com/phpdave11/gofpdf?tab=readme-ov-file) - Library for creating and working with PDF files on golang

### required version golang
version >= 1.22 
### Service features
- Automatic data saving when the service is stopped.
- Checking multiple links per request, completely asynchronous.
- PDF file generation
### Instructions

#### 1. Command to build the project
```bash
 go build cmd/linkChecker/main.go
```
#### 2. Edit congfig file
Open `./config/config.yaml` file and and make changes to it following the instructions. 

**WARNING!** If you change the path for backups, the folder where they will be stored must exist before starting the server.

#### 3. Run server
```bash
./main
```

### API endpoints
#### 1. Get links state `http://host:port/links/status`
- METHOD: `POST`
- Request body:
```json 
{
    "links": [
        "google.com", "ya.ru"
    ]
}
```
- Example curl request:
```bash
curl -X POST http://127.0.0.1:8008/links/status -d '{"links": ["google.com", "ya.ru"]}'
```
- Response:
```json
{
    "links": {
        "google.com":"available",
        "ya.ru":"available"
    },
    "links_num":1
}
```

#### 2. Get links_group report `http://host:port/links/group/report`

- METHOD: `GET`
- Request body:
```json 
{
    "links_list": [1, 2]
}
```
- Example curl request:
```bash
curl -X GET http://127.0.0.1:8008/links/group/report -d '{"links_list": [1, 2]}' --output report.pdf
```
receive pdf file report and move it `./report.pdf` file
- Response:
`PDF file`

## TODO 
- [ ] Create tests

## LiCENSE