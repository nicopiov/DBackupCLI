/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbackupcli",
	Short: "DBackupCLI è una CLI tool per facilitare le operazioni di backup sui database",
	Long: `DBackupCLi è una CLI tool per facilitare le operazioni di backup sui, attualmente è possibile operare
su database NoSQL Couchdb eseguendo:
- Elencare i database presenti nell'istanza di CouchDB
- Backup/Restore di un singolo database
- Backup/Restore di un'intera istanza CouchDB`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetUsageTemplate(`
Usage: {{.Use}} [operation]

Operation:
	backup		Perform a backup on a single database
	restore		Perform the restore of a dumped database
	listdbs		List all the databases in the specified CouchDB
	backupAll	Perform a backup of the entire CouchDB 

Use "{{.Use}} [operation]" -h" for more information about a module.
`)

	rootCmd.Flags().BoolP("help", "h", false, "Help message for dbackupcli")
}
