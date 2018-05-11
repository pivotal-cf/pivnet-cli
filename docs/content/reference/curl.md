# Curl an endpoint (aliases: c)

```
Usage:
  pivnet [OPTIONS] curl [curl-OPTIONS] URL

Application Options:
  -v, --version                  Print the version of this CLI and exit
      --format=[table|json|yaml] Format to print as (default: table)
      --verbose                  Display verbose output
      --profile=                 Name of profile (default: default)
      --config=                  Path to config file (default: /Users/pivotal/.pivnetrc)

Help Options:
  -h, --help                     Show this help message

[curl command options]
      -X, --request=             Custom method e.g. PATCH
      -d, --data=                Request data e.g. '{"foo":"bar"}'

[curl command arguments]
  URL:                           URL without host or API prefix e.g. /products/p-mysql/releases/3451

```
