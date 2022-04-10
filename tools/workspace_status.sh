#!/bin/bash
set -e -E -o pipefail

if [ "$CI" != "true" ]; then
  GITHUB_REPOSITORY=$(basename $(git rev-parse --show-toplevel))
  GITHUB_SHA=$(git rev-parse --short HEAD)
  GITHUB_REF_NAME=$(git rev-parse --abbrev-ref HEAD)
fi

set -u

GITHUB_REPOSITORY=$GITHUB_REPOSITORY
echo "STABLE_GIT_REPO_SLUG $GITHUB_REPOSITORY"

COMMIT_HASH=$GITHUB_SHA
echo "COMMIT_HASH $COMMIT_HASH"

echo "GITHUB_REF_NAME $GITHUB_REF_NAME"
