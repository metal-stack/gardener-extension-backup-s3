// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/gardener/gardener/pkg/logger"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/metal-stack/gardener-extension-backup-s3/cmd/gardener-extension-backup-s3/app"
)

func main() {
	runtimelog.SetLogger(logger.MustNewZapLogger(logger.InfoLevel, logger.FormatJSON))
	cmd := app.NewControllerManagerCommand(signals.SetupSignalHandler())

	if err := cmd.Execute(); err != nil {
		runtimelog.Log.Error(err, "Error executing the main controller command")
		os.Exit(1)
	}
}
