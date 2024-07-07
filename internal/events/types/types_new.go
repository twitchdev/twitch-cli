// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package types

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type MockEventBase struct {
	Subscription map[string]interface{} `json:"subscription"`
	Event        map[string]interface{} `json:"event"`
}

var mockEvents []MockAbstract

func RegisterAllEvents() error {
	// Find directory holding YAML EventSub templates
	exeDir, err := os.Executable()
	if err != nil {
		return err
	}
	templatesBaseDir := path.Join(filepath.Dir(exeDir), "templates", "events")

	// Go through eventsYamlDir to find all files
	files := []string{}
	err = filepath.Walk(templatesBaseDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return errors.New("Could not read EventSub yaml files: " + err.Error())
	}

	// Read and store all the files
	for _, f := range files {
		abstract, err := ParseEventYaml(f)
		if err != nil {
			return err
		}

		// Check for duplicates
		ok, duplicateFilepath := RegisterSubscriptionType(abstract)
		if !ok {
			return errors.New("Duplicate subscription type/version pair:\n - " + abstract.Filepath + "\n - " + duplicateFilepath)
		}

	}

	return nil
}

func RegisterSubscriptionType(eventAbstract MockAbstract) (bool, string) {
	// Look for duplicates
	for _, sub := range mockEvents {
		if sub.Metadata.Type == eventAbstract.Metadata.Type && sub.Metadata.Version == eventAbstract.Metadata.Version {
			return false, sub.Filepath
		}
	}

	mockEvents = append(mockEvents, eventAbstract)

	return true, ""
}

func NEW_GetByTriggerAndTransportAndVersion(trigger string, transport string, version string) (*MockAbstract, error) {
	validEventBadVersions := []string{}
	var latestEventSeen *MockAbstract

	for _, sub := range mockEvents {
		if trigger == sub.Metadata.Type {
			// Found an event type that match's user input

			// Check if transport is valid
			validTransport := false
			for _, t := range sub.Metadata.SupportedTransports {
				if transport == t {
					validTransport = true
					break
				}
			}
			if !validTransport {
				if strings.EqualFold(transport, "websocket") {
					return nil, errors.New("Invalid transport. This event is not available via WebSockets.")
				}
				return nil, fmt.Errorf("Invalid transport. This event supports the following transport types: %v", strings.Join(sub.Metadata.SupportedTransports, ", "))
			}

			// Check for matching verison; Assumes version is not empty but doesn't matter performance-wise
			if version == sub.Metadata.Version {
				return &sub, nil
			} else {
				validEventBadVersions = append(validEventBadVersions, sub.Metadata.Version)
				latestEventSeen = &sub
			}
		}
	}

	// When no version is given, and there's only one version available, use the default version.
	if version == "" && len(validEventBadVersions) == 1 {
		return latestEventSeen, nil
	}

	// Error for events with non-existent version used
	if len(validEventBadVersions) != 0 {
		sort.Strings(validEventBadVersions)
		errStr := fmt.Sprintf("Invalid version given. Valid version(s): %v", strings.Join(validEventBadVersions, ", "))
		if version == "" {
			errStr += "\nUse --version to specify"
		}
		return nil, errors.New(errStr)
	}

	// Default error
	return nil, errors.New("Invalid event") // TODO
}

func GenerateEventObject() (map[string]MockAbstractData, error) {
	return nil, nil
}
