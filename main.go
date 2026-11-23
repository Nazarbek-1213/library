package main

import (
	"fmt"
	"library/config"
)

func main() {
	fmt.Println(config.Load())
}
