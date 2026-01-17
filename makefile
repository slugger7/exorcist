dtos:
	mkdir -p ./apps/server/ts
	go run ./apps/server/cmd/enum-export
	tygo generate
	rm -rf ./apps/web/src/dtos
	mv ./apps/server/ts ./apps/web/src/dto
