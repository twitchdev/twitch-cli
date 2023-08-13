// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/generate"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_server"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var queryParameters []string
var body string
var prettyPrint bool
var autoPaginate int = 0
var port int
var verbose bool

var generateCount int

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Used to interface with the Twitch API",
}

var getCmd = &cobra.Command{
	Use:       "get",
	Short:     "Performs a GET request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("GET"),
	RunE:      cmdRun,
}
var postCmd = &cobra.Command{
	Use:       "post",
	Short:     "Performs a POST request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("POST"),
	RunE:      cmdRun,
}
var patchCmd = &cobra.Command{
	Use:       "patch",
	Short:     "Performs a PATCH request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("PATCH"),
	RunE:      cmdRun,
}
var deleteCmd = &cobra.Command{
	Use:       "delete",
	Short:     "Performs a DELETE request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("DELETE"),
	RunE:      cmdRun,
}
var putCmd = &cobra.Command{
	Use:       "put",
	Short:     "Performs a PUT request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("PUT"),
	RunE:      cmdRun,
}

var mockCmd = &cobra.Command{
	Use:   "mock-api",
	Short: "Used to interface with the mock Twitch API.",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Used to start the server for the mock API.",
	RunE:  mockStartRun,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Used to randomly generate data for use with the mock API. By default, this is run on the first invocation of the start command, however this allows you to generate further primitives.",
	RunE:  generateMockRun,
}

func init() {
	rootCmd.AddCommand(apiCmd, mockCmd)

	apiCmd.AddCommand(getCmd, postCmd, patchCmd, deleteCmd, putCmd)

	apiCmd.PersistentFlags().StringArrayVarP(&queryParameters, "query-params", "q", nil, "Available multiple times. Passes in query parameters to endpoints using the format of `key=value`.")
	apiCmd.PersistentFlags().StringVarP(&body, "body", "b", "", "Passes a body to the request. Alteratively supports CURL-like references to files using the format of `@data,json`.")
	apiCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Whether to display HTTP request and header information above the response of the API call.")

	// default here is false to enable -p commands to toggle off without explicitly defining -p=false as -p false will not work. The below commands invert the bool to pass the true default. Deprecated, so marking as hidden in favor of the unformatted flag.
	apiCmd.PersistentFlags().BoolVarP(&prettyPrint, "pretty-print", "p", false, "Whether to pretty-print API requests. Default is true.")
	apiCmd.PersistentFlags().MarkHidden("pretty-print")

	apiCmd.PersistentFlags().BoolVarP(&prettyPrint, "unformatted", "u", false, "Whether to have API requests come back unformatted/non-prettyprinted. Default is false.")

	getCmd.PersistentFlags().IntVarP(&autoPaginate, "autopaginate", "P", 0, "Whether to have API requests automatically paginate. Default is to not paginate.")
	getCmd.PersistentFlags().Lookup("autopaginate").NoOptDefVal = "0"

	mockCmd.AddCommand(startCmd, generateCmd)

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Defines the port that the mock API will run on.")

	generateCmd.Flags().IntVarP(&generateCount, "count", "c", 25, "Defines the number of fake users to generate.")
}

func cmdRun(cmd *cobra.Command, args []string) error {
	var path string

	if len(args) == 0 {
		cmd.Help()
		return fmt.Errorf("")
	} else if len(args) == 1 && args[0][:1] == "/" {
		path = args[0]
	} else {
		path = "/" + strings.Join(args[:], "/")
	}

	if body != "" && body[:1] == "@" {
		var err error
		body, err = getBodyFromFile(body[1:])
		if err != nil {
			return err
		}
	}

	if cmd.Name() == "get" && cmd.PersistentFlags().Lookup("autopaginate").Changed {
		return api.NewRequest(cmd.Name(), path, queryParameters, []byte(body), !prettyPrint, &autoPaginate, verbose)
	} else {
		return api.NewRequest(cmd.Name(), path, queryParameters, []byte(body), !prettyPrint, nil, verbose) // only set on when the user changed the flag
	}
}

func getBodyFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func mockStartRun(cmd *cobra.Command, args []string) error {
	log.Printf("Starting mock API server on http://localhost:%v", port)
	return mock_server.StartServer(port)
}

func generateMockRun(cmd *cobra.Command, args []string) error {
	generate.Generate(generateCount)
	return nil
}
