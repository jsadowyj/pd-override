.PHONY: build run clean test

build:
	go build .
run: build
	./pd-override M-F@0900-1700 M,T,R@0900-1200 M,T,R@1600-1900
clean:
	rm -f pd-override
test:
	go test .
