compile:
	go build -o bin/ory-poc *.go

build:
	docker-compose build kratos-selfservice-ui-node

run-all:
	docker-compose \
      up --build  -d

run-hydra:
	docker-compose -f hydra-docker-compose.yaml \
	  up -d --build --force-recreate

clean-hydra:
	docker-compose -f hydra-docker-compose.yaml \
          rm -s -v -f

run-ui:
	docker-compose up -d --build \
      kratos-selfservice-ui-node

clean:
	docker-compose \
      rm -s -v -f
