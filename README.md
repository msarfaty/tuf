# tuf

The "Terraform Un-F*ck Tool" or "Terraform Unweaver Framework" is a tool designed to simplify the migration of resources and modules between Terraform workspaces. It aims to make these migrations as seamless and configuration-free as possible.

`tuf` can:
- Move actual HCL (HashiCorp Configuration Language) blocks.
- Remediate resource and module outputs that are used in moved blocks by hardcoding their values from the original workspace in those modules.
- Handle Terraform state remediation.

# Key Opinions

To minimize configuration requirements, `tuf` operates with a set of predefined opinions:

### Comment Relevancy and HCL Block "Prettiness"
When moving resources or modules, `tuf` retains comments directly associated with the moved blocks.

### Data Block Portability
Data blocks are automatically copied when deemed necessary for the migration.

### Provider Same-ness
It is (temporarily) assumed that provider blocks with the same aliases are fully interoperable.

# Disclaimer

`tuf` is currently in pre-release. Use it with caution, as opinions and functionality may change without notice.

# Examples

### Initialize a Terraform Migration
```
tuf init --workspace /path/to/workspace/a \
  --workspace /path/to/workspace/b \
  --terraform-state-pull-command='AWS_PROFILE="PRODUCTION"; terraform state pull > terraform.tfstate' \
  --terraform-state-file-name=terraform.tfstate
```

### Move Resources and Modules Between Workspaces
```
tuf mv /path/to/workspace/a:module.example /path/to/workspace/b:module.example
tuf mv /path/to/workspace/b:aws_security_group.foo /path/to/workspace/b:aws_security_group.bar
```

### Finalize the Migration
```
tuf finalize
```

# Installation

Install `tuf` using the following command:

```
go install github.com/msarfaty/tuf@latest
```

Ensure your Go binary directory is included in your system's `PATH`.

# Pre-release Checklist

- [x] `tuf init`
- [ ] `tuf move`
  - [x] File manipulation
  - [ ] Move tracking
- [ ] `tuf finalize`
  - [ ] State remediation
  - [ ] Variable reference updates

# License

`tuf` is licensed under the MIT License.