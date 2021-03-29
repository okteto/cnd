// Copyright 2020 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/okteto/okteto/pkg/log"
)

// FileExists return true if the file exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		log.Infof("failed to check if %s exists: %s", name, err)
	}

	return true
}

// CopyFile copies a binary between from and to
func CopyFile(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}

	// skipcq GSC-G302 syncthing is a binary so it needs exec permissions
	toFile, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		return err
	}

	defer toFile.Close()

	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return err
	}

	return nil
}

//GetValidNameFromFolder returns a valid kubernetes name for a folder
func GetValidNameFromFolder(folder string) (string, error) {
	dir, err := filepath.Abs(folder)
	if err != nil {
		return "", fmt.Errorf("error inferring name: %s", err)
	}
	name := filepath.Base(dir)
	name = strings.ToLower(name)
	name = ValidKubeNameRegex.ReplaceAllString(name, "-")
	log.Infof("autogenerated name: %s", name)
	return name, nil
}

// GetFileByRegex return the first file that matches the regex
func GetFileByRegex(regexString string) string {
	file, err := os.Open(".")
	if err != nil {
		return ""
	}
	defer file.Close()

	fileNames, _ := file.Readdirnames(0)

	regex, err := regexp.Compile(regexString)
	for _, name := range fileNames {
		if regex.Match([]byte(name)) {
			return name
		}
	}
	return ""
}
