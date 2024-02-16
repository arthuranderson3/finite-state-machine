APP=fsm

build:
	CGO_ENABLED=false	go build -v -o ./$(APP) ./cmd/fsm

run:
	go run ./cmd/fsm

clean:
	rm ./$(APP)

get:
	go get gonum.org/v1/gonum/graph
	go get gonum.org/v1/gonum/graph/multi
