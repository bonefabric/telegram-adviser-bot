APP_NAME=adviser

test:
	go test ./...

build:
	docker build -t $(APP_NAME) .

stop:
	container_ids=$$(docker ps -a -q --filter ancestor=$(APP_NAME)); \
    	if [ ! -z "$$container_ids" ]; then \
    		docker stop $$container_ids; \
    	fi

run:
	docker run -d $(APP_NAME)

deploy: test build stop run