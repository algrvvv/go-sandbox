# форматирование
format:
	@goimports -w .

# сборка проекта
build: rb
	@docker volume create code-files
	@docker save -o go-runner.tar go-runner
	@docker build --no-cache -t go-sandbox .

# запуск приложения используя докер
run:
	@docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock -v code-files:/tmp/go-sandbox -v ./go-runner.tar:/go-runner.tar go-sandbox

# для разработки. сразу форматирование, билд и запуск
dev: format build run

sd:
	@docker exec -it $(filter-out $@,$(MAKECMDGOALS)) /bin/sh

# билд раннера. rb - runner build
rb:
	@docker build --no-cache -t go-runner ./executor

# запуск раннера. rr - runner run
rr:
	@docker run --rm -v /tmp/go-sandbox:/app go-runner $(filter-out $@,$(MAKECMDGOALS))

# запуск раннер в интерактивном режиме. rd - runner debug
rd:
	@docker run -it --rm -v /tmp/go-sandbox:/app --entrypoint /bin/sh go-runner

# запуск основно сервера без докера (gr - go run)
gr:
	@mkdir -p /tmp/go-sandbox/
	@go run main.go
