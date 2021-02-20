# Terraml
Terraml is a thin wrapper for Terraform that enables users to build complex and repetitive infrastructures through YAML. That is, a tool that liberates engineers from writing repeptitive Terraform/Terragrunt codes over and over again.

# Why Terraml
My job as a Platform Engineer @ [Quantium](https://quantium.com/) involves provisioning platform resources and infrastructures for data and application engineering teams to use my team's engineering platform offering. In order to keep our Terraform codebase "[DRY](https://terragrunt.gruntwork.io/docs/features/keep-your-terraform-code-dry/)", we adopted [Terragrunt](https://terragrunt.gruntwork.io/) as our deployment tool and built infrastructures using the GitOps approach. Everything seemed fine until more and more teams and projects began to jump onto our platform as we began to scale up our infrastrctures - that's when we started to realize many of Terragrunt's shortcomings:
- Most of our projects required indenpendent platform infrstructure/resource, but 95% of those infrastructure configuration were identical (with minor difference). Thus, our codebase got very repetitive.
- Everytime a root module got changed, we had to update every single Terragrunt file correspondingly.

Eventually, our infrastructure codebase got to a point where it was almost impossible to maintain/patch. A trivial change from cloud provider would force us to effectively modify hundreds of Terragrunt files.

Although Terragrunt wasn't perfect, it had its advantages, included handy state file management, variables loader, etc. After giving this problem some thoughts, I decided to develop an open source solution which is capable of handling repetitive infrastructure code, as well as incorporating with Terragrunt's many nice features.

That's how Terraml was made.

# How to Build & Install Terraml

Terraml is still in alpha. I am planning to release a package when it reaches ```v1.0```. For now, you can build and install Terraml on Linux via:
```bash
go build && mv terraml /usr/local/bin/terraml
```

# How to Use Terraml

To achieve our goal of minimizing repetitive infrastructure code, Terraml introduces the concept of ```infrastructure template``` - a template file that will be rendered and executed against user's variable inputs. We've provided a simple [template file](https://github.com/zakufish/terraml/blob/main/examples/sample.tpl.yml) and a [variable file](https://github.com/zakufish/terraml/blob/main/examples/variables.yml) inside ```examples/``` directory, which can be executed by running:

```
terraml --action init --file examples/sample.tpl.yml --variables examples/variables.yml
```

In this case, two code directories will be generated (since there will be two modules, given the iteration of ```clusters``` list defined in our template), each contains a ```main.tf.json``` file for terraform to consume. The ```main.tf.json``` files will also include two different types of ```terraform blocks```:
1. Terraform Configuration - ```backend``` and ```required_providers```. These configurations will become the "global" terraform setting for this particular terraml job.
2. Module - the provision block.

As Terraml runs these generated Terraform files, the code directories will also get cleaned up.

# TODO

Terraml is still in very early stage of development, and I am planning to get it to ```v1.0``` as soon as possible. There're a couple of things I need to take care of before reaching that point:
1. Expand templating and parsing to cover all Terraform code blocks.
2. Add tests.
3. Fix weird edge cases.
