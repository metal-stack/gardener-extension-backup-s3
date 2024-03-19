// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=github.com/metal-stack/gardener-extension-backup-s3/pkg/apis/config
// +k8s:openapi-gen=true
// +k8s:defaulter-gen=TypeMeta

// Package v1alpha1 contains the S3 backup provider configuration API resources.
// +groupName=s3.backup.extensions.config.gardener.cloud
package v1alpha1 // import "github.com/metal-stack/gardener-extension-backup-s3/pkg/apis/config/v1alpha1"
