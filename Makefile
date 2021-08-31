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

create-hydra-client:
	docker-compose -f hydra-docker-compose.yaml exec hydra hydra clients create \
        --endpoint http://127.0.0.1:4445 \
        --id auth-code-client \
        --secret secret \
        --grant-types authorization_code,refresh_token \
        --response-types code,id_token \
        --scope openid,offline \
        --callbacks http://127.0.0.1:4455/callback

clean-hydra:
	docker-compose -f hydra-docker-compose.yaml \
          rm -s -v -f

run-ui:
	docker-compose up -d --build \
      kratos-selfservice-ui-node

clean:
	docker-compose \
      rm -s -v -f
