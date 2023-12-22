
build:
	cd src && go build -o gitjournal.so -buildmode=c-shared .
	dart run ffigen
