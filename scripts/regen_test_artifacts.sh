#!/usr/bin/env bash
# Regenerate noir-samples test artifacts in both serialization formats.
#
# For each Nargo.toml under noir-samples/, produces in target/:
#   <name>.json         - default `nargo compile` (MsgpackCompact, byte 3)
#   <name>.tagged.json  - NOIR_SERIALIZATION_FORMAT=msgpack-tagged (byte 4)
#   <name>.gz           - witness from `nargo execute`
#
# Usage: scripts/regen_test_artifacts.sh [path-to-noir-samples]

set -uo pipefail

samples_root="${1:-$(cd "$(dirname "$0")/.." && pwd)/noir-samples}"

if [[ ! -d "$samples_root" ]]; then
    echo "samples root not found: $samples_root" >&2
    exit 1
fi

unset NOIR_SERIALIZATION_FORMAT

regen_one() {
    local nargo_toml="$1"
    local pkg_dir
    pkg_dir="$(dirname "$nargo_toml")"
    local name
    name=$(sed -n 's/^name *= *"\([^"]*\)".*/\1/p' "$nargo_toml" | head -n1)
    if [[ -z "$name" ]]; then
        echo "skip $pkg_dir: no package name" >&2
        return
    fi

    local target="$pkg_dir/target"
    local default_json="$target/$name.json"
    local tagged_json="$target/$name.tagged.json"

    pushd "$pkg_dir" >/dev/null || return

    # Step 1: nargo execute -> default-format .json and witness .gz.
    if ! (unset NOIR_SERIALIZATION_FORMAT; nargo execute --force >/dev/null 2>&1); then
        echo "skip $name: nargo execute failed" >&2
        popd >/dev/null
        return
    fi

    if [[ ! -f "$default_json" ]]; then
        echo "skip $name: $default_json not produced" >&2
        popd >/dev/null
        return
    fi

    # Step 2: stash default .json, recompile in tagged mode, rename.
    mv "$default_json" "$target/$name.default.json"
    if ! NOIR_SERIALIZATION_FORMAT=msgpack-tagged nargo compile --force >/dev/null 2>&1; then
        echo "fail $name: tagged compile failed; restoring default" >&2
        mv "$target/$name.default.json" "$default_json"
        popd >/dev/null
        return
    fi
    mv "$default_json" "$tagged_json"
    mv "$target/$name.default.json" "$default_json"

    echo "ok   $name"
    popd >/dev/null
}

while IFS= read -r -d '' toml; do
    regen_one "$toml"
done < <(find "$samples_root" -name Nargo.toml -print0)
