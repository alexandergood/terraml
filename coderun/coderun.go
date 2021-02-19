package coderun

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strings"
)

func DestroyCodeDirectories(codeDirectories []string) {
	for _, codeDirectory := range codeDirectories {
		_ = os.RemoveAll(codeDirectory)
	}
}

func RunTerraformCode(executeOrder []string, action string) error {
	defer DestroyCodeDirectories(executeOrder)

	execPath := os.Getenv("TERRAFORM_EXEC_PATH")
	if execPath == "" {
		return fmt.Errorf("undefined terraform exec path (TERRAFORM_EXEC_PATH unset)")
	}

	if action == "destroy" {
		sort.Sort(sort.Reverse(sort.StringSlice(executeOrder)))
	}

	for _, codeDirectory := range executeOrder {
		deploymentInfo := strings.Split(codeDirectory, "-")
		_ = os.Setenv("TF_VAR_TERRAML_RESOURCE_PATH", strings.Join(deploymentInfo, "/"))

		terraform, err := tfexec.NewTerraform(codeDirectory, execPath)
		if err != nil {
			return errors.WithStack(err)
		}

		if action == "init" {
			err = terraform.Init(context.Background())
		} else if action == "plan" {
			_, err = terraform.Plan(context.Background())
		} else if action == "apply" {
			err = terraform.Apply(context.Background())
		} else if action == "destroy" {
			err = terraform.Destroy(context.Background())
		} else {
			return errors.WithStack(fmt.Errorf("unrecognized action item"))
		}

		if err != nil  {
			return errors.WithStack(err)
		}
	}

	return nil
}