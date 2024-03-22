// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package s3

import (
	"context"
	"fmt"
	"strconv"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	s3client "github.com/metal-stack/gardener-extension-backup-s3/pkg/s3/client"
)

// NewClientFromSecretRef creates a new Client for the given S3 credentials from given k8s <secretRef>
func NewClientFromSecretRef(ctx context.Context, client client.Client, secretRef corev1.SecretReference) (*s3client.Client, error) {
	credentials, err := getCredentialsFromSecretRef(ctx, client, secretRef)
	if err != nil {
		return nil, err
	}
	return s3client.NewClient(credentials)
}

func getCredentialsFromSecretRef(ctx context.Context, client client.Client, secretRef corev1.SecretReference) (*s3client.Credentials, error) {
	secret, err := extensionscontroller.GetSecretByReference(ctx, client, &secretRef)
	if err != nil {
		return nil, err
	}
	return readCredentialsSecret(secret)
}

func readCredentialsSecret(secret *corev1.Secret) (*s3client.Credentials, error) {
	if secret.Data == nil {
		return nil, fmt.Errorf("secret does not contain any data")
	}

	accessKeyID, err := getSecretStringValue(secret, AccessKeyID, true)
	if err != nil {
		return nil, err
	}

	secretAccessKey, err := getSecretStringValue(secret, SecretAccessKey, true)
	if err != nil {
		return nil, err
	}

	region, err := getSecretStringValue(secret, Region, true)
	if err != nil {
		return nil, err
	}

	endpoint, _ := getSecretStringValue(secret, Endpoint, false)

	insecureSkipVerify, err := getSecretBoolValue(secret, InsecureSkipVerify)
	if err != nil {
		return nil, err
	}

	s3ForcePathStyle, err := getSecretBoolValue(secret, S3ForcePathStyle)
	if err != nil {
		return nil, err
	}

	trustedCaCert, _ := getSecretStringValue(secret, TrustedCaCert, false)

	return &s3client.Credentials{
		AccessKeyID:        *accessKeyID,
		SecretAccessKey:    *secretAccessKey,
		Region:             *region,
		Endpoint:           endpoint,
		InsecureSkipVerify: insecureSkipVerify,
		S3ForcePathStyle:   s3ForcePathStyle,
		TrustedCaCert:      trustedCaCert,
	}, nil
}

func getSecretStringValue(secret *corev1.Secret, key string, required bool) (*string, error) {
	if value, ok := secret.Data[key]; ok {
		return ptr.To(string(value)), nil
	}
	if required {
		return nil, fmt.Errorf("missing %q field in secret", key)
	}
	return nil, nil
}

func getSecretBoolValue(secret *corev1.Secret, key string) (bool, error) {
	if raw, ok := secret.Data[key]; ok {
		value, err := strconv.ParseBool(string(raw))
		if err != nil {
			return false, fmt.Errorf("cannot parse %q field in secret: %w", key, err)
		}
		return value, nil
	}
	return false, nil
}
