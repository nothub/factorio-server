# factorio-server

> [!WARNING]
> WiP

## Usage

```
The factory must grow!
Usage: ./factorio-server
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
