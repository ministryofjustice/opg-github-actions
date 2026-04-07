# Snyk IaC Composite Action

Run `snyk iac test` against a directory of Infrastructure as Code files.

The action writes SARIF output so results can be uploaded to the GitHub Security tab.

## OPG ECR Pull Requirement

The AWS authentication steps in the examples are required because this action runs Snyk via our OPG-managed container image, which is built here:

https://github.com/ministryofjustice/opg-snyk-image-builder

We may add an option in future to use the public Snyk Docker image.

## Usage

Basic usage for scanning a Terraform directory:

```yaml
- name: Configure AWS Credentials
  uses: aws-actions/configure-aws-credentials
  with:
    aws-region: eu-west-1
    role-to-assume: ${{ vars.OIDC_MANAGEMENT_ECR_ROLE }}
    role-duration-seconds: 3600
    role-session-name: GitHubActions

- name: ECR Login
  id: login-ecr
  uses: aws-actions/amazon-ecr-login
  with:
    registries: <managment-id>

- name: Snyk IaC
  id: snyk_iac
  uses: ministryofjustice/opg-github-actions/actions/snyk-iac
  with:
    iac_test_directory: "./terraform/environment/"
    SNYK_TOKEN: ${{ secrets.snyk_token }}
```

### Example With SARIF Upload

```yaml
- name: Snyk IaC
  id: snyk_iac
  uses: ministryofjustice/opg-github-actions/actions/snyk-iac
  with:
    iac_test_directory: "./terraform/environment/"
    SNYK_TOKEN: ${{ secrets.snyk_token }}

- name: Upload Snyk IaC results to GitHub Security tab
  uses: github/codeql-action/upload-sarif
  if: always()
  with:
    category: "IaC"
    sarif_file: "test-results/snyk.sarif"
```

### Example With Additional Snyk Arguments

```yaml
- name: Snyk IaC (fail on high severity only)
  uses: ministryofjustice/opg-github-actions/actions/snyk-iac
  with:
    iac_test_directory: "./terraform/environment/"
    SNYK_TOKEN: ${{ secrets.snyk_token }}
    snyk_test_args: "--severity-threshold=high"
```

## Inputs

### Required Inputs

- `SNYK_TOKEN`
  - Snyk API token used for authentication.

### Optional Inputs

- `iac_test_directory` (default: `./`)
  - Directory containing IaC files to test.

- `sarif_filepath` (default: `test-results`)
  - Directory where the SARIF output file is written.

- `sarif_filename` (default: `snyk.sarif`)
  - SARIF output file name.

- `snyk_policy_path` (default: `./.snyk`)
  - Path to a `.snyk` policy file to mount into the container.

- `snyk_test_args` (default: `""`)
  - Extra arguments passed directly to `snyk iac test`.

## Notes

- The action mounts the IaC directory read-only at `/app/iac`.
- The action mounts the policy file into `/root/.snyk` and sets `--policy-path=/root/.snyk`.
- Ensure your workflow has access to Docker and can run `docker compose`.
