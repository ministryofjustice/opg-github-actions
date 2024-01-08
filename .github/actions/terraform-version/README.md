# Terraform Version Composite Action

Used to find and then switch terraform versions in projects where multiple versions maybe in use (monorepos with many state files etc).

Parse the terraform versions file from the directory passed and return the required terraform version information.

**Requires `required_version` line to be present.**

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - id: terraform_version
      name: "Get terraform version"
      uses: 'ministryofjustice/opg-github-actions/.github/actions/terraform-version@v2.1.3'
      with:
        terraform_directory: "./terraform/"
```

## Inputs and Outputs

Inputs:
- `terraform_directory`
- `terraform_versions_file` (default: ./versions.tf)
- `simple_file`

Outputs:
- `terraform_directory`
- `terraform_versions_file`
- `simple_file`
- **`version`**



### Inputs

#### `terraform_directory`
Path to the root of the terraform state you wish to find the version for.

#### `terraform_versions_file` (default: ./versions.tf)
The file to inspect for the `required_version` information for terraform

#### `simple_file`
Override to allow a plain text file which contains just the version string and nothing else to be used.


### Outputs

#### `terraform_directory`, `terraform_versions_file` and `simple_file`
Mirror of the inputted param.

#### `version`
The version string found within the file.
