// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-FileCopyrightText: 2024 metal-stack Authors
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"regexp"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

var (
	unauthenticatedRegexp = regexp.MustCompile(`(?i)(AuthFailure|InvalidAccessKeyId|InvalidSecretAccessKey)`)
	unauthorizedRegexp    = regexp.MustCompile(`(?i)(Unauthorized|InvalidClientTokenId|SignatureDoesNotMatch|UnauthorizedOperation|AccessDenied)`)

	// KnownCodes maps Gardener error codes to respective regex.
	KnownCodes = map[gardencorev1beta1.ErrorCode]func(string) bool{
		gardencorev1beta1.ErrorInfraUnauthenticated: unauthenticatedRegexp.MatchString,
		gardencorev1beta1.ErrorInfraUnauthorized:    unauthorizedRegexp.MatchString,
	}
)
