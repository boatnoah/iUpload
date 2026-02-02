package main

import (
	"github.com/boatnoah/iupload/internal/auth"
	"github.com/boatnoah/iupload/internal/processor"
)

type app struct {
	svc  *processor.Processor
	auth *auth.Auth
}
