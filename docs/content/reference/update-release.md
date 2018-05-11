# Update release (aliases: ur)

```
Usage:
  pivnet [OPTIONS] update-release [update-release-OPTIONS]

Application Options:
  -v, --version                                                                                Print the version of this CLI and exit
      --format=[table|json|yaml]                                                               Format to print as (default: table)
      --verbose                                                                                Display verbose output
      --profile=                                                                               Name of profile (default: default)
      --config=                                                                                Path to config file (default: /Users/pivotal/.pivnetrc)

Help Options:
  -h, --help                                                                                   Show this help message

[update-release command options]
      -p, --product-slug=                                                                      Product slug e.g. p-mysql
      -r, --release-version=                                                                   Release version e.g. 0.1.2-rc1
          --availability=[admins|selected-user-groups|all]                                     Release availability. Optional.
          --release-type=[all-in-one|major|minor|service|maintenance|security|alpha|beta|edge] Release type. Optional.

```
