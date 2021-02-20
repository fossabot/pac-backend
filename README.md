# PAC Backend
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FMilos5611%2Fpac-backend.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FMilos5611%2Fpac-backend?ref=badge_shield)

Backend for the PAC 2020 application (Conferencing app)

## Running locally
To run locally, simply run the application with:

`go run .\main.go`

The application will start with an embedded sqlite3 database, on port 9090.

To initialize the database with test data, call the init endpoint with `curl -XPOST localhost:9090/initDB`. The call drops and recreates all tables and adds the test data into them. 

## Running as a part of PAC infrastructure
The infrastructure expects a docker image tagged as `pac-backend`. To build the image, run:

`docker build -t pac-backend .`

The image currently cannot run an embedded sqlite3 database, so it is intended to be used with mysql. The connection can be configured using environment variables.

## Environment variables

The PAC Backend application is configured using environment variables which are listed in the .env file in the root directory of the project.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FMilos5611%2Fpac-backend.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FMilos5611%2Fpac-backend?ref=badge_large)