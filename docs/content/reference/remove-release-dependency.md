# Remove release dependency (aliases: rrd)

```
Usage:
  pivnet [OPTIONS] remove-release-dependency [remove-release-dependency-OPTIONS]

Application Options:
  -v, --version                        Print the version of this CLI and exit
      --format=[table|json|yaml]       Format to print as (default: table)
      --verbose                        Display verbose output
      --profile=                       Name of profile (default: default)
      --config=                        Path to config file (default: /Users/pivotal/.pivnetrc)

Help Options:
  -h, --help                           Show this help message

[remove-release-dependency command options]
      -p, --product-slug=              Product slug e.g. p-mysql
      -r, --release-version=           Release version e.g. 0.1.2-rc1
      -s, --dependent-product-slug=    Dependent product slug e.g. p-mysql
      -u, --dependent-release-version= Dependent release version e.g. 0.1.2-rc1

```
