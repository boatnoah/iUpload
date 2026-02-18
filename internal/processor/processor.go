package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/boatnoah/iupload/internal/storage"
	"github.com/disintegration/imaging"
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

type OperationPayload struct {
	Transformation *Transformations `json:"transformation"`
}

type Transformations struct {
	Resize      *Resize  `json:"resize"`
	Crop        *Crop    `json:"crop"`
	Rotate      *float64 `json:"rotate"`
	ContentType *string  `json:"content_type"`
}

type Resize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Crop struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
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

func (p *Processor) TranformImage(ctx context.Context, id uuid.UUID, transformPayload OperationPayload) ([]byte, error) {
	reader, contentType, err := p.GetByImageId(ctx, id)

	if err != nil {
		return nil, err
	}

	img, err := imaging.Decode(reader)

	if err != nil {
		return nil, err
	}

	if transformPayload.Transformation.Resize != nil {
		width := transformPayload.Transformation.Resize.Width
		height := transformPayload.Transformation.Resize.Height
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	if transformPayload.Transformation.Crop != nil {
		width := transformPayload.Transformation.Crop.Width
		height := transformPayload.Transformation.Crop.Height
		x := transformPayload.Transformation.Crop.X
		y := transformPayload.Transformation.Crop.Y

		err = validateValues(width, height, x, y, img.Bounds().Dx(), img.Bounds().Dy())
		if err != nil {
			return nil, err
		}

		rect := image.Rect(x, y, x+width, y+height)
		img = imaging.Crop(img, rect)
	}

	if transformPayload.Transformation.Rotate != nil {
		angle := transformPayload.Transformation.Rotate
		img = imaging.Rotate(img, *angle, color.Opaque)
	}

	if transformPayload.Transformation.ContentType != nil {

		err := validateContentType(*transformPayload.Transformation.ContentType)
		if err != nil {
			return nil, err
		}
		contentType = *transformPayload.Transformation.ContentType
	}

	var b bytes.Buffer
	var format imaging.Format

	if contentType == "image/png" {
		format = imaging.PNG
	} else {
		format = imaging.JPEG
	}

	err = imaging.Encode(&b, img, format)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func validateContentType(contentType string) error {

	if contentType != "image/png" && contentType != "image/jpeg" {
		return errors.New("Unable to convert to specified type")
	}

	return nil
}

func validateValues(width, height, x, y, boundsWidth, boundsHeight int) error {
	if x < 0 {
		return errors.New("Can't have x be negative")
	}

	if y < 0 {
		return errors.New("Can't have y be negative")
	}

	if x+width > boundsWidth {
		return errors.New("Can't exceed width")
	}

	if y+height > boundsHeight {
		return errors.New("Can't exceed height")
	}

	if width <= 0 {
		return errors.New("Can't have width be negative")
	}

	if height <= 0 {
		return errors.New("Can't have height be negative")
	}

	return nil

}
