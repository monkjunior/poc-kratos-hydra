start-all:
	docker-compose \
      up --build  -d

stop-all:
	docker-compose stop

clean-all:
	docker-compose \
      rm -s -v -f

compile:
	go build -o bin/ory-poc *.go

build-ui:
	docker-compose build ui-node

update-ui:
	docker-compose up -d --build \
      ui-node

create-hydra-client:
	docker-compose exec hydra hydra clients create \
        --endpoint http://127.0.0.1:4445 \
        --id auth-code-client \
        --secret secret \
        --grant-types authorization_code,refresh_token \
        --response-types code,id_token \
        --scope openid,offline \
        --callbacks http://127.0.0.1:5555/callback
