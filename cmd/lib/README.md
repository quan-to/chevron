Chevron Shared Library
======================

## Build Linux

```bash
go build -o chevron.so -buildmode=c-shared
```

## Build MacOSX

```bash
go build -o chevron.dylib -buildmode=c-shared
```

## Build Windows

```bash
GOARCH=386 go build -o chevron32.dll -buildmode=c-shared
GOARCH=amd64 go build -o chevron.dll -buildmode=c-shared
```

### Cross-Compile for Windows:

```bash
#!/bin/bash
CC=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o chevron32.dll -buildmode=c-shared
CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o chevron.dll -buildmode=c-shared
```
