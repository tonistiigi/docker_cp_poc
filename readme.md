#### `docker cp` alternative API proof of concept

This repo contains three tools:

#### `docker-pack` 

```
docker-pack /file/or/dir
```

Makes a tar archive from path and outputs to stdout. Equivalent of `GET /containers/id/copy?path=`.

No different formats for files/dirs, all clean tars.


#### `docker-unpack` 

```
docker-pack path/to/unpack < content.tar
```

Extracts tar archive to a destination path from stdin. Equivalent of `PUT /containers/id/copy?path=`.

Files are simply extracted to the specified path. Directories are merged. Destination path must exist. Path suffixes have no special meaning.


#### `docker-cp` 

```
docker-cp /source/path /dest/path
```

Command equivalent with GNU Coreutils `cp -a`.

This tool does not make any requests to user filesystem and only uses the other 2 commands.
