run-service.auth:
	cd service.auth && go run . && cd -

build-service.auth:
	go build -o dist/cc-service-auth ./service.auth/cmd

build-static-service.auth:
	go build -ldflags "-s -w -extldflags '-static'" -tags "osusergo,netgo" -trimpath -o dist/cc-service-auth ./service.auth/cmd