# shim

> A utility for quickly creating symlink shims and adding them to a
> PATH-accessible location.

## Usage

`shim` is a very basic wrapper around `ln`, creating symbolic links to a
"shim directory" (`~/shims`) that can be added to your PATH for easy usage.

```sh
# Link ~/shims/myprogram -> ~/some/path/to/myprogram
$ shim ~/some/path/to/myprogram

# Link ~/shims/abcdefg -> ~/program/with/ugly/name/abcdefg_x64_debug.out
$ shim ~/program/with/ugly/name/abcdefg_x64_debug.out abcdefg

# List shims (equal to `ls ~/shims`)
$ shim
```

## Installation

From the scriptoid root:

```sh
$ ./scriptoid.py link examples/shim
```