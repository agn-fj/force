package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var cmdExport = &Command{
	Run:   runExport,
	Usage: "export [dir]",
	Short: "Export metadata to a local directory",
	Long: `
Export metadata to a local directory

Examples:

  force export

  force export org/schema
`,
}

func runExport(cmd *Command, args []string) {
	wd, _ := os.Getwd()
	root := filepath.Join(wd, "metadata")
	if len(args) == 1 {
		root, _ = filepath.Abs(args[0])
	}
	force, _ := ActiveForce()
	sobjects, err := force.ListSobjects()
	if err != nil {
		ErrorAndExit(err.Error())
	}
	stdObjects := make([]string, 1, len(sobjects)+1)
	stdObjects[0] = "*"
	for _, sobject := range sobjects {
		name := sobject["name"].(string)
		if !sobject["custom"].(bool) && !strings.HasSuffix(name, "__Tag") && !strings.HasSuffix(name, "__History") && !strings.HasSuffix(name, "__Share") {
			stdObjects = append(stdObjects, name)
		}
	}
	//jwf-hack; shortened list temporarily
	query := ForceMetadataQuery{
		{Name: "Report", Members: []string{"RUSH_Agent_Reports"}, AllowsWildcard: false, FolderTypeName: "ReportFolder"}, //brings back folder
		{Name: "ApexTrigger", Members: []string{"*"}, AllowsWildcard: true, FolderTypeName: ""},
		{Name: "Workflow", Members: []string{"*"}, AllowsWildcard: true, FolderTypeName: ""},

		//foldertypenames
		//"DashboardFolder"
		//"DocumentFolder"
		//"EmailFolder"
		//"ReportFolder"
	}

	files, err := force.Metadata.Retrieve(query)
	if err != nil {
		fmt.Printf("Encountered and error with retrieve...\n")
		ErrorAndExit(err.Error())
	}
	for name, data := range files {
		file := filepath.Join(root, name)
		dir := filepath.Dir(file)
		if err := os.MkdirAll(dir, 0755); err != nil {
			ErrorAndExit(err.Error())
		}
		if err := ioutil.WriteFile(filepath.Join(root, name), data, 0644); err != nil {
			ErrorAndExit(err.Error())
		}
	}
	fmt.Printf("Exported to %s\n", root)
}
