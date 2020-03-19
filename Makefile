.PHONY: test clean

trk: *.go
	go build

test:
	go test

clean:
	rm trk
