# paramate
CLI tool for AWS Parameter Store

## Install

```bash
$ go install github.com/zisui-sukitarou/paramate@latest
```

## Usage

```bash
$ paramate -h
paramstore is a command line tool for AWS Parameter Store

Usage:
  paramate [flags]
  paramate [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  diff        show difference in envs of the two paths
  help        Help about any command
  show        show envs of the path

Flags:
  -h, --help            help for paramate
  -r, --region string   AWS Region (default "ap-northeast-1")
```

### show

```bash
$ paramate show -h
show envs of the path

Usage:
  paramate show [path] [flags]

Flags:
  -h, --help            help for show
  -r, --region string   AWS region (default "ap-northeast-1")
```

### diff

```bash
$ paramate diff -h
show difference in envs of the two paths

Usage:
  paramate diff [path1] [path2] [flags]

Flags:
  -h, --help            help for diff
  -r, --region string   AWS region (default "ap-northeast-1")
```