package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

// func TestRun(t *testing.T) {
// 	err := run()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestRun(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancelCause(ctx)
	time.AfterFunc(time.Second, func() { cancel(fmt.Errorf("cancelled from testing")) })

	err := run(ctx, "7777")
	if err != nil {
		log.Fatalf("expected no error, but got: %s", err)
	}
}
