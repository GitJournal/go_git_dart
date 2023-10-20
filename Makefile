
build:
	go build -o gitjournal.so -buildmode=c-shared gitjournal.go
	dart run ffigen
