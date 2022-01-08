package main

import (
	"context"
)

func main() {

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
