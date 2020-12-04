
build-docker:
	 docker build -t "6d6-bot:latest" --build-arg discord_bot_token=$(shell cat ./bot.token) .

stop:
	docker stop 6d6-bot && docker rm 6d6-bot

run: stop
	docker run -d --name 6d6-bot 6d6-bot:latest
