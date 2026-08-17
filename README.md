# Gardener Extension for S3 Compatible Storage

[![GitHub License](https://img.shields.io/github/license/metal-stack/gardener-extension-backup-s3)](https://github.com/metal-stack/gardener-extension-backup-s3/blob/main/LICENSE)
[![Build](https://github.com/metal-stack/gardener-extension-backup-s3/actions/workflows/build.yaml/badge.svg)](https://github.com/metal-stack/gardener-extension-backup-s3/actions/workflows/build.yaml)

[Project Gardener](https://gardener.cloud/) implements the automated management and operation of [Kubernetes](https://kubernetes.io/) clusters as a service. This controller implements Gardener's extension contract for **S3 compatible storage**.

It reconciles the `BackupBucket` and `BackupEntry` resources of `type: S3` against any [S3](https://aws.amazon.com/s3/) compatible object store (e.g. AWS S3, MinIO, Ceph RGW, ONTAP...).

For more detailed documentation about the extension contract, please refer to the [Gardener docs](https://github.com/gardener/gardener/blob/master/docs/extensions/overview.md).

## Example

An example `ControllerRegistration` resource that can be used to register this controller to Gardener can be found [here](example/controller-registration.yaml).

## Development

The extension can be developed locally in the [mini-lab](https://github.com/metal-stack/mini-lab), which contains a working deployment of Minio.

## Feedback and Support

Feedback and contributions are always welcome! Please report bugs or suggestions as [GitHub issues](https://github.com/metal-stack/gardener-extension-backup-s3/issues) or reach out to our [community](https://metal-stack.io/community).
