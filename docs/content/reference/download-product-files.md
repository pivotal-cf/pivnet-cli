# Download product files (aliases: dlpf)

```
Usage:
  pivnet [OPTIONS] download-product-files [download-product-files-OPTIONS]

Application Options:
  -v, --version                  Print the version of this CLI and exit
      --format=[table|json|yaml] Format to print as (default: table)
      --verbose                  Display verbose output
      --profile=                 Name of profile (default: default)
      --config=                  Path to config file (default: /Users/pivotal/.pivnetrc)

Help Options:
  -h, --help                     Show this help message

[download-product-files command options]
      -p, --product-slug=        Product slug e.g. p-mysql
      -r, --release-version=     Release version e.g. 0.1.2-rc1
      -i, --product-file-id=     Product file ID e.g. 1234
      -g, --glob=                Glob to match product name e.g. *aws*
      -d, --download-dir=        Local existing directory to download files to e.g. /tmp/my-file/ (default: .)
          --accept-eula          Automatically accept EULA if necessary (available to pivots only)

```
