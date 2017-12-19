# manifest-manager

CLI that manages manifest files.

### Installation

`go get -u github.com/ultrabluewolf/manifest-manager`

Assuming your go bin is on your path to use the CLI see help for available commands:

`manifest-manager -h|--help|help` i.e. `manifest-manager -h`

### Development

_using docker compose_

Create dotenv file and update env values as desired

`cp .env.example .env`

Build and start container

`docker-compose build`

`docker-compose up`

`docker-compose exec go ./bin/build.sh`

And in a new tab run CLI:

`docker-compose exec go manifest-manager -h`

or

_without docker compose_

grab the dependencies see [dep](https://github.com/golang/dep)

`dep ensure`

`go run cmd/manifest-manager/main.go -h`

Available env vars:

- LOGLEVEL

#### Tests

`docker-compose exec go ./bin/test.sh`

or

`./bin/test.sh`
