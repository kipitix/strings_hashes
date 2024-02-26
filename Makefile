###

export GOOSE_DRIVER := postgres
export GOOSE_DBSTRING := host=localhost port=5432 user=postgres password=postgres dbname=stringhashes sslmode=disable


protoc:
	@protoc ./hashcalc/api/*.proto --go-grpc_out=./hashcalc/pkg --go_out=./hashcalc/pkg

grpc-ui:
	@grpcui -plaintext localhost:50051

swagger-pull:
	@docker pull quay.io/goswagger/swagger

swagger-gen-srv:
	@etc/scripts/swagger.sh generate server \
	--spec=./hashkeeper/pkg/interfaces/restapi/spec/hashkeeper.yaml \
	--name="hashkeeper" \
	--target=./hashkeeper/internal/interfaces/restapi \
	--server-package=server \
	--exclude-main
	#
	# --main-package=../../../../cmd

swagger-gen-cnt:
	@etc/scripts/swagger.sh generate client \
	--spec=./hashkeeper/pkg/interfaces/restapi/spec/hashkeeper.yaml \
	--name="hashkeeper" \
	--target=./hashkeeper/pkg/interfaces/restapi \
	--client-package=client

swagger-serve:
	@etc/scripts/swagger.sh serve -p 8080 --no-open ./hashkeeper/api/hashkeeper.yaml

goose-install:
	@go install github.com/pressly/goose/v3/cmd/goose@latest

migrate:
	@goose -dir=./migrations up

unmigrate:
	@goose -dir=./migrations down

debug-up:
	@cd ./deployments/debug ; docker-compose up -d

debug-down:
	@cd ./deployments/debug ; docker-compose down

docker-up:
	@cd ./deployments ; docker-compose up -d ; cd - ; goose -dir=./migrations up

docker-build-up:
	@cd ./deployments ; docker-compose up -d --build ; cd - ; goose -dir=./migrations up

docker-down:
	@cd ./deployments ; docker-compose down
	
full-rebuild:
	@cd ./deployments ; docker-compose up -d --build --force-recreate --remove-orphans ; cd - ; goose -dir=./migrations up

unittests:
	@cd ./hashcalc/internal/strhash && go test -v && cd -

test-rest-send:
	@curl -X 'POST' \
	'http://localhost:8080/send' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '["string", "line", "строка"]'

test-rest-check:
	@curl -X 'GET' \
	'http://localhost:8080/check?ids=1,2,3,4,5' \
	-H 'accept: application/json'

run-hashcalc-race:
	@go run -race ./hashcalc/cmd/

run-hashkeeper-race:
	@export PORT=8080 && go run -race ./hashkeeper/cmd/

go-wrk-install:
	@go install github.com/tsliwowicz/go-wrk@latest

stress-rest-send:
	@go-wrk -M 'POST' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-body '["string", "line", "строка"]' \
	-c 10 \
	-d 10 \
	'http://localhost:8080/send'

stress-rest-check:
	@go-wrk -M 'GET' \
	-H 'accept: application/json' \
	-c 10 \
	-d 10 \
	'http://localhost:8080/check?ids=1,2,3,4,5'
