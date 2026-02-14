package processor

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/boatnoah/iupload/internal/storage"
	"github.com/google/uuid"
)

var (
	ErrorNotFound = errors.New("Unable to find image")
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

type ImagePayload struct {
	UserID      uuid.UUID
	FileName    string
	Body        io.ReadCloser
	ContentType string
}

func New(storage *storage.Storage, objectStore ObjectStore) *Processor {
	return &Processor{storage: storage, objectStore: objectStore}
}

func (p *Processor) UploadImage(ctx context.Context, imagePayload ImagePayload) (*storage.Image, error) {
	prefix := uuid.New()
	key := fmt.Sprintf("%s/%s-%s", imagePayload.UserID, imagePayload.FileName, prefix)
	image, err := p.storage.ImageStorage.Create(ctx, imagePayload.UserID, key, imagePayload.ContentType)

	if err != nil {
		return nil, err
	}

	err = p.objectStore.Put(ctx, key, imagePayload.Body, imagePayload.ContentType)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (p *Processor) GetByImageId(ctx context.Context, id uuid.UUID) (io.ReadCloser, string, error) {

	imageMetaData, err := p.storage.ImageStorage.GetById(ctx, id)

	if err != nil {

		return nil, "", err
	}

	key := imageMetaData.ObjectKey

	image, err := p.objectStore.Get(ctx, key)

	if err != nil {
		return nil, "", err
	}

	return image, imageMetaData.ContentType, nil
}

func (p *Processor) DeleteByImageId(ctx context.Context, id uuid.UUID) error {
	imageMetaData, err := p.storage.ImageStorage.GetById(ctx, id)

	if err != nil {
		return ErrorNotFound
	}

	err = p.storage.ImageStorage.DeleteById(ctx, id)

	if err != nil {
		return err
	}

	key := imageMetaData.ObjectKey

	err = p.objectStore.Delete(ctx, key)

	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) TranformImage(operation string) {
}
