package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/inspii/edgex-thingsboard/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	internal.Main(ctx, cancel, mux.NewRouter(), nil)
}
