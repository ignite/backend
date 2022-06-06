package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ignite-hq/blockchain-backend/cmd"
	"github.com/ignite-hq/cli/ignite/pkg/clictx"
)

// TODO: add event/attribute blacklisting (to be able to exclude binary ones)
// TODO: add logging support (add logs to the collector)

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
