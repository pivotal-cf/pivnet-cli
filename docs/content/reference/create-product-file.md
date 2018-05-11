# Create product file (aliases: cpf)

```
Usage:
  pivnet [OPTIONS] create-product-file [create-product-file-OPTIONS]

Application Options:
  -v, --version                  Print the version of this CLI and exit
      --format=[table|json|yaml] Format to print as (default: table)
      --verbose                  Display verbose output
      --profile=                 Name of profile (default: default)
      --config=                  Path to config file (default: /Users/pivotal/.pivnetrc)

Help Options:
  -h, --help                     Show this help message

[create-product-file command options]
      -p, --product-slug=        Product slug e.g. 'p-mysql'
          --name=                Name e.g. 'p-mysql 1.7.13'
          --aws-object-key=      AWS Object Key e.g. 'product_files/P-MySQL/p-mysql-1.7.13.pivotal'
          --file-type=           File Type e.g. 'Software'
          --file-version=        File Version e.g. '1.7.13'
          --md5=                 MD5 of file
          --description=         Description of file
          --docs-url=            URL of docs for file
          --included-file=       Name of included file
          --platform=            Platform of file
          --released-at=         When file is marked for release e.g. '2016/01/16'
          --system-requirement=  System-requirement of file

```
