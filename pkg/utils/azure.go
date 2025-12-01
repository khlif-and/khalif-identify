package utils

import (
	"context"
	"mime/multipart"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

)

type AzureUploader struct {
	Client        *azblob.Client
	ContainerName string
}

func NewAzureUploader(connStr, containerName string) (*AzureUploader, error) {
	client, err := azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}
	return &AzureUploader{Client: client, ContainerName: containerName}, nil
}

func (a *AzureUploader) UploadFile(file multipart.File, filename string) (string, error) {
	ctx := context.Background()
	_, err := a.Client.UploadStream(ctx, a.ContainerName, filename, file, nil)
	if err != nil {
		return "", err
	}
	return "https://" + a.Client.URL() + "/" + a.ContainerName + "/" + filename, nil
}