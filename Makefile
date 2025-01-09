# Build padrão (docker)
build-docker:
	@SSH_KEY=$(shell cat /home/saimon/.ssh/id_ed25519_no_passphrase | base64 -w 0) && \
	docker build --build-arg SSH_KEY="$${SSH_KEY}" --build-arg ENVIRONMENT=docker -t hubly:docker . && \
	docker tag hubly:docker saimonribeiros/hubly:docker && \
	docker push saimonribeiros/keeper:docker

# Build para desenvolvimento (dev)
build-dev:
	@SSH_KEY=$(shell cat /home/saimon/.ssh/id_ed25519_no_passphrase | base64 -w 0) && \
	docker build --build-arg SSH_KEY="$${SSH_KEY}" --build-arg ENVIRONMENT=dev -t hubly:dev . && \
	docker tag hubly:dev saimonribeiros/hubly:dev && \
	docker push saimonribeiros/keeper:dev
