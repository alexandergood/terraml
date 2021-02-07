package main

import (
	TerramlParser "github.com/zakufish/terraml/terramlparser"
	"fmt"
	"flag"
)

func main() {
	var filePath string
	var variableFilePath string
	flag.StringVar(&filePath, "file", "", "Path to file")
	flag.StringVar(&variableFilePath, "variables", "", "Path to variables file")

	flag.Parse()

	if filePath == "" {
		panic(fmt.Errorf("no file supplied to execute"))
	}

	deploymentOrder, deploymentManifest, err := TerramlParser.GetDeploymentManifest(filePath, variableFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Deployment Order: ", deploymentOrder)
	fmt.Println("Generated Deployment Manifest: ", string(deploymentManifest))
}