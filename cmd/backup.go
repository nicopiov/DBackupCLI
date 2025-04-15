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

	"github.com/spf13/cobra"
)

// backupCouchdbCmd represents the backupCouchdb command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Operates a backup of a database on CouchDB",
	Long:  `Operates a backup of the specified database on a certain CouchDB.`,
	Run: func(cmd *cobra.Command, args []string) {
		tmpFile, err := scripts.GetEmbeddedScripts()
		defer os.Remove(tmpFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		user, host, password, port := commons.GetAuthFlagValues(cmd)
		file, _ := cmd.Flags().GetString("file")
		if commons.CheckFlags(append([]string{}, file, user, password, host)) {
			fmt.Println("missing on or more flags/arguments\n\nCheck using 'dbackupcli backup -h'")
			os.Exit(1)
		}

		dbsList, err := commons.GetDBs(host, port, user, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		var selectedDatabase string
		if err = commons.SelectDatabase(append(dbsList, "exit"), &selectedDatabase); err != nil {
			fmt.Println(err)
			return
		}

		if selectedDatabase == "exit" {
			os.Exit(0)
		}

		var cmdArgs []string = []string{"-b", "-d", selectedDatabase}
		cmdArgs = append(cmdArgs, "-f", file)

		err = commons.OverWriteFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		cmdArgs = commons.PrepareCmdAuthArgs(cmdArgs, user, password, host, port)
		cmdExec := exec.Command("bash", append([]string{tmpFile}, cmdArgs...)...)
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdout = os.Stdout
		if err := cmdExec.Run(); err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("Backup completed successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.SetUsageTemplate(`
Usage: dbackupcli {{.Use}} [flags]

Flags:
 -h, --help		Show this help message
 -f, --file		The filename where to dump the backup (e.g dump.json)
 -u, --user		The CouchDB username for the auth
 -p, --password		The CouchDB password for the auth
 --host			The host of the remote CouchDB, with or without 'https://'
 --port			The port of the remote CouchDB, default is 5984

After entering the command a prompt will let you select the database to backup

Examples:
 dbackupcli backup -f dump.json -u admin -p root --host 127.0.0.1 --port 9876
 dbackupcli backup --file dump.json --user admin -p root --host 127.0.0.1.
`)
	backupCmd.Flags().BoolP("help", "h", false, "Help message")
	backupCmd.Flags().String("host", "", "The remote CouchDB host, can be provided with or without 'https://'")
	backupCmd.Flags().Int("port", 5984, "The remote CouchDB port (Default: 5984)")
	backupCmd.Flags().StringP("file", "f", "", "The name of the file where to backup (Default: empty)")
	backupCmd.Flags().StringP("user", "u", "", "The username to authenticate to the CouchDB (Default: empty)")
	backupCmd.Flags().StringP("password", "p", "", "The password to authenticate to the CouchDB (Default: empty)")
}
