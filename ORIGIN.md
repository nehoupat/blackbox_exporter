# Fork for origin.cz

Fork of [prometheus/blackbox_exporter](https://github.com/prometheus/blackbox_exporter) for the [origin.cz](https://www.origin.cz) probes. Why and in what shape: ADR 0026 in the `origin` repository.

* Branch `origin` is an upstream tag plus our commits. `master` mirrors upstream and is not used.
* Image `ghcr.io/nehoupat/blackbox_exporter:<tag>` (amd64 and arm64) is built by the `origin.yml` workflow from a tag `vX.Y.Z-origin.N`, where `vX.Y.Z` is the upstream version.
* Upstream workflows are disabled in this repository (`gh workflow disable`), a `v*` tag would otherwise trigger their release.

## What the fork adds

The code is in package `origin/`, upstream code only gets its wiring in `main.go`.

* **Cache busting:** `/probe?...&cache_buster=_cb` adds `_cb=<16 hex characters>` to the target URL, a new value on every probe. The rest of the URL stays as it is. The name may contain only `A-Za-z0-9_.-`, otherwise the probe returns 400.

## New upstream version

```sh
git fetch upstream --tags
git rebase --onto vX.Y.Z vA.B.C origin   # vA.B.C = previous upstream tag
go test ./...
git push --force-with-lease origin origin
git tag vX.Y.Z-origin.1 && git push origin vX.Y.Z-origin.1
```

## Local build

```sh
docker build -f Dockerfile.origin -t blackbox_exporter:dev .
```
