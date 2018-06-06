# Lifecycle

A reference implementation of Buildpack API v3.

## Commands

### Build

* `detector` - chooses buildpacks (via `/bin/detect`)
* `analyzer` - restores launch layer metadata from the previous build
* `builder` -  executes buildpacks (via `/bin/build`)
* `exporter` - remotely patches images with new layers (via rebase & append)
* `launcher` - invokes choice of process

### Develop

* `detector` - chooses buildpacks (via `/bin/detect`)
* `developer` - executes buildpacks (via `/bin/develop`)
* `launcher` - invokes choice of process

### Cache

* `restorer` - restores cache
* `cacher` - updates cache

## Notes

Only the `detector` and `builder` are currently implemented.

The `analyzer` and `exporter` are partially implemented in [packs](https://github.com/sclevine/packs).

Cache implementations (`restorer` and `cacher`) are intended to be interchangable and platform-specific.
A platform may choose not to deduplicate cache layers for efficiency.