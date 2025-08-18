package main

import (
	"fmt"
	"io"
	"os"
)

const metadataPath = "/home/nichepenguin/public_html/kno/bugs.json"

func main() {
	data, err := io.ReadAll(os.Stdin)

	if err != nil {
		httpError(500, fmt.Sprintf("Failed to read stdin: %v", err))
		return
	}

	req, err := parse(data)
	if err != nil {
		httpError(400, fmt.Sprintf("Failed to parse json: %v", err))
		return
	}

	response, err := process(req)
	if err != nil {
		httpError(500, fmt.Sprintf("Failed to process request: %v", err))
		return
	}

	fmt.Println("HTTP/1.1 200 OK")
	fmt.Println("Content-Type: application/json")
	fmt.Println()
	fmt.Println(response)
}

var statusText = map[uint]string{
	500: "Internal Server Error",
	400: "Bad Request",
}

func httpError(code uint, msg string) {
	fmt.Printf("HTTP/1.1 %d %s\n", code, statusText[code])
	fmt.Printf("Content-Type: text/plain\n\n")
	fmt.Println(msg)
}
