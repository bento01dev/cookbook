package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bento01dev/cookbook/internal/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Stdout, os.Stdout, os.Args, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
