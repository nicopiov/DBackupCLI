/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/
package cmd

import (
	"dbackupcli/cmd/commons"
	"dbackupcli/cmd/scripts"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// restoreAllCmd represents the restoreAll command
var restoreAllCmd = &cobra.Command{
	Use:   "restoreAll",
	Short: "Operates a restore of a dump of an entire CouchDB istance",
	Long:  `Operates a restore of a dump of an entire CouchDB istance contained inside a directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		tmpFile, err := scripts.GetEmbeddedScripts()
		defer os.Remove(tmpFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		user, host, password, port := commons.GetAuthFlagValues(cmd)
		dir, _ := cmd.Flags().GetString("filedir")
		if commons.CheckFlags(append([]string{}, dir, user, password, host)) {
			fmt.Println("missing on or more flags/arguments\n\nCheck using 'dbackupcli restoreAll -h'")
			os.Exit(1)
		}

		var cmdArgs []string = []string{"-r", "-c"}
		cmdArgs = commons.PrepareCmdAuthArgs(cmdArgs, user, password, host, port)
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		for _, file := range files {
			dbName, _ := strings.CutSuffix(file.Name(), ".json")
			restoreArgs := append(cmdArgs, "-d", dbName, "-f", dir+"/"+file.Name())
			cmdExec := exec.Command("bash", append([]string{tmpFile}, restoreArgs...)...)
			cmdExec.Stderr = os.Stderr
			cmdExec.Stdout = os.Stdout
			if err := cmdExec.Run(); err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("Restore completed successfully!")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(restoreAllCmd)
	restoreAllCmd.SetUsageTemplate(`
Usage: dbackupcli {{.Use}} [flags]

Flags:
 -h, --help		Show this help message
 -d, --database		The database where to restore the dump
 -f, --filedir		The directory containing the istance's dump to restore (e.g dump-directory)
 -u, --user		The CouchDB username for the auth
 -p, --password		The CouchDB password for the auth
 --host			The host of the remote CouchDB, with or without 'https://'
 --port			The port of the remote CouchDB, default is 5984

Examples:
 dbackupcli restore -d my-db -f backup_dir -u admin -p root --host 127.0.0.1 --port 9876
 dbackupcli restore --database my-db --filedir backup_dir --user admin -p root --host 127.0.0.1 -c.
`)
	restoreAllCmd.Flags().BoolP("help", "h", false, "Help message")
	restoreAllCmd.Flags().String("host", "", "The remote CouchDB host, can be provided with or without 'https://'")
	restoreAllCmd.Flags().Int("port", 5984, "The remote CouchDB port (Default: 5984)")
	restoreAllCmd.Flags().StringP("filedir", "f", "", "The name of the directory containing the istance's dump to restore (Default: empty)")
	restoreAllCmd.Flags().StringP("user", "u", "", "The user to authenticate to the CouchDB (Default: empty)")
	restoreAllCmd.Flags().StringP("password", "p", "", "The password to authenticate to the CouchDB (Default: empty)")
	restoreAllCmd.Flags().BoolP("createdb", "c", false, "Create the database if it does not exist on the remote couchdb (Default: false)")
}
