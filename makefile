APP_NAME=adviser
DEV_CONFIG=config.dev.yaml

test:
	go test ./...

build:
	docker build -t $(APP_NAME) .
	docker scan $(APP_NAME)

stop:
	container_ids=$$(docker ps -a -q --filter ancestor=$(APP_NAME)); \
    	if [ ! -z "$$container_ids" ]; then \
    		docker stop $$container_ids; \
    	fi

run:
	docker run -d $(APP_NAME)

deploy: test build stop run

dev:
	go run bonefabric/adviser/cmd/adviser -config $(DEV_CONFIG)