package util

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Utility to load a config file from various places
func LoadFile(ctx context.Context, file string) ([]byte, error) {
	u, err := url.Parse(file)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		return loadFromHttp(ctx, file)
	case "s3":
		return loadFromS3(ctx, u)
	default:
		return ioutil.ReadFile(file)
	}
}

// Utility to write a config file to various places
func WriteFile(ctx context.Context, file string, payload []byte, perms fs.FileMode) error {
	u, err := url.Parse(file)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https":
		return fmt.Errorf("Not supported scheme: %s", u.Scheme)
	case "s3":
		return writeToS3(ctx, u, payload)
	default:
		return ioutil.WriteFile(file, payload, perms)
	}
}

func loadFromHttp(ctx context.Context, file string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func loadFromS3(ctx context.Context, url *url.URL) ([]byte, error) {
	sess := session.Must(session.NewSession())
	client := s3manager.NewDownloader(sess)
	buf := aws.NewWriteAtBuffer([]byte{})
	size, err := client.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(url.Host),
		Key:    aws.String(url.Path),
	})
	if err != nil {
		return nil, err
	}

	bufr := buf.Bytes()
	return bufr[0:size], nil
}

func writeToS3(ctx context.Context, url *url.URL, payload []byte) error {
	sess := session.Must(session.NewSession())
	client := s3manager.NewUploader(sess)
	_, err := client.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(url.Host),
		Body:   bytes.NewBuffer(payload),
		Key:    aws.String(url.Path),
	})
	return err
}
