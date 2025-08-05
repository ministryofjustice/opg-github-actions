# 7. Terraform Version Action Changes

Date: 2025-07-24

## Status

Accepted

## Context

Before the community migrated to using our standard `.envrc` file that handles terraform version switching as well and environment variable setup some projects made use of a `.terrform-version` file, as this was supported by `tfswtich` and could be read within pipelines easily.

However, having this file as well as the value in `required_version` property in actual terraform files lead to drift and inconsistencies, so the community chose to move away from that approach and instead parse the contents of `versions.tf` file so ther was only one location for that configuration.

For a while, the `terraform-version` action maintained backward compatability with the older `.terraform-version` plain text file for easier adoption.

## Decision

Our code has moved to use the community approach, so this version will no longer support the `.terraform-version` file.

## Consequences

If your pipeline relies the `.terraform-version` file, then you will need to update you setup to use `versions.tf` approach.
