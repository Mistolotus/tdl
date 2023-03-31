package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Mistolotus/tdl/cmd"
	"github.com/fatih/color"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.New().ExecuteContext(ctx); err != nil {
		color.Red("%v", err)
	}

}
