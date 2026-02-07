package main

import (
	"fmt"
	"github.com/boatnoah/iupload/internal/auth"
	"github.com/boatnoah/iupload/internal/processor"
	"os"
)

type app struct {
	svc  *processor.Processor
	auth *auth.Auth
}

func RequiredEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Errorf("missing required environment variable: %s", key))
	}

	return v
}
