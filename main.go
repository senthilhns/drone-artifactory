// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/drone/drone-artifactory/plugin"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func CheckAndCertFileReadFile(filePath string) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("file does not exist")
		return "", fmt.Errorf("file does not exist: %s", filePath)
	} else if err != nil {
		fmt.Println("error checking file")
		return "", fmt.Errorf("error checking file: %v", err)
	}

	// Read the file contents
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("error reading file")
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return string(content), nil
}

func main() {

	fmt.Println("============= no cert file =============")
	filePath := `C:\users\ContainerAdministrator\.jfrog\security\certs\cert.windows.pem`
	_, e := CheckAndCertFileReadFile(filePath)
	if e != nil {
		fmt.Println(e)
	}

	logrus.SetFormatter(new(formatter))

	var args plugin.Args
	if err := envconfig.Process("", &args); err != nil {
		logrus.Fatalln(err)
	}

	switch args.Level {
	case "debug":
		logrus.SetFormatter(textFormatter)
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetFormatter(textFormatter)
		logrus.SetLevel(logrus.TraceLevel)
	}

	if err := plugin.Exec(context.Background(), args); err != nil {
		logrus.Fatalln(err)
	}
}

// default formatter that writes logs without including timestamp or level information.
type formatter struct{}

func (*formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

// text formatter that writes logs with level information
var textFormatter = &logrus.TextFormatter{
	DisableTimestamp: true,
}
