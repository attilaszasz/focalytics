#!/usr/bin/env sh

set -eu

version=${1:?version is required}
goos=${2:?goos is required}
goarch=${3:?goarch is required}
out_dir=${4:?output directory is required}

repo_root=$(CDPATH= cd -- "$(dirname "$0")/../.." && pwd)
src_dir="$repo_root/src"

archive_name=$(cd "$src_dir" && go run ./tools/releasegen asset-name --version "$version" --goos "$goos" --goarch "$goarch" --field archive)
binary_name=$(cd "$src_dir" && go run ./tools/releasegen asset-name --version "$version" --goos "$goos" --goarch "$goarch" --field binary)

build_dir=$(mktemp -d "${TMPDIR:-/tmp}/focalytics-release.XXXXXX")
cleanup() {
  rm -rf "$build_dir"
}
trap cleanup EXIT INT TERM

mkdir -p "$out_dir"

(
  cd "$src_dir"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$build_dir/$binary_name" .
)

case "$archive_name" in
  *.zip)
    (
      cd "$build_dir"
      zip -q "$out_dir/$archive_name" "$binary_name"
    )
    ;;
  *.tar.gz)
    (
      cd "$build_dir"
      tar -czf "$out_dir/$archive_name" "$binary_name"
    )
    ;;
  *)
    echo "unsupported archive extension for $archive_name" >&2
    exit 1
    ;;
esac

printf '%s\n' "$out_dir/$archive_name"