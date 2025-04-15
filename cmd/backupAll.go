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

// backupAllCmd represents the backupAll command
var backupAllCmd = &cobra.Command{
	Use:   "backupAll",
	Short: "Operates a backup of the entire CouchDB istance",
	Long:  `Operates a backup of the entire CouchDB istance that has been specified with the flags`,
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
			fmt.Println("missing on or more flags/arguments\n\nCheck using 'dbackupcli backup -h'")
			os.Exit(1)
		}

		dbsList, err := commons.GetDBs(host, port, user, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		var cmdArgs []string = []string{"-b"}
		cmdArgs = commons.PrepareCmdAuthArgs(cmdArgs, user, password, host, port)
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Println("error:", err)
			return
		}

		for _, db := range dbsList {
			if !strings.HasPrefix(db, "_") && db != "" {
				dbArgs := append(cmdArgs, "-d", db, "-f", dir+"/"+db+".json")
				cmdExec := exec.Command("bash", append([]string{tmpFile}, dbArgs...)...)
				cmdExec.Stderr = os.Stderr
				cmdExec.Stdout = os.Stdout
				if err := cmdExec.Run(); err != nil {
					fmt.Println("Error: ", err)
				} else {
					fmt.Println("Backup completed successfully!")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(backupAllCmd)
	backupAllCmd.SetUsageTemplate(`
Usage: dbackupcli {{.Use}} [flags]

Flags:
 -h, --help		Show this help message
 -f, --filedir	The directory where to dump the entire couchdb istance,
				if not present the backup will create the directory
 -u, --user		The CouchDB username for the auth
 -p, --password		The CouchDB password for the auth
 --host			The host of the remote CouchDB, with or without 'https://'
 --port			The port of the remote CouchDB, default is 5984

After entering the command a prompt will let you select the database to backup

Examples:
 dbackupcli backupAll -d backup-core -u admin -p root --host 127.0.0.1 --port 9876
 dbackupcli backupAll --dir backup-intraner --user admin -p root --host 127.0.0.1.
`)
	backupAllCmd.Flags().BoolP("help", "h", false, "Help message")
	backupAllCmd.Flags().String("host", "", "The remote CouchDB host, can be provided with or without 'https://'")
	backupAllCmd.Flags().Int("port", 5984, "The remote CouchDB port (Default: 5984)")
	backupAllCmd.Flags().StringP("filedir", "f", "", "The name of the directory where to backup (Default: empty)")
	backupAllCmd.Flags().StringP("user", "u", "", "The username to authenticate to the CouchDB (Default: empty)")
	backupAllCmd.Flags().StringP("password", "p", "", "The password to authenticate to the CouchDB (Default: empty)")
}
