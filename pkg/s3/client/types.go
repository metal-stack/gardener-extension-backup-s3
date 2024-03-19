// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Client struct {
	s3 s3iface.S3API
}

type Credentials struct {
	AccessKeyID        string
	Region             string
	SecretAccessKey    string
	Endpoint           *string
	InsecureSkipVerify bool
	S3ForcePathStyle   bool
	TrustedCaCert      *string
}
