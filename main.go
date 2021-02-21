package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	TerramlCodeGenerator "github.com/zakufish/terraml/codegen"
	TerramlCodeRunner "github.com/zakufish/terraml/coderun"
	TerramlParser "github.com/zakufish/terraml/terramlparser"
)

func runAction(c *cli.Context) error {
	action := c.Command.Name
	template := c.String("template")
	variables := c.String("variables")

	if template == "" {
		return errors.New("no template file supplied to execute")
	}

	if variables == "" {
		return errors.New("no variables file supplied to execute")
	}

	fmt.Println("Rendering and generating the deployment manifests ...")
	deploymentOrder, deploymentManifest, err := TerramlParser.GetDeploymentManifest(template, variables)
	if err != nil {
		return errors.Wrap(err, "terraml: template rendering error")
	}

	fmt.Println("Generating Terraform code directories ...")
	executeOrder, err := TerramlCodeGenerator.GenerateTerraformCodeDirectories(deploymentOrder, deploymentManifest)
	if err != nil {
		return errors.Wrap(err, "terraml: terraform code generation error")
	}

	fmt.Println("Executing Terraform commands ...")
	if err := TerramlCodeRunner.RunTerraformCode(executeOrder, action); err != nil {
		return errors.Wrap(err, "terraml: terraform command execution error")
	}
	return nil
}

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "template",
			Aliases:  []string{"t"},
			Usage:    "Load template from `FILE`",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "variables",
			Aliases:  []string{"v"},
			Usage:    "Load variables from `FILE`",
			Required: true,
		},
	}

	terraml := &cli.App{
		Name:  "terraml",
		Usage: "Build complex Terraform through YAML",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "initialize the working directory for Terraform.",
				Action:  runAction,
				Flags:   flags,
			},
			{
				Name:    "plan",
				Aliases: []string{"p"},
				Usage:   "create an execution plan.",
				Action:  runAction,
				Flags:   flags,
			},
			{
				Name:    "apply",
				Aliases: []string{"a"},
				Usage:   "apply the changes required to reach the desired state of the configuration.",
				Action:  runAction,
				Flags:   flags,
			},
			{
				Name:    "destroy",
				Aliases: []string{"d"},
				Usage:   "destroy the Terraform-managed infrastructure.",
				Action:  runAction,
				Flags:   flags,
			},
		},
		EnableBashCompletion: true,
	}

	err := terraml.Run(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
