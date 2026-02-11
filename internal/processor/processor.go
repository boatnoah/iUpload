package processor

import (
	"context"
	"fmt"
	"io"

	"github.com/boatnoah/iupload/internal/storage"
	"github.com/google/uuid"
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

func (p *Processor) UploadImage(ctx context.Context, userID uuid.UUID, fileName string, body io.ReadCloser, contentType string) (*storage.Image, error) {

	key := fmt.Sprintf("%s/%s", userID, fileName)
	image, err := p.storage.ImageStorage.Create(ctx, userID, key, contentType)

	if err != nil {
		return nil, err
	}

	err = p.objectStore.Put(ctx, key, body, contentType)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (p *Processor) GetByImageId() {
}

func (p *Processor) DeleteByImageId() {
}

func (p *Processor) TranformImage(operation string) {
}
