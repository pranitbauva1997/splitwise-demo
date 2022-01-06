# Introduction

This is a demo splitwise app built for Setu's interview.

## Setup Local Machine

[Install go](https://go.dev/doc/install).

Install postgres and start the service.

If on Mac:

```shell
make install_mac_dependencies
```

If on Ubuntu:
```shell
make install_ubuntu_dependencies
```

Setup GitHooks:

```shell
make setup_githooks
```

Setup local database:

```shell
make setup_db
```

How to build code:

```shell
make build
```

How to run code:

```shell
./web
```

Create a docker image:

```shell
make docker-build
```

Run the created docker image in a container:

```shell
make docker-run
```

Note: While the docker build is successful, docker run still throws errors while
connecting to the local postgres instance. To fix this, I have to introduce `DB_HOST`,
`DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` environment variables and then pass
them into the docker run command and set the `DB_HOST` as `host.docker.internal` but
I haven't found the time to implement this. Currently, the connection parameters are
hard-coded which also have to be picked up from environment variables.

After the `web` process is up and running, it's time to set up Swagger to interact
with the server.

Pull the swagger docker image:

```shell
make docker-pull-swagger
```

Run swagger:

```shell
make docker-run-swagger
```

Now make the API calls and enjoy!