package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	TerramlCodeGenerator "github.com/zakufish/terraml/codegen"
	TerramlCodeRunner "github.com/zakufish/terraml/coderun"
	TerramlParser "github.com/zakufish/terraml/terramlparser"
)

func main() {
	var filePath string
	var variableFilePath string
	var action string
	flag.StringVar(&filePath, "file", "", "Path to file")
	flag.StringVar(&variableFilePath, "variables", "", "Path to variables file")
	flag.StringVar(&action, "action", "", "Action item")

	flag.Parse()

	if filePath == "" {
		fmt.Println(fmt.Errorf("no file supplied to execute"))
	}

	fmt.Println("Rendering and generating the deployment manifests ...")
	deploymentOrder, deploymentManifest, err := TerramlParser.GetDeploymentManifest(filePath, variableFilePath)
	if err != nil {
		fmt.Println(errors.Wrap(err, "terraml: template rendering error"))
	}

	fmt.Println("Generating Terraform code directories ...")
	executeOrder, err := TerramlCodeGenerator.GenerateTerraformCodeDirectories(deploymentOrder, deploymentManifest)
	if err != nil {
		fmt.Println(errors.Wrap(err, "terraml: terraform code generation error"))
	}

	fmt.Println("Executing Terraform commands ...")
	if err := TerramlCodeRunner.RunTerraformCode(executeOrder, action); err != nil {
		fmt.Println(errors.Wrap(err, "terraml: terraform command execution error"))
	}
}