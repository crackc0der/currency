reset:
	@docker ps -aq | xargs docker rm -f || true
	if [ -d "data" ]; then \
		sudo rm -r data; \
	fi
