// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
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

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Used to interface with the Twitch API",
}

var getCmd = &cobra.Command{
	Use:       "get",
	Short:     "Performs a GET request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("GET"),
	Run:       cmdRun,
}
var postCmd = &cobra.Command{
	Use:       "post",
	Short:     "Performs a POST request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("POST"),
	Run:       cmdRun,
}
var patchCmd = &cobra.Command{
	Use:       "patch",
	Short:     "Performs a PATCH request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("PATCH"),
	Run:       cmdRun,
}
var deleteCmd = &cobra.Command{
	Use:       "delete",
	Short:     "Performs a DELETE request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("DELETE"),
	Run:       cmdRun,
}
var putCmd = &cobra.Command{
	Use:       "put",
	Short:     "Performs a PUT request on the specified command.",
	Args:      cobra.MaximumNArgs(3),
	ValidArgs: api.ValidOptions("PUT"),
	Run:       cmdRun,
}

var mockCmd = &cobra.Command{
	Use:   "mock-api",
	Short: "Used to interface with the mock Twitch API.",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Used to start the server for the mock API.",
	Run:   mockStartRun,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Used to randomly generate data for use with the mock API. By default, this is run on the first invocation of the start command, however this allows you to generate further primitives.",
	Run:   generateMockRun,
}

func init() {
	rootCmd.AddCommand(apiCmd, mockCmd)

	apiCmd.AddCommand(getCmd, postCmd, patchCmd, deleteCmd, putCmd)

	apiCmd.PersistentFlags().StringArrayVarP(&queryParameters, "query-params", "q", nil, "Available multiple times. Passes in query parameters to endpoints using the format of `key=value`.")
	apiCmd.PersistentFlags().StringVarP(&body, "body", "b", "", "Passes a body to the request. Alteratively supports CURL-like references to files using the format of `@data,json`.")

	// default here is false to enable -p commands to toggle off without explicitly defining -p=false as -p false will not work. The below commands invert the bool to pass the true default. Deprecated, so marking as hidden in favor of the unformatted flag.
	apiCmd.PersistentFlags().BoolVarP(&prettyPrint, "pretty-print", "p", false, "Whether to pretty-print API requests. Default is true.")
	apiCmd.PersistentFlags().MarkHidden("pretty-print")

	apiCmd.PersistentFlags().BoolVarP(&prettyPrint, "unformatted", "u", false, "Whether to have API requests come back unformatted/non-prettyprinted. Default is false.")

	getCmd.PersistentFlags().IntVarP(&autoPaginate, "autopaginate", "P", 0, "Whether to have API requests automatically paginate. Default is to not paginate.")
	getCmd.PersistentFlags().Lookup("autopaginate").NoOptDefVal = "0"

	mockCmd.AddCommand(startCmd, generateCmd)

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Defines the port that the mock API will run on.")

	generateCmd.Flags().IntVarP(&count, "count", "c", 10, "Defines the number of fake users to generate.")
}

func cmdRun(cmd *cobra.Command, args []string) {
	var path string

	if len(args) == 0 {
		cmd.Help()
		return
	} else if len(args) == 1 && args[0][:1] == "/" {
		path = args[0]
	} else {
		path = "/" + strings.Join(args[:], "/")
	}

	if body != "" && body[:1] == "@" {
		body = getBodyFromFile(body[1:])
	}

	if cmd.PersistentFlags().Lookup("autopaginate").Changed {
		api.NewRequest(cmd.Name(), path, queryParameters, []byte(body), !prettyPrint, &autoPaginate)
	} else {
		api.NewRequest(cmd.Name(), path, queryParameters, []byte(body), !prettyPrint, nil) // only set on when the user changed the flag
	}

}

func getBodyFromFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func mockStartRun(cmd *cobra.Command, args []string) {
	log.Printf("Starting mock API server on http://localhost:%v", port)
	mock_server.StartServer(port)
}

func generateMockRun(cmd *cobra.Command, args []string) {
	generate.Generate(count)
}
