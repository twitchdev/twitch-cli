// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"sort"
)

func ValidTransports() []string {
	names := []string{}

	for name, enabled := range transportSupported {
		if enabled == true {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	return names
}
