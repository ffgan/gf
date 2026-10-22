#!/usr/bin/env bash
# Build gf for every platform supported by the Go toolchain on PATH.
#
# Release CI installs the oldest Go release that is still maintained
# (see .go-version) and runs this script. The platform list is that
# toolchain's `go tool dist list`, plus the ABI variants that decide
# which machines can actually run the binary:
#   - 32-bit ARM: GOARM=5, 6 and 7 (Android stays armv7)
#   - MIPS: hardfloat and softfloat
#
# android/386, android/amd64 and android/arm need an NDK (external
# linking). Set ANDROID_NDK_HOME. ios/arm64 needs macOS with Xcode and
# is skipped on other hosts unless it was requested explicitly.
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

version="${VERSION:-0.2.0}"
version="${version#v}"
dist="${DIST_DIR:-$root/dist}"
mode=build
only_csv=""

usage() {
	cat <<EOF
Usage: VERSION=${version} ANDROID_NDK_HOME=... $0 [--print-targets | --expect-dist] [goos/goarch ...]

  (no args)        build every target this host can produce
  goos/goarch      build only those dist-list pairs (ARM and MIPS variants included)
  --print-targets  print archive names, do not build
  --expect-dist    exit 1 if dist/ is missing any archive
EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
	--print-targets) mode=print ;;
	--expect-dist) mode=expect ;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		only_csv="${only_csv} $1"
		;;
	esac
	shift
done

if ! command -v go >/dev/null 2>&1; then
	echo "go is not on PATH" >&2
	exit 1
fi

selected() {
	[[ -z "${only_csv// /}" ]] && return 0
	[[ " ${only_csv} " == *" $1 "* ]]
}

# Lines are: goos, goarch, label, kind, extra-env.
# kind is pure, ios-amd64, cgo-android, or cgo-ios.
list_targets() {
	local pair goos goarch v fp
	while read -r pair; do
		[[ -n "$pair" ]] || continue
		goos=${pair%/*}
		goarch=${pair#*/}
		case "$goarch" in
		arm)
			if [[ "$goos" == android ]]; then
				printf '%s\t%s\t%s\t%s\t%s\n' android arm armv7 cgo-android GOARM=7
			else
				for v in 5 6 7; do
					printf '%s\t%s\t%s\t%s\t%s\n' "$goos" arm "armv${v}" pure "GOARM=${v}"
				done
			fi
			;;
		mips | mipsle)
			for fp in hardfloat softfloat; do
				printf '%s\t%s\t%s\t%s\t%s\n' "$goos" "$goarch" "${goarch}_${fp}" pure "GOMIPS=${fp}"
			done
			;;
		mips64 | mips64le)
			for fp in hardfloat softfloat; do
				printf '%s\t%s\t%s\t%s\t%s\n' "$goos" "$goarch" "${goarch}_${fp}" pure "GOMIPS64=${fp}"
			done
			;;
		*)
			local kind=pure
			case "$goos/$goarch" in
			android/386 | android/amd64) kind=cgo-android ;;
			ios/amd64) kind=ios-amd64 ;;
			ios/arm64) kind=cgo-ios ;;
			esac
			printf '%s\t%s\t%s\t%s\t%s\n' "$goos" "$goarch" "$goarch" "$kind" ""
			;;
		esac
	done < <(go tool dist list)
}

archive_name() {
	local goos=$1 label=$2 ext=tar.gz
	[[ "$goos" == windows ]] && ext=zip
	printf 'gf_%s_%s_%s.%s' "$version" "$goos" "$label" "$ext"
}

ndk_bin() {
	local pre="${ANDROID_NDK_HOME}/toolchains/llvm/prebuilt" host
	case "$(uname -s)-$(uname -m)" in
	Linux-x86_64) host=linux-x86_64 ;;
	Linux-aarch64 | Linux-arm64) host=linux-aarch64 ;;
	Darwin-x86_64) host=darwin-x86_64 ;;
	Darwin-arm64) host=darwin-arm64 ;;
	*)
		echo "no NDK host prebuilt for $(uname -s)-$(uname -m)" >&2
		exit 1
		;;
	esac
	if [[ ! -x "$pre/$host/bin/clang" ]]; then
		echo "NDK clang not found at $pre/$host/bin/clang" >&2
		exit 1
	fi
	printf '%s\n' "$pre/$host/bin"
}

android_cc() {
	local triple
	case "$1" in
	386) triple=i686-linux-android21-clang ;;
	amd64) triple=x86_64-linux-android21-clang ;;
	arm) triple=armv7a-linux-androideabi21-clang ;;
	*)
		echo "no Android compiler for $1" >&2
		exit 1
		;;
	esac
	printf '%s/%s\n' "$(ndk_bin)" "$triple"
}

ios_cc() {
	local wrap="$workdir/ios-clang" sdk clang
	if [[ -x "$wrap" ]]; then
		printf '%s\n' "$wrap"
		return
	fi
	sdk=$(xcrun --sdk iphoneos --show-sdk-path)
	clang=$(xcrun --sdk iphoneos --find clang)
	cat >"$wrap" <<EOF
#!/bin/sh
exec $(printf '%q' "$clang") -arch arm64 -isysroot $(printf '%q' "$sdk") -mios-version-min=12.0 "\$@"
EOF
	chmod +x "$wrap"
	printf '%s\n' "$wrap"
}

if [[ "$mode" == print || "$mode" == expect ]]; then
	missing=0
	while IFS=$'\t' read -r goos goarch label kind extra; do
		name=$(archive_name "$goos" "$label")
		if [[ "$mode" == print ]]; then
			printf '%s\n' "$name"
			continue
		fi
		if [[ ! -f "$dist/$name" ]]; then
			echo "missing $name" >&2
			missing=1
		fi
	done < <(list_targets)
	exit "$missing"
fi

workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT
rm -rf "$dist"
mkdir -p "$dist"

ldflags="-s -w -X github.com/ffgan/gf/cmd.Version=${version}"
built=0

while IFS=$'\t' read -r goos goarch label kind extra; do
	pair="$goos/$goarch"
	selected "$pair" || continue

	if [[ "$kind" == cgo-ios && "$(uname -s)" != Darwin ]]; then
		if [[ -n "${only_csv// /}" ]]; then
			echo "ios/arm64 has to be built on macOS with Xcode" >&2
			exit 1
		fi
		echo "skip ios/arm64 (needs macOS and Xcode)"
		continue
	fi
	if [[ "$kind" == cgo-android && -z "${ANDROID_NDK_HOME:-}" ]]; then
		echo "ANDROID_NDK_HOME is required to build $pair" >&2
		exit 1
	fi

	binname=gf
	[[ "$goos" == windows ]] && binname=gf.exe
	[[ "$goarch" == wasm ]] && binname=gf.wasm

	stage="$workdir/stage"
	rm -rf "$stage"
	mkdir -p "$stage"
	binpath="$stage/$binname"

	# Keep the go-build argv non-empty for bash 3.2 (macOS /bin/bash) under set -u.
	args=(go build -trimpath -ldflags "$ldflags" -o "$binpath" .)
	case "$kind" in
	ios-amd64) args=(go build -trimpath -buildmode=exe -ldflags "$ldflags" -o "$binpath" .) ;;
	cgo-ios) args=(go build -trimpath -buildmode=pie -ldflags "$ldflags" -o "$binpath" .) ;;
	esac

	run=(env -u CC -u CXX)
	case "$kind" in
	cgo-android)
		run+=(CGO_ENABLED=1 CC="$(android_cc "$goarch")")
		;;
	cgo-ios)
		run+=(CGO_ENABLED=1 CC="$(ios_cc)")
		;;
	*)
		run+=(CGO_ENABLED=0)
		;;
	esac
	run+=(GOOS="$goos" GOARCH="$goarch")
	[[ -n "$extra" ]] && run+=("$extra")

	echo "build ${goos}/${label}"
	"${run[@]}" "${args[@]}"

	cp "$root/LICENSE" "$root/README.md" "$stage/"
	name=$(archive_name "$goos" "$label")
	if [[ "$goos" == windows ]]; then
		python3 - "$dist/$name" "$binname" "$stage" <<'PY'
import os, sys, zipfile
dest, binname, stage = sys.argv[1], sys.argv[2], sys.argv[3]
with zipfile.ZipFile(dest, "w", compression=zipfile.ZIP_DEFLATED) as zf:
    for fn in (binname, "LICENSE", "README.md"):
        zf.write(os.path.join(stage, fn), fn)
PY
	else
		tar -C "$stage" -czf "$dist/$name" "$binname" LICENSE README.md
	fi
	built=$((built + 1))
done < <(list_targets)

if [[ "$built" -eq 0 ]]; then
	echo "no targets built" >&2
	exit 1
fi

host_archive="$dist/gf_${version}_linux_amd64.tar.gz"
if [[ "$(uname -s)" == Linux && -f "$host_archive" ]]; then
	smoke="$workdir/smoke"
	mkdir -p "$smoke"
	tar -C "$smoke" -xzf "$host_archive"
	got=$("$smoke/gf" --version | tr -d '\r')
	if [[ "$got" != "gf ${version}" ]]; then
		echo "linux/amd64 --version = ${got}, want gf ${version}" >&2
		exit 1
	fi
fi

echo "built ${built} archives in ${dist}"
