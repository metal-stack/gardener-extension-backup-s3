// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package backupbucket

import (
	"context"

	"github.com/gardener/gardener/extensions/pkg/controller/backupbucket"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/metal-stack/gardener-extension-backup-s3/pkg/s3"
	"github.com/metal-stack/gardener-extension-backup-s3/pkg/s3/helper"
)

type actuator struct {
	backupbucket.Actuator
	client client.Client
}

func newActuator(mgr manager.Manager) backupbucket.Actuator {
	return &actuator{
		client: mgr.GetClient(),
	}
}

func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, backupBucket *extensionsv1alpha1.BackupBucket) error {
	s3Client, err := s3.NewClientFromSecretRef(ctx, a.client, backupBucket.Spec.SecretRef)
	if err != nil {
		log.Error(err, "unable to create S3 client")
		return util.DetermineError(err, helper.KnownCodes)
	}

	err = s3Client.CreateBucketIfNotExists(ctx, backupBucket.Name, backupBucket.Spec.Region)
	if err != nil {
		log.Error(err, "Failed to CreateBucketIfNotExists", "bucket", backupBucket.Name)
		return util.DetermineError(err, helper.KnownCodes)
	}

	return nil
}

func (a *actuator) Delete(ctx context.Context, _ logr.Logger, backupBucket *extensionsv1alpha1.BackupBucket) error {
	s3Client, err := s3.NewClientFromSecretRef(ctx, a.client, backupBucket.Spec.SecretRef)
	if err != nil {
		return util.DetermineError(err, helper.KnownCodes)
	}

	return util.DetermineError(s3Client.DeleteBucketIfExists(ctx, backupBucket.Name), helper.KnownCodes)
}
