package util

import (
	"bytes"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSCommon struct {
	Session *session.Session
}

var MyAWS *AWSCommon

func InitAwsSession(Region, accessKeyID, AccessSecret string) error {

	MyAWS = nil
	newSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, AccessSecret, ""),
	})
	if err != nil {
		return err
	}
	MyAWS = &AWSCommon{newSession}
	return nil
}

func UploadFileToS3(file *os.File, params map[string]interface{}) (string, error) {
	if MyAWS == nil {
		return "", errors.New("No config aws")
	}
	//bucket and key
	bucket := params["bucket_name"].(string)
	key := params["key"].(string)
	//file info
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	//upload with sdk
	_, err := s3.New(MyAWS.Session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return fileInfo.Name(), err
	}
	return fileInfo.Name(), nil
}
