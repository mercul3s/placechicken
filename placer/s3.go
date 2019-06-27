package placer

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

// List returns all objects in an s3 bucket.
func (s *S3) List(b string) ([]Image, error) {
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

// Save stores an object in s3.
func (s *S3) Save() {

}

// s3 workflow for getting an image, resizing it, and returning to user:
// - get random image from the bucket (list all and do a random)
// - download it to a temp file
// - resize the image
// - upload it to s3
// - return the s3 url to the user (load image from s3)
