
prefix_plugin:
	CGO_ENABLED=1 CC="/usr/local/bin/x86_64-linux-musl-gcc" GOOS=linux GOARCH=amd64 go build -buildmode=plugin  -o prefix_plugin.so prefix/*

suffix_plugin:
	CGO_ENABLED=1 CC="/usr/local/bin/x86_64-linux-musl-gcc" GOOS=linux GOARCH=amd64 go build -buildmode=plugin  -o suffix_plugin.so suffix/*
