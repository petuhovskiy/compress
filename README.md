# compress

This is compress utility, which can compress and decompress files.

## Build

```bash
go build .
file ./compress
```

## Usage

```bash
# Get help
./compress

# Compress using Arithmetic coding 
./compress c pushkin.txt cmp.bin ppm

# Compress using best of all (best result, very slow)
./compress c pushkin.txt cmp.bin best

# Decompress
./compress d cmp.bin orig.txt
```

## Run without build
```bash
go run main.go c pushkin.txt cmp.bin mock
```