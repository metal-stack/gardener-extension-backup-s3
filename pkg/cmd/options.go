// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	extensionsbackupbucketcontroller "github.com/gardener/gardener/extensions/pkg/controller/backupbucket"
	extensionsbackupentrycontroller "github.com/gardener/gardener/extensions/pkg/controller/backupentry"
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	extensionsheartbeatcontroller "github.com/gardener/gardener/extensions/pkg/controller/heartbeat"

	backupbucketcontroller "github.com/metal-stack/gardener-extension-backup-s3/pkg/controller/backupbucket"
	backupentrycontroller "github.com/metal-stack/gardener-extension-backup-s3/pkg/controller/backupentry"
)

// ControllerSwitchOptions are the controllercmd.SwitchOptions for the provider controllers.
func ControllerSwitchOptions() *controllercmd.SwitchOptions {
	return controllercmd.NewSwitchOptions(
		controllercmd.Switch(extensionsbackupbucketcontroller.ControllerName, backupbucketcontroller.AddToManager),
		controllercmd.Switch(extensionsbackupentrycontroller.ControllerName, backupentrycontroller.AddToManager),
		controllercmd.Switch(extensionsheartbeatcontroller.ControllerName, extensionsheartbeatcontroller.AddToManager),
	)
}
