package main

import (
	docker "LHS/backend/Processors/Docker"
)

func main() {
	docker.InitDockerClient()
	docker.CloseClient()
}
