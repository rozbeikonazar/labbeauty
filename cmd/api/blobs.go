package main

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureBlobStorage struct {
	client *azblob.Client
	ctx    context.Context
}

func NewAzureBlobStorage(blobURL string, credential azcore.TokenCredential, ctx context.Context) (*AzureBlobStorage, error) {
	client, err := azblob.NewClient(blobURL, credential, nil)
	if err != nil {
		return nil, err
	}
	return &AzureBlobStorage{client: client, ctx: ctx}, nil
}

func (abs *AzureBlobStorage) UploadBlob(blobName string, file *multipart.File) error {
	// Check if the blob with that name already exists
	buffer, err := io.ReadAll(*file)
	if err != nil {
		return err
	}

	for i := 1; i <= 3; i++ {
		_, err := abs.client.UploadBuffer(abs.ctx, containerName, blobName, buffer, &azblob.UploadBufferOptions{})
		if nil == err {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("failed to upload a blob after 3 attempts")

}

func (abs *AzureBlobStorage) DeleteBlob(blobName string) error {

	for i := 1; i <= 3; i++ {
		_, err := abs.client.DeleteBlob(abs.ctx, containerName, blobName, nil)
		if nil == err {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("failed to delete a blob after 3 attempts")

}
