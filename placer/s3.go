package placer

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 struct exists as a placeholder to allow abstracting aws' s3 methods.
type S3 struct {
}

// List returns all objects in an s3 bucket.
func (s *S3) List(p string) ([]os.FileInfo, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		fmt.Println("unable to create session")
		return nil, err
	}
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	return resp, err
}
