/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/

package cmd

import (
	"dbackupcli/cmd/commons"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listCouchdbCmd represents the listCouchdb command
var listCmd = &cobra.Command{
	Use:   "listdbs",
	Short: "List the databases contained in the specified CouchDB",
	Long:  `List the databases contained in the specified CouchDB`,
	Run: func(cmd *cobra.Command, args []string) {
		user, host, password, port := commons.GetAuthFlagValues(cmd)
		if commons.CheckFlags(append([]string{}, user, password, host)) {
			fmt.Println("missing on or more flags/arguments\n\nCheck using 'dbackupcli backup -h'")
			os.Exit(1)
		}

		dbNames, err := commons.GetDBs(host, port, user, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		if dbNames != nil {
			fmt.Printf("Found %d databases\n", len(dbNames))
			fmt.Println("List of databases: \n", dbNames)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.SetUsageTemplate(`
Usage: dbackupcli {{.Use}} [flags]

Flags:
 -h, --help		Show this help message
 -u, --user		The CouchDB username for the auth
 -p, --password		The CouchDB password for the auth
 --host			The host of the remote CouchDB, with or without 'https://'
 --port			The port of the remote CouchDB, default is 5984

Examples:
 dbackupcli listdbs -u admin -p root --host 127.0.0.1 --port 9876
 dbackupcli listdbs --user admin -p root --host 127.0.0.1.
`)
	listCmd.Flags().BoolP("help", "h", false, "Help message")
	listCmd.Flags().String("host", "", "The remote CouchDB host, can be provided with or without 'https://'")
	listCmd.Flags().Int("port", 5984, "The remote CouchDB port (Default: 5984)")
	listCmd.Flags().StringP("user", "u", "", "The user to authenticate to the CouchDB (Default: empty)")
	listCmd.Flags().StringP("password", "p", "", "The password to authenticate to the CouchDB (Default: empty)")
}
