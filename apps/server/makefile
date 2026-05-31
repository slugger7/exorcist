run:
	go run ./cmd/exorcist

get:
	go get ./cmd/exorcist

update:
	go get -u ./...

exorcist:
	go build -o ./build/exorcist ./cmd/exorcist

build: exorcist

clean:
	rm -rf build
	rm -rf ts/*

run-migrations:
	./scripts/run-migrations.sh

undo-migration:
	./scripts/undo-migration.sh

models:
	./scripts/update-models.sh

generate-diagrams:
	./scripts/generate-diagrams.sh
