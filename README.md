# Shark Chat
###### By [Alex Beaver](https://github.com/alexander-beaver), [Alex Krantz](https://github.com/akrantz01), and [Aidan Sacco](https://github.com/asacco796940)

This is a simple chat app with a backend written in [Golang](https://golang.org) and a cross-platform [React Native](https://react-native.org) frontend.
The name of the project was chosen by our teacher.
The app uses a [PostgreSQL](https://postgresql.org) database to store all the chat, message, and user data.
It also uses a SMTP server to send emails about verification and password reset.

## Why?
This project was for AP Computer Science Principles during the 2019-2020 school year.

## Deployment
The server can be deployed using either a [Docker](https://docker.com) container or a standalone binary.
Either environment variables or a configuration file can be used to configure the server.
If you use environment variables, each key will be in all caps, prefixed with `CHAT_`, and have the periods (`.`) replaced with underscores (`_`).
For example, the key `database.host` would be `CHAT_DATABASE_HOST`.
If you want to use a configuration file, it should be named `config` and contain all the keys you want to change.
It supports either JSON, YAML or TOML.

To see a reference for all the keys, see the [documentation](docs/api/configuration.md#configuration-keys).

### Docker Container
This is the preferred way to deploy the server and arguably the easiest.
You can create the API container manually or use a [`docker-compose.yaml`](docker-compose.yaml) file.

#### Manual
Running with a configuration file:
```shell script
docker run --name chat-app -v $PWD/config.yaml:/config.yaml docker.pkg.github.com/akrantz01/apcsp/chat-app-api:latest
```
Running with environment variables:
```shell script
docker run --name chat-app -e CHAT_DATABASE_HOST=127.0.0.1 [other environment variables] docker.pkg.github.com/akrantz01/apcsp/chat-app-api:latest
```

#### Docker Compose
The compose file uses [MailHog](https://hub.docker.com/r/mailhog/mailhog) as a mock SMTP server without a proper domain.
It also has a web UI to see all outgoing emails.
The API will be accessible on port `8080` on the IP address of your machine.
```shell script
wget -O docker-compose.yaml https://raw.githubusercontent.com/akrantz01/apcsp/master/docker-compose.yaml
docker-compose up -d
```

**NOTE:** The API container may fail to start saying that it cannot access the Postgres container.
This is because the Postgres container can sometimes take longer to be ready to accept connections.
To remedy this, simply restart the API container with the following command: `docker container start <container name>`

### Standalone Binary
To run from a standalone binary, download it from the [releases tab](https://github.com/akrantz01/apcsp/releases/latest).
Currently, only binaries for 64-bit Linux are built, but more architectures and operating systems are planned.

Steps to build:
1. Download the server binary
1. Mark as executable with `chmod +x chat-app_[arch]_[os]`
1. Ensure a mail server and Postgres instance are accessible
1. Run the server with `./chat-app`
