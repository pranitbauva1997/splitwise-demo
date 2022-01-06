build:
	go build -o web ./cmd/web
	chmod +x web

run:
	./web

clean:
	rm -f web
	go clean

build_and_run: clean build run

test:
	go test -cover ./...

setup_db: create_db create_db_tables

create_db:
	psql -U postgres -c 'CREATE DATABASE splitwise_demo_test;'

create_db_tables:
	psql -U postgres -d splitwise_demo_test -f migrations/pg/0001_init.up.sql

destroy_db:
	psql -U postgres -c 'DROP DATABASE splitwise_demo_test;'

init: setup_githooks install_dev_dependencies

install_mac_dependencies:
	brew install shellcheck

install_ubuntu_dependencies:
	apt-get update --fix-missing
	apt-get install shellcheck -y

install_dev_dependencies:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

setup_githooks:
	rm -f .git/hooks/pre-commit
	ln -s ${PWD}/scripts/pre-commit.sh .git/hooks/pre-commit

shellcheck:
	shellcheck scripts/*.sh

lint:
	golangci-lint run ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

all-static-checks: vet lint staticcheck shellcheck

docker-pull-swagger:
	docker pull swaggerapi/swagger-ui

docker-build:
	docker build . -t splitwisedemo:latest

docker-run:
	docker run --network=host -p 8000:8000 splitwisedemo:latest
