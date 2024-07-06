# factorio-server

> [!WARNING]
> WiP

## Usage

```
Usage of /tmp/factorio-server:
  -factorio-token string
    	factorio.com token
  -factorio-user string
    	factorio.com username
  -server-dir string
    	Server base dir and process pwd (default "server")
```

## Build

To build binaries, run:

```sh
goreleaser build --clean --snapshot
```

## Release

To build a local snapshot release, run:

```sh
goreleaser release --clean --snapshot
```

To build and publish a full release, push a semver tag (with 'v' prefix) to the 'main' branch on GitHub.
