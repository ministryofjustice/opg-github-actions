# Terraform Version Composite Action

Used to find and then switch terraform versions in projects where multiple versions maybe in use (monorepos with many state files etc).

Parse the terraform versions file from the directory passed and return the required terraform version information.

**Requires `required_version` line to be present.**

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - id: terraform_version
      name: "Get terraform version"
      uses: 'ministryofjustice/opg-github-actions/actions/terraform-version@v1.2.3'
      with:
        terraform_directory: "./terraform/"
```

## Inputs and Outputs

Inputs:
- `terraform_directory`
- `terraform_versions_file` (default: `./versions.tf`)

Outputs:
- **`version`**

### Inputs

#### `terraform_directory`
Path to the root of the terraform state you wish to find the version for.

#### `terraform_versions_file` (default: ./versions.tf)
The file to inspect for the `required_version` information for terraform


### Outputs

#### `version`
The version string found within the file.
