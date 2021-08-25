build:
	go build -o bin/ory-poc *.go
run:
	docker-compose \
      up --build --force-recreate -d

clean:
	docker-compose \
      rm -s -v -f
