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

func (c *TerraformCodeGenerator) WriteTerraformFile(terraformFileBlock interface{}, codeDirectory string) error {
	terraformFileContentString, err := json.MarshalIndent(terraformFileBlock, "", strings.Repeat(" ", 4))
	if err != nil {
		return errors.WithStack(err)
	}

	fileName := codeDirectory + "/main.tf.json"
	err = ioutil.WriteFile(fileName, terraformFileContentString, 0777)

	return nil
}

func (c *TerraformCodeGenerator) GenerateTerraformDeploymentDirectory(deploymentType string, deploymentName string) (string, error) {
	codeDirectory := uuid.New().String()
	if err := os.Mkdir(codeDirectory, 0777); err != nil {
		return "", errors.WithStack(err)
	}

	deploymentElementInputs := make(map[string]interface{})
	deploymentElementInputs[deploymentName] = c.DeploymentManifest[deploymentType].(map[string]interface{})[deploymentName]

	terraformFileBlock := make(map[string]interface{})
	terraformFileBlock["terraform"] = c.DeploymentManifest["terraform"]
	terraformFileBlock["provider"] = c.DeploymentManifest["provider"]
	terraformFileBlock[deploymentType] = deploymentElementInputs

	if err := c.WriteTerraformFile(terraformFileBlock, codeDirectory); err != nil {
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