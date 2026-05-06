FROM golang:1.25-alpine

label maintainer="ksokos"
label app="go-web-server"
label version=1.0
label description="ascii art web"

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
