/*
Copyright © 2025 Nicolò Piovan <nicopiovan@gmail.com>
*/

package commons

import (
	"bufio"
	"dbackupcli/cmd/struct/couchdb"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	urlProtocol       = "http://"
	headerContentType = "Content-Type"
	valueJSON         = "application/json"
)

const (
	ErrCreateHTTPRequest  = "error while creating the http request: %v"
	ErrPerformHTTPRequest = "error performing the http request: %v"
	ErrReadResponseBody   = "error reading response body: %v"
	ErrUnmarshalJSON      = "error unmarshalling JSON: %v"
	ErrRemoveFile         = "error removing file %s: %v"
)

func CheckFlags(args []string) bool {
	isMissing := false
	for _, arg := range args {
		if arg == "" {
			isMissing = true
			break
		}
	}
	return isMissing
}

func OverWriteFile(fileName string) error {
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("File %s already exists. Do you want to overwrite it? (y/n): ", fileName)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "y" && input != "Y" {
			return errors.New("operation canceled. The file will not be overwritten")
		}

		fmt.Println("Overwriting file...")
		err := os.Remove(fileName)
		if err != nil {
			return fmt.Errorf(ErrRemoveFile, fileName, err)
		} else {
			fmt.Printf("File %s removed.\n", fileName)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking file %s: %v", fileName, err)
	}
	return nil
}

func GetDBs(host string, port int, user string, password string) ([]string, error) {
	url := urlProtocol + host + ":" + strconv.Itoa(port) + "/_all_dbs"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateHTTPRequest, err)
	}

	req.Header.Add(headerContentType, valueJSON)
	req.SetBasicAuth(user, password)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(ErrPerformHTTPRequest, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(ErrReadResponseBody, err)
	}
	var dbNames []string

	err = json.Unmarshal(body, &dbNames)
	if err != nil {
		return nil, fmt.Errorf(ErrUnmarshalJSON, err)
	}
	return dbNames, nil
}

func GetDB(host string, port int, user string, password string, dbName string) (int, couchdb.Database, error) {
	url := urlProtocol + host + ":" + strconv.Itoa(port) + "/" + dbName
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, couchdb.Database{}, fmt.Errorf(ErrCreateHTTPRequest, err)
	}
	req.Header.Add(headerContentType, valueJSON)
	req.SetBasicAuth(user, password)
	res, err := client.Do(req)
	if err != nil {
		return 0, couchdb.Database{}, fmt.Errorf(ErrPerformHTTPRequest, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, couchdb.Database{}, fmt.Errorf(ErrReadResponseBody, err)
	}

	var db couchdb.Database
	if res.StatusCode != 200 {
		return res.StatusCode, couchdb.Database{}, fmt.Errorf("%s", string(body))
	} else {
		err = json.Unmarshal(body, &db)
		if err != nil {
			return res.StatusCode, couchdb.Database{}, fmt.Errorf(ErrUnmarshalJSON, err)
		}
	}
	return res.StatusCode, db, nil
}

func DeleteDatabase(host string, port int, user string, password string, dbName string) error {
	url := urlProtocol + host + ":" + strconv.Itoa(port) + "/" + dbName
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf(ErrCreateHTTPRequest, err)
	}
	req.Header.Add(headerContentType, valueJSON)
	req.SetBasicAuth(user, password)
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(ErrPerformHTTPRequest, err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error deleting the database: %s", res.Status)
	}
	fmt.Printf("Database %s deleted.\n", dbName)
	return nil
}

func SelectDatabase(options []string, selectedDb *string) error {
	prompt := &survey.Select{
		Message: "Select a database to backup:",
		Options: options,
		Default: options[0],
	}
	err := survey.AskOne(prompt, selectedDb)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

func GetAuthFlagValues(cmd *cobra.Command) (user string, host string, password string, port int) {
	user, _ = cmd.Flags().GetString("user")
	password, _ = cmd.Flags().GetString("password")
	host, _ = cmd.Flags().GetString("host")
	port, _ = cmd.Flags().GetInt("port")
	return
}

func PrepareCmdAuthArgs(args []string, user string, password string, host string, port int) []string {
	args = append(args, "-u", user, "-p", password, "-H", host)

	if port != 5984 {
		args = append(args, "--port", strconv.Itoa(port))
	}
	return args
}
