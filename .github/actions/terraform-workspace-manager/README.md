# Terraform Workspace Manager Composite Action

Handles registering and listing ephemeral environments, typically within development accounts, backed by a dynamodb.

Allows custom time to live setup.

**Note: This does not handle the destruction of any aged environments**

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - id: terraform_workspace
      name: "Register workspace"
      uses: 'ministryofjustice/opg-github-actions/.github/actions/terraform-workspace-manager@v2.1.3'
      with:
          aws_access_key_id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws_secret_access_key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws_account_id: '1111111110'
          aws_iam_role: 'gh-reusable-actions-ci'
          register_workspace: "testworkspace"
          time_to_protect: 1
```

## Inputs and Outputs

Inputs:
- `aws_access_key_id`
- `aws_secret_access_key`
- `aws_account_id`
- `aws_iam_role`
- `register_workspace`
- `time_to_protect` (default: 24)

Outputs:
- `protected_workspaces`
- `workspace_name`


### Inputs

#### `aws_access_key_id` and `aws_secret_access_key`
AWS secrets for auth to access the `aws_account_id` where the dynamodb data source is.

#### `aws_account_id` and `aws_iam_role`
The account and role to use to access the dynamodb data store for adding and retriving workspace data.

#### `register_workspace`
The name of the terraform workspace to register.

#### `time_to_protect` (default: 24)
How long the `register_workspace` should be marked as protected for.


### Outputs

#### `protected_workspaces`
List of existing workspaces within the dynamodb

#### `workspace_name`
The workspace that has been registered
