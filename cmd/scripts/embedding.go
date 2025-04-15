/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/

package scripts

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed *.sh

var embeddedScripts embed.FS

func GetEmbeddedScripts() (string, error) {
	var tmpFile string
	var scriptFile string

	tmpFile = filepath.Join(os.TempDir(), "couch-script.sh")
	scriptFile = "couch-dump-restore.sh"

	scriptContent, err := embeddedScripts.ReadFile(scriptFile)
	if err != nil {
		return "", fmt.Errorf("error reading embedded script: %s", err)
	}
	err = os.WriteFile(tmpFile, scriptContent, 0755)
	if err != nil {
		return "", fmt.Errorf("error writing embedded script into temp file: %s", err)
	}
	return tmpFile, nil
}
