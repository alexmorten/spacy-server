run:
	go run main/spacy.go

proto:
	docker build -t spacy-proto-c -f proto.Dockerfile .
	docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) al-proto-c

docker-run:
	docker build -t spacy-server .
	docker run --rm -p 4000:4000 --name spacy-server spacy-server
