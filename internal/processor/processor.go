package processor

import "github.com/boatnoah/iupload/internal/storage"

type Processor struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Processor {
	return &Processor{storage: storage}
}
