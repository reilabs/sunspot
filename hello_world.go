package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func Hello() string {
	return "Hello, World!"
}

func main() {
	fmt.Println("Hello, World!")
	data, err := os.ReadFile("noir-samples/hello_world/target/hello_world.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var acirFile ACIRFile
	if err := json.Unmarshal(data, &acirFile); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
}
