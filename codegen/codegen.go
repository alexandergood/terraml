package codegen

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
)

type TerraformCodeGenerator struct {
	DeploymentOrder    []string
	DeploymentManifest map[string]interface{}
}

func GenerateTerraformCodeDirectories(deploymentOrder []string, deploymentManifest map[string]interface{}) ([]string, error) {
	c := &TerraformCodeGenerator{deploymentOrder, deploymentManifest}
	return c.GenerateTerraformCodeDirectories()
}

func (c *TerraformCodeGenerator) WriteTerraformFile(terraformFileType string, terraformFileContent interface{}, codeDirectory string) error {
	terraformFileContentString, err := json.MarshalIndent(terraformFileContent, "", strings.Repeat(" ", 4))
	if err != nil {
		return errors.WithStack(err)
	}

	fileName := codeDirectory + "/" + terraformFileType + ".tf.json"
	err = ioutil.WriteFile(fileName, terraformFileContentString, 0777)

	return nil
}

func (c *TerraformCodeGenerator) GenerateTerraformDeploymentDirectory(deploymentType string, deploymentName string) (string, error) {
	codeDirectory := uuid.New().String()
	if err := os.Mkdir(codeDirectory, 0777); err != nil {
		return "", errors.WithStack(err)
	}

	terraformConfElement := make(map[string]interface{})
	terraformConfElement["terraform"] = c.DeploymentManifest["terraform"]

	providerElement := make(map[string]interface{})
	providerElement["provider"] = c.DeploymentManifest["provider"]

	deploymentElement := make(map[string]interface{})
	deploymentElementInputs := make(map[string]interface{})
	deploymentElementInputs[deploymentName] = c.DeploymentManifest[deploymentType].(map[string]interface{})[deploymentName]
	deploymentElement[deploymentType] = deploymentElementInputs

	if err := c.WriteTerraformFile("terraform", terraformConfElement, codeDirectory); err != nil {
		return "", err
	}
	if err := c.WriteTerraformFile("provider", providerElement, codeDirectory); err != nil {
		return "", err
	}
	if err := c.WriteTerraformFile(deploymentType, deploymentElement, codeDirectory); err != nil {
		return "", err
	}

	return codeDirectory, nil
}

func (c *TerraformCodeGenerator) GenerateTerraformCodeDirectories() ([]string, error) {
	codeDirectories := make([]string, 0)

	for _, deploymentItem := range c.DeploymentOrder {
		deploymentItems := strings.Split(deploymentItem, "/")
		deploymentType, deploymentName := deploymentItems[0], deploymentItems[1]

		codeDirectory, err := c.GenerateTerraformDeploymentDirectory(deploymentType, deploymentName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		codeDirectories = append(codeDirectories, codeDirectory)
	}

	return codeDirectories, nil
}