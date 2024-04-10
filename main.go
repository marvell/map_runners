package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := LoadConfig()
	exitOnError(err)

	ctx := context.Background()

	app := NewApplication(cfg)

	defer func() {
		if err := app.Close(ctx); err != nil {
			log.Printf("ERR %s", err)
		}
	}()

	err = app.Run(ctx)
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Function: %s\n", funcName)
		fmt.Fprintf(os.Stderr, "File: %s:%d\n", file, line)

		os.Exit(1)
	}
}
