# version-manager

GitHub Action that creates and pushes the next semver tag of the repository.

```yaml
- uses: rmacciuci/version-manager@v1.8.0
  id: tag
  with:
    minor: true
    prefix: api/   # optional, for monorepos: creates api/v1.2.3
# ${{ steps.tag.outputs.tag }}      -> api/v1.2.3 (also in ./VERSION)
# ${{ steps.tag.outputs.version }}  -> v1.2.3 (Docker image tag)
```

- Without `prefix`, only plain `vX.Y.Z` tags are considered.
- Several bumps apply in order major, minor, patch (`minor`+`patch` on v1.4.2 gives v1.5.1).

Development: `go test ./...` (Go 1.26).
