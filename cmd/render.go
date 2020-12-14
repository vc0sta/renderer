/*
Copyright © 2020 NAME HERE vinicius.costa.92@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	// "net/http"
	// "os"
	// "time"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"log"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func init() {

	// Cmd.AddCommand(check.PostgresCmd, check.HttpCmd)
	// Cmd.AddCommand(check.HttpCmd)
	rootCmd.AddCommand(RenderCmd)
}

var instructions = `Check if a http/https is reachable.
It will not test if the returning code is 2XX.`

// RenderCmd represents the render command
var RenderCmd = &cobra.Command{
	Use:     "render",
	Short:   "Renders a given template with provided values",
	Long:    `Renders a given template with provided values`,
	Example: "renderer render -values values.json -template template",
	Run: func(cmd *cobra.Command, args []string) {
		CallRendering(args)
	},
}

var (
	values    string
	templates string
	recursive bool
)

func init() {
	RenderCmd.Flags().StringVar(&values, "values", "values.json", "file containing the message with the values")
	RenderCmd.Flags().StringVar(&templates, "templates", "template", "folder containing the templates to be renderized")
	RenderCmd.Flags().BoolVar(&recursive, "recursive", true, "should we use templates from all subdirs?")
}

func CallRendering(args []string) {

	jsonData, _ := ReadFile(values)

	var v interface{}
	json.Unmarshal(jsonData, &v)

	data := v.(map[string]interface{})

	result, err := RenderFolder(data, templates, recursive)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Prints renderized template
	for _, v := range result {
		
		fmt.Println(string(v))
	}

}

// CreaRenderizedFile 
func CreateRenderizedFile(path string, result []byte) error {
	str := strings.Split(path, "/")
	str = str[:len(str) - 1]
	original_path := strings.Join(str, "/")

	err := os.MkdirAll(original_path, 0755)
	if err != nil {
		log.Fatal(err)
}
	os.Create(path)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Fprintln(f, string(result))
	return nil
}	


// ReadFile receives a canonical path to a file and returns its containts as a []byte
func ReadFile(file string) ([]byte, error) {
	dat, err := ioutil.ReadFile(file)
	return dat, err
}

// RenderTemplate receives any struct as a value and renderize a given template using its information.
func RenderTemplate(values interface{}, file string) ([]byte, error) {
	// The template name must match a file in ParseFiles which is the main template
	// See: https://stackoverflow.com/questions/10199219/go-template-function
	tmpl, err := template.New(filepath.Base(file)).ParseFiles(file)
	if err != nil {
		// log.Error(err)
		return []byte{}, err
	}

	var tpl bytes.Buffer
	// Here we use the template and conf to make to generate textual output
	// We are using 'os.Stdout` to output to screen, a file can be used instead
	err = tmpl.Execute(&tpl, values)

	if err != nil {
		// log.Error(err)
		return []byte{}, err
	}
	err = CreateRenderizedFile("out/" + file, tpl.Bytes())
	if err != nil {
		return []byte{}, err
	}

	return tpl.Bytes(), nil
}

// RenderFolder get all yaml files from a given folder and create a list of renderized templates.
func RenderFolder(values interface{}, path string, subdir bool) ([][]byte, error) {
	var err error
	var files []string

	if subdir {
		files, err = getFilesRecursive(path)
	} else {
		files, err = getFiles(path)
	}

	if err != nil {
		// log.Error(err)
		return [][]byte{}, err
	}

	var yamlFiles [][]byte

	for _, file := range files {
		yaml, err := RenderTemplate(values, file)
		if err != nil {
			// log.Error(err)
			return yamlFiles, err
		}

		yamlFiles = append(yamlFiles, yaml)

	}

	return yamlFiles, nil
}

// getFiles is used to properly read a given directory and return a list of YAML file names
func getFiles(dir string) ([]string, error) {
	var fileList []string

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return []string{}, err
	}

	for _, file := range files {

		if filepath.Ext(filepath.Base(file.Name())) == ".yaml" || filepath.Ext(filepath.Base(file.Name())) == ".yml" {
			fileList = append(fileList, dir+"/"+file.Name())
		}
	}

	return fileList, nil
}

// getFilesRecursive is used to properly read a given directory and its subdirectories and return a list of YAML file names
func getFilesRecursive(dir string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(filepath.Base(path)) == ".yaml" || filepath.Ext(filepath.Base(path)) == ".yml" {
				fileList = append(fileList, path)
			}
			return nil
		})
	if err != nil {
		// log.Println(err)
		return fileList, err
	}

	return fileList, err
}
