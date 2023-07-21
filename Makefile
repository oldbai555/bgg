DockerDir := "docker"
DockerCmd := "docker"
DockerBuild := $(DockerCmd) build --platform linux/amd64


DockerImageName = ""
.PHONY: docker.build
docker.build:
	$(DockerBuild) -t $(DockerImageName) -f $(DockerDir)/Dockerfile-$(DockerImageName) .

.PHONY: docker.clean
docker.clean:
	docker stop $(DockerImageName)
	docker rm $(DockerImageName)
	docker rmi $(DockerImageName)

.PHONY: d_b_gateway
d_b_gateway: DockerImageName := gateway
d_b_gateway: docker.build