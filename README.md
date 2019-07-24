# thc

thc (Temporal HTTP Client) is a Golang library for interacting with Temporal's HTTP API. Additionally it includes a command line tool that lets you use a subset of the API from the command line! This makes it extremely easy to use Temporal as a backup tool, or even for easily hosting your website on IPFS.

# Download (CLI)

We curently distribute pre-built binaries in a wide variety of platforms, including:
* linux (32+64 bit)
* arm
* darwin/mac (32+64 bit)
* windows (32+64 bit)

To access the prebuilt binaries, including the sha256 checksums, please see [here](https://gateway.temporal.cloud/ipfs/QmVPdxNGFg1drdFmrwcih8zzGKd4ywti6izLp2jvf6KY1L)

# Usage (CLI)

To use the CLI, regardless of the command you're using, you must provide the flags `--user.name`, and `--user.pass` for the Temporal account that you want to use. The values for these flags are used to authenticate with the API, and generate the JWT used to authenticate calls like pins, lens indexing, etc...

Currently the supported commands are:

* File uploads
  * This is used to uplaod files, and is essentially a wrapper around `ipfs add`
* Directory uploads
  * This is used to add directories, and is essentially a wrapper around `ipfs add -r`
* Pin adds
  * This is used to add a pin to Temporal's IPFS nodes
* Pin extend
  * This is used to extend the pin duration of your ipfs hashes
* Lens search
  * This is used to search the lens search engine
  * Please note this command is very experimental, and the output is quite  ugly
* Lens index
  * This is used to submit content for indexing by the lens search engine

## Examples

Upload a directory called `release`, and recursively pin the directory hash for 18 months

```
thc --user.name foo --user.pass bar upload dir --dir release --hold.time 18
```

Pin a hash `QmCheese` for 2 months

```
thc --user.name foo --user.pass bar pin --hash QmCheese --hold.time 2
```