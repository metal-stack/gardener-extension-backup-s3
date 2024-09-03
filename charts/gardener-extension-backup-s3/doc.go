// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

//go:generate sh -c "bash $GARDENER_HACK_DIR/generate-controller-registration.sh backup-s3 . $(cat ../../VERSION) ../../example/controller-registration.yaml BackupBucket:S3 BackupEntry:S3"

// Package chart enables go:generate support for generating the correct controller registration.
package chart
