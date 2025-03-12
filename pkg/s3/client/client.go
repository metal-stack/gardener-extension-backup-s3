// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewClient returns a new S3 client from the given credentials.
func NewClient(cred *Credentials) (*Client, error) {
	httpClient := http.DefaultClient

	if cred.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if cred.TrustedCaCert != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(*cred.TrustedCaCert))
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			},
		}
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials:      credentials.NewStaticCredentials(cred.AccessKeyID, cred.SecretAccessKey, ""),
			Endpoint:         cred.Endpoint,
			Region:           aws.String(cred.Region),
			S3ForcePathStyle: aws.Bool(cred.S3ForcePathStyle),
			HTTPClient:       httpClient,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		s3: s3.New(sess),
	}, nil
}

// DeleteObjectsWithPrefix deletes the s3 objects with the specific <prefix> from <bucket>. If it does not exist,
// no error is returned.
func (c *Client) DeleteObjectsWithPrefix(ctx context.Context, bucket, prefix string) error {
	in := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	var delErr error
	if err := c.s3.ListObjectsPagesWithContext(ctx, in, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		objectIDs := make([]*s3.ObjectIdentifier, 0)
		for _, key := range page.Contents {
			obj := &s3.ObjectIdentifier{
				Key: key.Key,
			}
			objectIDs = append(objectIDs, obj)
		}

		if len(objectIDs) != 0 {
			if _, delErr = c.s3.DeleteObjectsWithContext(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &s3.Delete{
					Objects: objectIDs,
					Quiet:   aws.Bool(true),
				},
			}); delErr != nil {
				return false
			}
		}
		return !lastPage
	}); err != nil {
		return err
	}

	if delErr != nil {
		if aerr, ok := delErr.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil
		}
		return delErr
	}
	return nil
}

// CreateBucketIfNotExists creates the s3 bucket with name <bucket> in <region>. If it already exists,
// no error is returned.
func (c *Client) CreateBucketIfNotExists(ctx context.Context, bucket, region string) error {
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}

	if _, err := c.s3.HeadBucket(input); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				return c.CreateBucketIfNotExists2(ctx, bucket, region)
			}
		}
		return err
	}
	return nil
}

func (c *Client) CreateBucketIfNotExists2(ctx context.Context, bucket, region string) error {
	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(s3.BucketCannedACLPrivate),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	}

	if _, err := c.s3.CreateBucketWithContext(ctx, createBucketInput); err != nil {
		if aerr, ok := err.(awserr.Error); !ok {
			return err
		} else if aerr.Code() != s3.ErrCodeBucketAlreadyExists && aerr.Code() != s3.ErrCodeBucketAlreadyOwnedByYou {
			return err
		}
	}

	return nil
}

// DeleteBucketIfExists deletes the s3 bucket with name <bucket>. If it does not exist,
// no error is returned.
func (c *Client) DeleteBucketIfExists(ctx context.Context, bucket string) error {
	if _, err := c.s3.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				return nil
			}
			if aerr.Code() == "BucketNotEmpty" {
				if err := c.DeleteObjectsWithPrefix(ctx, bucket, ""); err != nil {
					return err
				}
				return c.DeleteBucketIfExists(ctx, bucket)
			}
		}
		return err
	}
	return nil
}
