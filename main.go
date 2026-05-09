package main

import (
	"fmt"
	"os"

	"ascii-art-web-dockerize/src"
)

func main() {
	ipAddr := "0.0.0.0:8080"
	fmt.Println("Server listening on ", ipAddr)
	if err := src.StartServer(ipAddr); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
