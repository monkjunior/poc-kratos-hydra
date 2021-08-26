compile:
	go build -o bin/ory-poc *.go

build:
	docker-compose build kratos-selfservice-ui-node

run-all:
	docker-compose \
      up --build --force-recreate -d

run:
	docker-compose up \
	  kratos-selfservice-ui-node -d

clean:
	docker-compose \
      rm -s -v -f
