package main

import (
	"context"
	"io"
	"mime/multipart"

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
	_, err = abs.client.UploadBuffer(abs.ctx, containerName, blobName, buffer, &azblob.UploadBufferOptions{})
	if err != nil {
		return err
	}

	return nil

}

func (abs *AzureBlobStorage) DeleteBlob(blobName string) error {
	_, err := abs.client.DeleteBlob(abs.ctx, containerName, blobName, nil)
	if err != nil {
		return err
	}
	return nil

}
