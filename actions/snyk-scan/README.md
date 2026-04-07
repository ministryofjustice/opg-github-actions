# Snyk Scan Composite Action

Run `snyk container test` against a Docker image. This action is intended for workflows that build an image and want to scan it before pushing to a registry.

The action writes SARIF output so results can be uploaded to the GitHub Security tab.

## OPG ECR Pull Requirement

The AWS authentication steps in the examples are required because this action runs Snyk via our OPG-managed container image, which is built here:

https://github.com/ministryofjustice/opg-snyk-image-builder

We may add an option in future to use the public Snyk Docker image.

## Usage

Basic usage for scanning a docker image:

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

- name: Snyk Scan
  id: snyk_scan
  uses: ministryofjustice/opg-github-actions/actions/snyk-scan
  with:
    SNYK_TEST_IMAGE: my-image:latest
    SNYK_TOKEN: ${{ secrets.snyk_token }}

- name: Upload Snyk scan results to GitHub Security tab
  id: snyk_upload_sarif
  uses: github/codeql-action/upload-sarif
  if: always()
  with:
    category: "Docker"
    sarif_file: "test-results/snyk.sarif"
```

### Example With Additional Snyk Arguments

```yaml
- name: Snyk Scan (fail on high severity only)
  uses: ministryofjustice/opg-github-actions/actions/snyk-scan
  with:
    SNYK_TEST_IMAGE: my-image:latest
    SNYK_TOKEN: ${{ secrets.snyk_token }}
    snyk_test_args: "--severity-threshold=high"
```

## Inputs

### Required Inputs

- `SNYK_TEST_IMAGE`
  - The Docker image name to test, including tag (for example `my-image:latest`).

- `SNYK_TOKEN`
  - Snyk API token used for authentication.

### Optional Inputs

- `sarif_filepath` (default: `test-results`)
  - Directory where the SARIF output file is written.

- `sarif_filename` (default: `snyk.sarif`)
  - SARIF output file name.

- `snyk_policy_path` (default: `./.snyk`)
  - Path to a `.snyk` policy file to mount into the container.

- `snyk_test_args` (default: `""`)
  - Extra arguments passed directly to `snyk container test`.

## Notes

- The action mounts the policy file into `/root/.snyk` and sets `--policy-path=/root/.snyk`.
- The action checks image identity before and after scanning and fails if the image was modified.
- Ensure your workflow has access to Docker and can run `docker compose`.
