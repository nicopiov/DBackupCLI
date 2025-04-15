/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/
package cmd

import (
	"bufio"
	"dbackupcli/cmd/commons"
	"dbackupcli/cmd/scripts"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// restoreCouchdbCmd represents the restoreCouchdb command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Operates a restore of a dump file in a database on CouchDB",
	Long:  `Operates a restore of a dump.json in the specified database on a certain CouchDB.`,
	Run: func(cmd *cobra.Command, args []string) {
		tmpFile, err := scripts.GetEmbeddedScripts()
		defer os.Remove(tmpFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		user, host, password, port := commons.GetAuthFlagValues(cmd)
		database, _ := cmd.Flags().GetString("database")
		file, _ := cmd.Flags().GetString("file")
		if commons.CheckFlags(append([]string{}, database, file, user, password, host)) {
			fmt.Println("missing on or more flags/arguments\n\nCheck using 'dbackupcli couchdb backup -h'")
			os.Exit(1)
		}

		statusCode, Database, err := commons.GetDB(host, port, user, password, database)
		if err != nil {
			fmt.Printf("%v", err)
		}

		if Database.DocCount != 0 {
			if statusCode == 200 {
				fmt.Printf("Database %s already exists. Do you want to overwrite it? (y/n): ", Database.DbName)
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if input != "y" && input != "Y" {
					fmt.Printf("operation canceled. The database will not be overwritten")
					return
				}
				fmt.Println("Overwriting database...")
				if err := commons.DeleteDatabase(host, port, user, password, database); err != nil {
					fmt.Printf("%v", err)
					return
				}
			} else {
				fmt.Printf("%v", err)
				return
			}
		}

		var cmdArgs []string = []string{"-r", "-d", database, "-f", file, "-c"}
		cmdArgs = commons.PrepareCmdAuthArgs(cmdArgs, user, password, host, port)

		cmdExec := exec.Command("bash", append([]string{tmpFile}, cmdArgs...)...)
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdout = os.Stdout
		if err := cmdExec.Run(); err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("Restore completed successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.SetUsageTemplate(`
Usage: dbackupcli {{.Use}} [flags]

Flags:
 -h, --help		Show this help message
 -d, --database		The database where to restore the dump
 -f, --file		The filename containing the dump to restore (e.g dump.json)
 -u, --user		The CouchDB username for the auth
 -p, --password		The CouchDB password for the auth
 --host			The host of the remote CouchDB, with or without 'https://'
 --port			The port of the remote CouchDB, default is 5984

Examples:
 dbackupcli restore -d my-db -f dump.json -u admin -p root --host 127.0.0.1 --port 9876
 dbackupcli restore --database my-db --file dump.json --user admin -p root --host 127.0.0.1 -c.
`)
	restoreCmd.Flags().BoolP("help", "h", false, "Help message")
	restoreCmd.Flags().String("host", "", "The remote CouchDB host, can be provided with or without 'https://'")
	restoreCmd.Flags().Int("port", 5984, "The remote CouchDB port (Default: 5984)")
	restoreCmd.Flags().StringP("database", "d", "", "The name of the database where to restore the dump (Default: empty)")
	restoreCmd.Flags().StringP("file", "f", "", "The name of the file containing the dump to restore (Default: empty)")
	restoreCmd.Flags().StringP("user", "u", "", "The user to authenticate to the CouchDB (Default: empty)")
	restoreCmd.Flags().StringP("password", "p", "", "The password to authenticate to the CouchDB (Default: empty)")
	restoreCmd.Flags().BoolP("createdb", "c", false, "Create the database if it does not exist on the remote couchdb (Default: false)")
}
