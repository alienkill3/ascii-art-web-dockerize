## Build the image using the project Dockerfile

```sh
docker image build -f Dockerfile -t go-web-server:1.0 .
```

**To check the docker images run:**

```sh
docker images
```

**example output:**
```
IMAGE               ID             DISK USAGE   CONTENT SIZE   EXTRA
go-web-server:1.0   4fbb4067b8c5        467MB         97.2MB        
```

## Start the container using the above image

```sh
docker container run -p 8080:8080 --detach --name go-web-container go-web-server:1.0
```

**To check the running containers run:**

```sh
docker ps -a
```

**example output:**
```
CONTAINER ID   IMAGE               COMMAND    CREATED         STATUS         PORTS                                         NAMES
c2e33603e9b4   go-web-server:1.0   "./main"   2 minutes ago   Up 2 minutes   0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp   go-web-container
```

## Run a command in the default working directory of the container.

```sh
docker exec -it go-web-container /bin/sh
```

## How to check metadata

```sh
docker inspect --format='{{json .Config.Labels}}' go-web-container
```

## Simplify the above process

We simplify the steps of building and running with docker compose. Just run:

```sh
docker compose up
```

## Clean up

**Commands to clean/remove containers:**
```sh
docker container rm <container-name>
docker image rm <image-name>
docker system prune
```