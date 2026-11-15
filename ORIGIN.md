# Fork pro origin.cz

Fork [prometheus/blackbox_exporter](https://github.com/prometheus/blackbox_exporter) pro sondy [origin.cz](https://www.origin.cz). Proč a v jaké podobě: ADR 0026 v repu `origin`.

* Větev `origin` je upstream tag plus naše commity. `master` je kopie upstreamu, nepoužívá se.
* Image `ghcr.io/nehoupat/blackbox_exporter:<tag>` (amd64 a arm64) staví workflow `origin.yml` z tagu `vX.Y.Z-origin.N`, kde `vX.Y.Z` je verze upstreamu.
* Workflow upstreamu jsou v tomto repu vypnuté (`gh workflow disable`), tag `v*` by jinak spustil jejich publikování.

## Nová verze upstreamu

```sh
git fetch upstream --tags
git rebase --onto vX.Y.Z vA.B.C origin   # vA.B.C = předchozí upstream tag
go test ./...
git push --force-with-lease origin origin
git tag vX.Y.Z-origin.1 && git push origin vX.Y.Z-origin.1
```

## Lokální build

```sh
docker build -f Dockerfile.origin -t blackbox_exporter:dev .
```
