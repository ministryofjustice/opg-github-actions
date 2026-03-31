# Github Deploy Key Composite Action

Used to start the SSH Agent, and then load a provided Github Deployment SSH Key.

**Requires `required_version` line to be present.**

## Usage

Within you github workflow job you can place a step such as:

```yaml
    - id: deploy-key
      uses: 'ministryofjustice/opg-github-actions/actions/github-deploy-key@821b6f92327f0f195276860676aa8133d63f39dd # v4.5.1'
      with:
        deploy_key: ${{ secrets.GITHUB_SSH_PRIVATE_KEY }}
```

## Inputs and Outputs

Inputs:
- `deploy_key`

### Inputs

#### `deploy_key`
Secret containing the Github Deployment Key you want to load.
