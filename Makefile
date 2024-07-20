format:
	@goimports -w .

run:
	docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock go-sandbox

dev: format run
