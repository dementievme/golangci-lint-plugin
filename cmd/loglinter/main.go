package main

import (
	"fmt"

	"github.com/dementievme/golangci-lint-plugin/internal/config"
)

func main() {
	config := config.MustLoad()
	fmt.Printf("%v", config)
}
