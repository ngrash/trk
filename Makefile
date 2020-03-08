.PHONY: test clean

trk: main.go
	go build

test:
	go test

clean:
	rm trk
