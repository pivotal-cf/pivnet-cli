# Create dependency specifier (aliases: cds)

```
Usage:
  pivnet [OPTIONS] create-dependency-specifier [create-dependency-specifier-OPTIONS]

Application Options:
  -v, --version                     Print the version of this CLI and exit
      --format=[table|json|yaml]    Format to print as (default: table)
      --verbose                     Display verbose output
      --profile=                    Name of profile (default: default)
      --config=                     Path to config file (default: /Users/pivotal/.pivnetrc)

Help Options:
  -h, --help                        Show this help message

[create-dependency-specifier command options]
      -p, --product-slug=           Product slug e.g. p-mysql
      -r, --release-version=        Release version e.g. 0.1.2-rc1
      -s, --dependent-product-slug= Dependent product slug e.g. p-mysql
      -u, --specifier=              Specifier e.g. 1.2.*

```      
