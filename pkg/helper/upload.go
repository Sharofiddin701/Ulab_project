package helper

import (
	"context"
	"e-commerce/models"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

func UploadFiles(file *multipart.Form) (*models.MultipleFileUploadResponse, error) {

	i := 0
	var resp models.MultipleFileUploadResponse

	for _, v := range file.File["file"] {

		id := uuid.New().String()
		var (
			url models.Url
		)

		url.Url = fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/ecommece-e1b2e.appspot.com/o/%s?alt=media&token=%s", id, id)

		resp.Url = append(resp.Url, &url)
		resp.Url[i].Id = id

		imageFile, err := v.Open()
		if err != nil {
			return nil, err
		}
		tempFile, err := os.Create(id)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(tempFile, imageFile)
		if err != nil {
			return nil, err
		}

		defer tempFile.Close()

		opt := option.WithCredentialsFile("serviceAccountKey.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		client, err := app.Storage(context.TODO())
		if err != nil {
			return nil, err
		}

		bucketHandle, err := client.Bucket("ecommece-e1b2e.appspot.com")
		if err != nil {
			return nil, err
		}

		f, err := os.Open(tempFile.Name())
		if err != nil {
			return nil, err
		}
		defer os.Remove(tempFile.Name())
		objectHandle := bucketHandle.Object(tempFile.Name())

		writer := objectHandle.NewWriter(context.Background())

		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		defer writer.Close()

		if _, err := io.Copy(writer, f); err != nil {
			return nil, err
		}

		defer os.Remove(tempFile.Name())
		i++
	}

	return &resp, nil
}

func DeleteFile(id string) error {
	// Initialize a context and Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Bucket name and object path to delete

	fmt.Println(id)
	bucketName := "ecommece-e1b2e.appspot.com"
	objectPath := id

	// Delete the object
	err = client.Objects.Delete(bucketName, objectPath).Do()
	if err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}

	fmt.Printf("Object %s deleted successfully from bucket %s\n", objectPath, bucketName)
	return nil
}

func UploadFile(file *os.File) (*models.MultipleFileUploadResponse, error) {

	i := 0
	var resp models.MultipleFileUploadResponse

	id := uuid.New().String()
	var (
		url models.Url
	)

	url.Url = fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/ecommece-e1b2e.appspot.com/o/%s?alt=media&token=%s", id, id)

	resp.Url = append(resp.Url, &url)
	resp.Url[i].Id = id

	imageFile, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}
	tempFile, err := os.Create(id)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(tempFile, imageFile)
	if err != nil {
		return nil, err
	}

	defer tempFile.Close()

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	client, err := app.Storage(context.TODO())
	if err != nil {
		return nil, err
	}

	bucketHandle, err := client.Bucket("ecommece-e1b2e.appspot.com")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())
	objectHandle := bucketHandle.Object(tempFile.Name())

	writer := objectHandle.NewWriter(context.Background())

	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

	defer writer.Close()

	if _, err := io.Copy(writer, f); err != nil {
		return nil, err
	}

	defer os.Remove(tempFile.Name())

	resp.Url[0].Id = id

	return &resp, nil
}
