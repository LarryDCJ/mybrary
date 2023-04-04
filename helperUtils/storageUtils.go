package helperUtils

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"
)

const (
	bucketName = "roots_images" // FILL IN WITH YOURS
)

func InitStorageClient() *storage.Client {
	os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Panic(err)
	}

	defer func(client *storage.Client) {
		if err := client.Close(); err != nil {
			log.Print(err)
		}
	}(client)

	log.Printf("Storage connected\n")

	return client
}

// UploadFile uploads an object
func UploadFile(client *storage.Client, file *multipart.FileHeader) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucketName).Object(file.Filename)

	o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)

	var f multipart.File
	var err error

	if f, err = file.Open(); err != nil {
		return err
	}
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
