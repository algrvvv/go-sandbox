format:
	@goimports -w .

build:
	@docker build -t go-runner ./executor
	@docker save -o go-runner.tar go-runner
	@docker build -t go-sandbox .

run:
	docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock -v ./go-runner.tar:/go-runner.tar go-sandbox

dev: build run

dev: format run
