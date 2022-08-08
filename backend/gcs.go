package backend

import (
    "context"
    "fmt"
    "io"

    "around/constants"

    "cloud.google.com/go/storage"
)

var (
    GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
    client *storage.Client
    bucket string
}

func InitGCSBackend() {
    client, err := storage.NewClient(context.Background())
    if err != nil {
        panic(err)
    }

    GCSBackend = &GoogleCloudStorageBackend{
        client: client,
        bucket: constants.GCS_BUCKET,
    }
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
    //1.upload           
    //client --> GCS & ES
    //2. search 
    //client ---> ES(research by user/keyword to find the file link in GCS)
    
    ctx := context.Background()
    object := backend.client.Bucket(backend.bucket).Object(objectName)
    wc := object.NewWriter(ctx)
    //r is local source file
    if _, err := io.Copy(wc, r); err != nil {
        return "", err
    }

    if err := wc.Close(); err != nil {
        return "", err
    }

    //ACL: access control list
    //Set(ctx, storage.AllUsers, storage.RoleReader) modifies read authoritives that
    //all users can read becuase frontend needs to read (private->semi-public, can read can't write)
    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        return "", err
    }

    attrs, err := object.Attrs(ctx)
    if err != nil {
        return "", err
    }

    fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
}
