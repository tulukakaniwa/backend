package runner

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func createBucket(name string, compressedFlow string) (*s3.CreateBucketOutput, string, error) {
	// Create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Printf("error communicating with aws: %s", err.Error())
		return nil, "", err
	}
	svc := s3.New(sess)

	// Create bucket
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("openroad-cloud"),
	})
	if err != nil {
		log.Printf("error creating bucket: %s", err.Error())
		return nil, "", err
	}
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String("openroad-cloud"),
	})
	if err != nil {
		log.Printf("error while waiting for bucket creattion: %s", err.Error())
		return nil, "", err
	}

	// Upload current flow dir
	err = uploadFlowDir(name, compressedFlow)
	if err != nil {
		log.Printf("error while waiting for bucket creattion: %s", err.Error())
		return nil, "", err
	}

	// Generate URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("openroad-cloud"),
		Key:    aws.String(name + ".tar"),
	})
	urlString, err := req.Presign(7 * 24 * time.Hour)
	if err != nil {
		log.Printf("error signing the URL: %s", err.Error())
		return nil, "", err
	}
	return nil, urlString, nil
}

func uploadFlowDir(bucketName string, compressedFlow string) error {
	// Create session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Printf("error communicating with aws: %s", err.Error())
		return err
	}

	file, err := os.Open(compressedFlow)
	if err != nil {
		log.Printf("Unable to open flow file: %s", err.Error())
		return err
	}

	defer file.Close()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("openroad-cloud"),
		Key:    aws.String(bucketName + ".tar"),
		Body:   file,
	})
	if err != nil {
		log.Printf("error while uploading flow directory: %s", err.Error())
		return err
	}
	return nil
}
