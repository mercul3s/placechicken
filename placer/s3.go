package placer

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 struct exists as a placeholder to allow abstracting aws' s3 methods.
type S3 struct {
	Session *session.Session
}

// S3Config returns an s3 config populated with a session.
func S3Config() (S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return S3{}, err
	}
	return S3{Session: session}, err
}

// RandImg gets a list of objects in s3, and returns a random object for image
// manipulation.
func (s *S3) RandImg(b string) (Image, error) {
	var fileName string
	objects, err := s.list(b)
	i := Image{}
	if err != nil {
		return i, err
	}

	if len(objects) > 0 {
		randIdx := rand.Intn(len(objects))
		randObject := objects[randIdx]
		fileName, err = s.download(randObject, b)
		if err != nil {
			return i, err
		}
	}
	i.Name = fileName
	return i, nil
}

func (s *S3) list(b string) ([]Image, error) {
	i := []Image{}
	svc := s3.New(s.Session)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(b)})
	if err != nil {
		return i, err
	}
	for _, object := range resp.Contents {
		if strings.Contains(*object.Key, "original") {
			i = append(i, Image{Name: *object.Key})
		}
	}
	return i, err
}

func (s *S3) download(i Image, b string) (string, error) {
	downloader := s3manager.NewDownloader(s.Session)
	file, err := os.Create("/tmp/placechicken/" + i.Name)
	if err != nil {
		return "", err
	}
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(b),
			Key:    aws.String(i.Name),
		})
	if err != nil {
		return "", err
	}

	fmt.Printf("object %s downloaded to %s file size %d", i.Name, file.Name(), numBytes)
	return file.Name(), nil

}

func (s *S3) save() {

}

// s3 workflow for getting an image, resizing it, and returning to user:
// - get random image from the bucket (list all and do a random)
// - download it to a temp file
// - resize the image
// - upload it to s3
// - return the s3 url to the user (load image from s3)
