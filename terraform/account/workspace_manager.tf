module "workspace_cleanup" {
  # previous: github.com/TomTucka/terraform-workspace-manager/terraform/workspace_cleanup
  source  = "github.com/ministryofjustice/opg-terraform-workspace-manager//terraform/workspace_cleanup?ref=v0.3.4"
  enabled = true
}
