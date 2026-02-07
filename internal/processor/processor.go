package processor

import (
	"context"
	"io"

	"github.com/boatnoah/iupload/internal/storage"
)

type ObjectStore interface {
	Put(ctx context.Context, key string, body io.Reader, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

type Processor struct {
	storage     *storage.Storage
	objectStore ObjectStore
}

func New(storage *storage.Storage, objectStore ObjectStore) *Processor {
	return &Processor{storage: storage, objectStore: objectStore}
}

func (p *Processor) UploadImage() {
}

func (p *Processor) GetByImageId() {
}

func (p *Processor) DeleteByImageId() {
}

func (p *Processor) TranformImage(operation string) {
}
