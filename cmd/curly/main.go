package main

import (
	"curly/pkg/core"
	"fmt"
)

func main() {
	_, err := core.NewApp()
	if err != nil {
		panic(err)
	}

	fmt.Println("exiting curly, thanks for using me.")
}
