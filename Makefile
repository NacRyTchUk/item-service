gen:
	buf generate

run:
	docker-compose up back --build


run-env:
	docker-compose up nats --build

run-all: run-env run

drop:
	docker-compose down --volumes

