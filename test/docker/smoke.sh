#!/bin/sh
# Checks that the Docker image can run a typed rule, as root and as a non-root user.
set -eu

image="${1:?usage: smoke.sh IMAGE}"
dir="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
want='main.go:11:14: use t1.Equal(t2) instead of "==" operator'

for user in root 1000:1000; do
	status=0
	got="$(docker run --rm --user "$user" --volume "$dir":/app --workdir /app "$image" \
		-config revive.toml -set_exit_status ./...)" || status=$?

	if [ "$got" != "$want" ]; then
		printf 'as %s: got output %s, want %s\n' "$user" "${got:-<none>}" "$want" >&2
		exit 1
	fi
	if [ "$status" -ne 1 ]; then
		printf 'as %s: got exit status %s, want 1\n' "$user" "$status" >&2
		exit 1
	fi
done
