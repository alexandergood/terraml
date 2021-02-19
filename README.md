# Terraml
Terraml is a thin wrapper for Terraform that enables users to build complex and repetitive infrastructures through YAML. That is, a tool that liberates engineers from writing repeptitive Terraform/Terragrunt codes over and over again.

# Why Terraml
My job as a Platform Engineer @ [Quantium](https://quantium.com/) involves provisioning platform resources and infrastructures for data and application engineering teams to use my team's engineering platform offering. In order to keep our Terraform codebase "[DRY](https://terragrunt.gruntwork.io/docs/features/keep-your-terraform-code-dry/)", we adopted [Terragrunt](https://terragrunt.gruntwork.io/) as our deployment tool and built infrastructures using the GitOps approach. Everything seemed fine until more and more teams and projects began to jump onto our platform as we began to scale up our infrastrctures - that's when we started to realize many of Terragrunt's shortcomings:
- Most of our projects required indenpendent platform infrstructure/resource, but 95% of those infrastructure configuration were identical (with minor difference). Thus, our codebase got very repetitive.
- Everytime a root module got changed, we had to update every single Terragrunt file correspondingly.

Eventually, our infrastructure codebase got to a point where it was almost impossible to maintain/patch. A trivial change from cloud provider would force us to effectively modify hundreds of Terragrunt files.

Although Terragrunt wasn't perfect, it had its advantages, included handy state file management, variables loader, etc. So, we decided to develop a solution which incorporated with Terragrunt's nice features, as well as capable of handling repetitive infrastructure code.
