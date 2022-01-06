package main

import (
	"context"
	"fmt"
	"io/ioutil"
)

func main() {

	IOReadDir("/go/")
	IOReadDir("/go/templates/")

	ctx := context.Background()
	minion, bootstrapError := bootstrap(nil, ctx)
	exitOnError(bootstrapError)

	executionError := minion.Run(ctx)
	exitOnError(executionError)
}

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func IOReadDir(root string) {
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range fileInfo {
		fmt.Println(file.Mode(), " ", root, "/", file.Name())
	}
}
