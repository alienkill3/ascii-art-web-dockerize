
| Concept | Meaning |
|:--:|:--:|
| Dockerfile | Recipe |
| Image | Built app snapshot |
| Container | Running app |
| Metadata | Labels/info |
| Cleanup | Remove unused stuff |

### Dockerfile

A Dockerfile is a text file that contains instructions for building your app into a Docker image.
It defines:
- Base environment (e.g., Go version)
- Copying your code
- Installing dependencies
- Running your app

_example of `Dockerfile`:_
```
# Use official Go image
FROM golang:1.21-alpine

# Add metadata
LABEL maintainer="yourname"
LABEL app="go-web-server"
LABEL version="1.0"

# Set working directory
WORKDIR /app

# Copy files
COPY . .

# Build the app
RUN go build -o main .

# Expose port
EXPOSE 8080

# Run the app
CMD ["./main"]
```

### Image

A Docker image is a packaged version of your app created from the Dockerfile.
Think of it like a snapshot or template of your app
It's read-only
You can create many containers from it

_example:_
```
my-go-server:1.0
```
_build the docker image:_
```
docker build -t go-web-server:1.0 .
```
_how to check the images;_
```
docker images
```


### Container
A container is a running instance of your image.
Think of it like a live version of your app
You can start, stop, delete containers

_how to create and run the container:_
```
docker run -d -p 8080:8080 --name go-container go-web-server:1.0
```


### Metadata (labels)
Metadata = extra information you attach to Docker objects.
In Docker, this is done using labels.
Example:
```
LABEL maintainer="yourname"
LABEL version="1.0"
LABEL description="Go web server project"
```
Why?

-   Helps identify, organize, filter containers/images

_how to check the metadata:_
```
docker inspect go-web-server:1.0
```

### Garbage collection (cleanup unused objects)

Docker stores:

-   images
-   containers
-   volumes
-   networks

Over time → unused ones pile up.
You must **clean them**.

_example:_
```
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Clean everything unused
docker system prune
```

## Optional: use smaller image (multi-stage build)

```
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

LABEL maintainer="yourname"

EXPOSE 8080
CMD ["./main"]
```

## Docker compose

Docker Compose is a tool that lets you define your application using a YAML file.
Instead of running long commands you define everything in a file, then run:
```
docker compose up
```

_example of `docker-compose.yml`:_
```

version: "3.9"

services:
  web:
    build: .
    image: go-web-server:1.
    container_name: go-web-container
    ports:
      - "8080:8080"
    labels:
      maintainer: "yourname"
      app: "go-web-server"
      version: "1.0"
```

_meaning of commands:_
```
docker compose up --build   # build + run
docker compose up -d        # run in background
docker compose down         # stop/remove
docker system prune         # cleanup
```