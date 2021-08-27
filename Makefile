compile:
	go build -o bin/ory-poc *.go

build:
	docker-compose build kratos-selfservice-ui-node

run-all:
	docker-compose \
      up --build --force-recreate -d

run:
	docker-compose up -d --build \
      kratos-selfservice-ui-node

clean:
	docker-compose \
      rm -s -v -f
