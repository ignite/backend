package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ignite/cli/ignite/pkg/clictx"

	"github.com/ignite/backend/cmd"
)

func main() {
	ctx := clictx.From(context.Background())

	if err := cmd.New().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		fmt.Println("aborted")
	}
}
