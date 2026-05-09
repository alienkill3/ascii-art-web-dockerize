FROM golang:1.25-alpine

LABEL maintainer="ktsioupr"
LABEL app="go-web-server"
LABEL version=1.0
LABEL description="ascii art web"

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
