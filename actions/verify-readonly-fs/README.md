# Verify Read-only Filesystem

Checks that `read_only` and `tmpfs` configuration in a docker compose file matches the ECS mounts config for each mapped service. Fails with a per-service error list when any mismatch is found.

For each service the action verifies:

- `read_only: true` is set
- every tmpfs path/size from the mounts config is present in compose with the correct byte size
- compose has no extra tmpfs paths not present in the mounts config

## Usage

Intended for use with mounts configuration stored in SSM but any mapped json will do.

```yaml
- name: Fetch mounts config
  id: ssm
  run: |
    mounts_json=$(aws ssm get-parameter \
      --name "/<service>/<environment>/ecs-mounts" \
      --region eu-west-1 \
      --query "Parameter.Value" \
      --output text)
    echo "mounts_json<<EOF" >> "$GITHUB_OUTPUT"
    echo "$mounts_json" >> "$GITHUB_OUTPUT"
    echo "EOF" >> "$GITHUB_OUTPUT"

- name: Verify read-only filesystem config
  uses: 'ministryofjustice/opg-github-actions/actions/verify-readonly-fs@<SHA> # <version>'
  with:
    mounts_json: ${{ steps.ssm.outputs.mounts_json }}
    service_json: |
      {
        "api-app":   "api-app",
        "queue":     "queue-consumer",
        "api":       "nginx",
        "frontend":  "nginx",
        "finance-hub": null
      }
```

## Inputs

#### `service_json` (required)

JSON object mapping each docker compose service name to its key in `mounts_json`.

The same key can be reused for multiple services (e.g. `"nginx"` for both `api` and `frontend`).

Set the value to `null` for services that have no tmpfs mounts — only `read_only: true` is verified for those services, and any unexpected tmpfs entries in compose are flagged.

#### `mounts_json` (required)

JSON object containing the ECS tmpfs mount config, keyed by service name.

Structure:
```json
{
  "api-app":  { "tmpfs": [{ "containerPath": "/tmp", "size": 64 }] },
  "nginx":    { "tmpfs": [{ "containerPath": "/var/cache/nginx", "size": 64 }, { "containerPath": "/tmp", "size": 64 }] }
}
```

`size` is in **MiB**. The action converts to bytes for comparison against `docker compose config` output.

#### `compose_file`

Path to the docker compose file to verify. Default: `docker-compose.yml`.

## Outputs

#### `passed`

`true` when every service matched, `false` otherwise.

#### `errors`

Newline-delimited list of mismatches. Empty when `passed` is `true`.

## Error cases

| Scenario | Error |
|---|---|
| `read_only` not set | `<svc>: missing read_only: true` |
| tmpfs path missing from compose | `<svc>: missing tmpfs path <path>` |
| tmpfs size wrong | `<svc>: <path> size <actual> bytes != expected <expected> bytes` |
| unexpected tmpfs path in compose | `<svc>: unexpected tmpfs path <path> ...` |
| `tmpfs:` key present but empty | `<svc>: compose config invalid: ...tmpfs must be a string` |
| key not found in `mounts_json` | `<svc>: key '<key>' not found in mounts config` |

## Requirements

`docker` (compose plugin) and `jq` must be available on the runner. Both are present on the default GitHub-hosted Ubuntu runners.
