#!/bin/sh
set -eu

IMAGE="$1"

IDENTITY="$(docker image inspect "$IMAGE" --format '{{if .RepoDigests}}{{index .RepoDigests 0}}{{else}}{{.Id}}{{end}}')"

if [ -z "$IDENTITY" ]; then
  echo "ERROR: Unable to determine image identity for $IMAGE" >&2
  exit 1
fi

echo "$IDENTITY"
