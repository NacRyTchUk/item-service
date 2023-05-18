BUF_PATH = /usr/local/bin/buf-Linux-x86_64


.PHONY: gen
gen:
	@$(BUF_PATH) generate

run:
	docker-compose up back --build


run-dummy-broker:
	docker-compose up nats --build

run-all: run run-dummy-broker

drop:
	docker-compose down --volumes


db_clean:
	sudo rm -rf db/postgres/01-init.sql/ db/postgres/postgres-data/