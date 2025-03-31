# include .env
# export

run:
	cd cmd/ && go run main.go

migration_create:
	@migrate create -dir db/migrations -ext sql -seq $(name)

migration_up:
	@migrate -path db/migrations/ -database  $(DB_CONFIG) -verbose up

migration_down:
	@migrate -path db/migrations/ -database  $(DB_CONFIG) -verbose down $(n)

migration_force:
	@migrate -path db/migrations/ -database  $(DB_CONFIG) -verbose force $(n)

migration_up_test:
	@migrate -path db/migrations/ -database  $(TEST_DB_CONFIG) -verbose up	

test:
	gotestsum --format testname -- -v ./...