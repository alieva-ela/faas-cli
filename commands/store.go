// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"text/tabwriter"

	"github.com/spf13/cobra"
)

var storeAddress string

const defaultStore = "https://cdn.rawgit.com/openfaas/store/master/store.json"

type storeItem struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	Name        string            `json:"name"`
	Fprocess    string            `json:"fprocess"`
	Network     string            `json:"network"`
	RepoURL     string            `json:"repo_url"`
	Environment map[string]string `json:"environment"`
}

func init() {
	// Setup flags that are used by multiple commands (variables defined in faas.go)
	storeListCmd.Flags().StringVarP(&storeAddress, "store", "g", defaultStore, "Store URL starting with http(s)://")
	storeInspectCmd.Flags().StringVarP(&storeAddress, "store", "g", defaultStore, "Store URL starting with http(s)://")
	storeDeployCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL starting with http(s)://")
	storeDeployCmd.Flags().StringVar(&handler, "handler", "", "Directory with handler for function, e.g. handler.js")
	storeDeployCmd.Flags().StringVar(&language, "lang", "", "Programming language template")

	// Setup flags that are used only by deploy command (variables defined above)
	storeDeployCmd.Flags().StringArrayVarP(&envvarOpts, "env", "e", []string{}, "Set one or more environment variables (ENVVAR=VALUE)")
	storeDeployCmd.Flags().StringArrayVarP(&labelOpts, "label", "l", []string{}, "Set one or more label (LABEL=VALUE)")
	storeDeployCmd.Flags().BoolVar(&replace, "replace", true, "Replace any existing function")
	storeDeployCmd.Flags().BoolVar(&update, "update", false, "Update existing functions")
	storeDeployCmd.Flags().StringArrayVar(&constraints, "constraint", []string{}, "Apply a constraint to the function")
	storeDeployCmd.Flags().StringArrayVar(&secrets, "secret", []string{}, "Give the function access to a secure secret")

	// Set bash-completion.
	_ = storeDeployCmd.Flags().SetAnnotation("handler", cobra.BashCompSubdirsInDir, []string{})

	storeCmd.AddCommand(storeListCmd)
	storeCmd.AddCommand(storeInspectCmd)
	storeCmd.AddCommand(storeDeployCmd)
	faasCmd.AddCommand(storeCmd)
}

var storeCmd = &cobra.Command{
	Use:   `store`,
	Short: "OpenFaaS store commands",
	Long:  "Allows browsing and deploying OpenFaaS store functions",
}

var storeListCmd = &cobra.Command{
	Use:     `list [--store STORE_URL]`,
	Short:   "List OpenFaaS store items",
	Long:    "Lists the available items in OpenFaas store",
	Example: `  faas-cli store list --store https://domain:port`,
	RunE:    runStoreList,
}

var storeInspectCmd = &cobra.Command{
	Use:     `inspect FUNCTION_NAME [--store STORE_URL]`,
	Short:   "Show OpenFaaS store function details",
	Long:    "Prints the detailed informations of the specified OpenFaaS function",
	Example: `  faas-cli store inspect NodeInfo --store https://domain:port`,
	RunE:    runStoreInspect,
}

var storeDeployCmd = &cobra.Command{
	Use: `deploy FUNCTION_NAME
							[--lang <ruby|python|node|csharp>]
							[--gateway GATEWAY_URL]
							[--handler HANDLER_DIR]
							[--env ENVVAR=VALUE ...]
							[--label LABEL=VALUE ...]
							[--replace=false]
							[--update=false]
							[--constraint PLACEMENT_CONSTRAINT ...]
							[--regex "REGEX"]
							[--filter "WILDCARD"]
							[--secret "SECRET_NAME"]`,

	Short: "Deploy OpenFaaS functions from the store",
	Long:  `Same as faas-cli deploy except pre-loaded with arguments from the store`,
	Example: `  faas-cli store deploy figlet
  faas-cli store deploy figlet --label canary=true
  faas-cli store deploy figlet --filter "*gif*" --secret dockerhuborg
  faas-cli store deploy figlet --regex "fn[0-9]_.*"
  faas-cli store deploy figlet --replace=false
  faas-cli store deploy figlet --update=true
  faas-cli deploy --image=alexellis/faas-url-ping --name=url-ping
  faas-cli store deploy figlet --handler=/path/to/fn/
                  --gateway=http://remote-site.com:8080 --lang=python
                  --env=MYVAR=myval`,
	RunE: runStoreDeploy,
}

func runStoreList(cmd *cobra.Command, args []string) error {
	items, err := storeList(storeAddress)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Printf("The store is empty.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "FUNCTION\tDESCRIPTION")

	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\n", item.Title, item.Description)
	}

	fmt.Fprintln(w)
	w.Flush()

	return nil
}

func runStoreInspect(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide the function name")
	}

	item, err := findFunction(args[0])
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "FUNCTION\tDESCRIPTION\tIMAGE\tFUNCTION PROCESS\tREPO")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
		item.Title,
		item.Description,
		item.Image,
		item.Fprocess,
		item.RepoURL,
	)

	fmt.Fprintln(w)
	w.Flush()
	return nil
}

func runStoreDeploy(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide the function name")
	}

	item, err := findFunction(args[0])
	if err != nil {
		return err
	}

	// Add the store environement variables to the provided ones from cmd
	if item.Environment != nil {
		for _, env := range item.Environment {
			envvarOpts = append(envvarOpts, env)
		}
	}

	RunDeploy(
		cmd,
		args,
		item.Image,
		item.Fprocess,
		envvarOpts,
	)
	return nil
}

func storeList(store string) ([]storeItem, error) {
	var results []storeItem

	store = strings.TrimRight(store, "/")

	timeout := 60 * time.Second
	client := makeHTTPClient(&timeout)

	getRequest, err := http.NewRequest(http.MethodGet, store, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to OpenFaaS store on URL: %s", store)
	}

	res, err := client.Do(getRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to OpenFaaS store on URL: %s", store)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK:

		bytesOut, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read result from OpenFaaS store on URL: %s", store)
		}
		jsonErr := json.Unmarshal(bytesOut, &results)
		if jsonErr != nil {
			return nil, fmt.Errorf("cannot parse result from OpenFaaS store on URL: %s\n%s", store, jsonErr.Error())
		}
	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return nil, fmt.Errorf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}
	return results, nil
}

func findFunction(functionName string) (storeItem, error) {
	var item storeItem

	items, err := storeList(storeAddress)
	if err != nil {
		return item, err
	}

	for _, item = range items {
		if item.Name == functionName {
			return item, nil
		}
	}

	return item, fmt.Errorf("function '%s' not found", functionName)
}

func makeHTTPClient(timeout *time.Duration) http.Client {
	if timeout != nil {
		return http.Client{
			Timeout: *timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout: *timeout,
					// KeepAlive: 0,
				}).DialContext,
				// MaxIdleConns:          1,
				// DisableKeepAlives:     true,
				IdleConnTimeout:       120 * time.Millisecond,
				ExpectContinueTimeout: 1500 * time.Millisecond,
			},
		}
	}

	// This should be used for faas-cli invoke etc.
	return http.Client{}
}
