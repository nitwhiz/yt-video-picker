.PHONY: build
build: build_app build_server

.PHONY: build_app
build_app:
	cd app && yarn build

.PHONY: build_server
build_server:
	cd server && go build -o dist/ytvp-server
