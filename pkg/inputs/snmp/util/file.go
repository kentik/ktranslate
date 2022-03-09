package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
		return loadFromS3(ctx, u, getS3Downloader())
	case "s3m":
		return loadFromS3(ctx, u, getMockS3Client())
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
		return writeToHttp(ctx, file, payload)
	case "s3":
		return writeToS3(ctx, u, payload, getS3Uploader())
	case "s3m":
		return writeToS3(ctx, u, payload, getMockS3Client())
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

func writeToHttp(ctx context.Context, file string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, file, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Cannot send to %s: %d", file, resp.StatusCode)
	}
	return nil
}

func loadFromS3(ctx context.Context, url *url.URL, client s3ClientDown) ([]byte, error) {
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

func writeToS3(ctx context.Context, url *url.URL, payload []byte, client s3ClientUp) error {
	_, err := client.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(url.Host),
		Body:   bytes.NewBuffer(payload),
		Key:    aws.String(url.Path),
	})
	return err
}

type s3ClientDown interface {
	DownloadWithContext(aws.Context, io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
}

type s3ClientUp interface {
	UploadWithContext(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type s3Full interface {
	s3ClientUp
	s3ClientDown
}

func getS3Downloader() s3ClientDown {
	sess := session.Must(session.NewSession())
	client := s3manager.NewDownloader(sess)
	return client
}

func getS3Uploader() s3ClientUp {
	sess := session.Must(session.NewSession())
	client := s3manager.NewUploader(sess)
	return client
}

type mockS3 struct {
	lastContent []byte
}

var (
	sMock *mockS3
)

func (m *mockS3) UploadWithContext(ctx aws.Context, in *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	m.lastContent = in.Body.(*bytes.Buffer).Bytes()
	return nil, nil
}

func (m *mockS3) DownloadWithContext(ctx aws.Context, w io.WriterAt, in *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	w.WriteAt(m.lastContent, 0)
	return int64(len(m.lastContent)), nil
}

func getMockS3Client() s3Full {
	return sMock
}
