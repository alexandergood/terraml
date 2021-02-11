package terramlparser

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type TerramlParser struct {
	TerramlFileContent  TerramlFileContent
	Variables			interface{}
	DeploymentManifest  map[string]interface{}
	RenderedFilePath    string
}

func DetermineTaskType(task Task) string {
	if !task.Module.IsEmpty() && !task.Resource.IsEmpty() {
		return "multiple_task_type"
	} else if !task.Module.IsEmpty() {
		return "module"
	} else if !task.Resource.IsEmpty() {
		return "resource"
	}

	return "unsupported_task_type"
}

func GetModuleBlock(moduleProvisionTask Module) (string, map[string]interface{}) {
	moduleBlock := make(map[string]interface{})

	source := moduleProvisionTask.Src
	name := moduleProvisionTask.Name
	variables := moduleProvisionTask.Var

	moduleBlock["source"] = source
	for key, value := range variables {
		moduleBlock[key] = value
	}

	return name, moduleBlock
}

func GetResourceBlock(resourceProvisionTask Resource)  (string, map[string]interface{}) {
	resourceBlock := make(map[string]interface{})

	resourceType := resourceProvisionTask.ResourceType
	name := resourceProvisionTask.Name
	config := resourceProvisionTask.Config

	resourceBlock[name] = make(map[string]interface{})
	resourceBlockConfig := make(map[string]interface{})
	for key, value := range config {
		resourceBlockConfig[key] = value
	}

	resourceBlock[name] = resourceBlockConfig
	return resourceType, resourceBlock
}

func GetTaskBlock(task Task) (string, string, map[string]interface{}, error) {
	taskType := DetermineTaskType(task)
	if taskType == "module" {
		taskRef, moduleBlock := GetModuleBlock(task.Module)
		return taskType, taskRef, moduleBlock, nil
	} else if taskType == "resource" {
		taskRef, resourceBlock := GetResourceBlock(task.Resource)
		return taskType, taskRef, resourceBlock, nil
	} else {
		return "", "", nil, fmt.Errorf("invalid task type: %v", taskType)
	}
}

func (p *TerramlParser) RenderTerramlFileWithVariables(filePath string)  error {
	tmpl, err := template.New(filePath).Funcs(sprig.FuncMap()).ParseFiles(filePath)

	if err != nil {
		return errors.WithStack(err)
	}

	renderedFilePath := filepath.Base(filePath) + "-" + uuid.New().String()
	f, err := os.Create(renderedFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tmpl.Execute(f, p.Variables)
	if err != nil {
		return errors.WithStack(err)
	}

	p.RenderedFilePath = renderedFilePath
	return nil
}

func (p *TerramlParser) Cleanup() {
	err := os.Remove(p.RenderedFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDeploymentManifest(filePath string, variableFilePath string) ([]string, map[string]interface{}, error) {
	p := &TerramlParser{}
	p.DeploymentManifest = make(map[string]interface{})

	if err := p.LoadVariables(variableFilePath); err != nil {
		return nil, nil, err
	}

	err := p.RenderTerramlFileWithVariables(filePath)
	if err != nil {
		return nil, nil, err
	}

	// defer p.Cleanup()
	if err := p.ParseTerramlFile(); err != nil {
		return nil, nil, err
	}

	p.GetTerraformConfBlock()
	p.GetProviderBlock()
	deploymentOrder, err := p.GetProvisionBlock()
	if err != nil {
		return nil, nil, fmt.Errorf("error processing provision blocks: %v", err)
	}

	return deploymentOrder, p.DeploymentManifest, nil
}

func (p *TerramlParser) ValidateInput() error {
	if err := p.ValidateTerraformConfiguration(); err != nil {
		return errors.WithStack(err)
	}
	if p.TerramlFileContent.FileType == "playbook" {
		return p.ValidatePlaybookInput()
	} else if p.TerramlFileContent.FileType == "template" {
		return p.ValidateTemplate()
	} else {
		return errors.WithStack(fmt.Errorf("unrecognized file type"))
	}
}

func (p *TerramlParser) ValidateTerraformConfiguration() error {
	return nil
}

func (p *TerramlParser) ValidatePlaybookInput() error {
	return nil
}

func (p *TerramlParser) ValidateTemplate() error {
	return nil
}

func (p *TerramlParser) ParseTerramlFile() error {
	terramalFile, err := ioutil.ReadFile(p.RenderedFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = yaml.Unmarshal(terramalFile, &p.TerramlFileContent)
	if err != nil {
		return errors.WithStack(err)
	}

	err = p.ValidateInput()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (p *TerramlParser) LoadVariables(varFilePath string) error {
	variableFile, err := ioutil.ReadFile(varFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = CustomizedJSONUnmarshal(variableFile, &p.Variables)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (p *TerramlParser) GetTerraformConfBlock() {
	backendType := p.TerramlFileContent.TerraformConf.RemoteState.BackendType
	config := p.TerramlFileContent.TerraformConf.RemoteState.Config

	terraformConfBlock := make(map[string]interface{})
	terraformConfBackendBlock := make(map[string]interface{})
	terraformConfBackendBlock[backendType] = config

	terraformConfBlock["backend"] = terraformConfBackendBlock

	p.DeploymentManifest["terraform"] = terraformConfBlock
}

func (p *TerramlParser) GetProviderBlock() {
	providerBlock := make(map[string]interface{})
	for _, provider := range p.TerramlFileContent.Providers {
		providerBlock[provider.Name] = provider.Config
	}

	p.DeploymentManifest["provider"] = providerBlock
}

func (p *TerramlParser) GetProvisionBlock() ([]string, error) {
	provisionModuleBlock := make(map[string]interface{})
	provisionResourceBlock := make(map[string]interface{})
	var deploymentOrder []string

	for _, task := range p.TerramlFileContent.Provision {
		taskType, taskRef, taskBlock, err := GetTaskBlock(task)
		if err != nil {
			return nil, err
		} else {
			deploymentOrder = append(deploymentOrder, taskType + "/" + taskRef)
			if taskType == "module" {
				provisionModuleBlock[taskRef] = taskBlock
			} else if taskType == "resource" {
				provisionResourceBlock[taskRef] = taskBlock
			}
		}
	}

	p.DeploymentManifest["module"] = provisionModuleBlock
	p.DeploymentManifest["resource"] = provisionResourceBlock
	return deploymentOrder, nil
}