# Pivnet CLI

Interact with [Pivotal Network](https://network.pivotal.io) from the command-line.

## Developing

### Prerequisites

A valid install of golang >= 1.6 is required.

### Dependencies

Dependencies are vendored in the `vendor` directory, according to the
[golang 1.5 vendor experiment](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&uact=8&ved=0ahUKEwi7puWg7ZrLAhUN1WMKHeT4A7oQFggdMAA&url=https%3A%2F%2Fgolang.org%2Fs%2Fgo15vendor&usg=AFQjCNEPCAjj1lnni5apHdA7rW0crWs7Zw).

No action is require to fetch the vendored dependencies.

### Running the tests

Install the ginkgo executable with:

```
go get -u github.com/onsi/ginkgo/ginkgo
```

The tests require a valid Pivotal Network API token and host.

Refer to the
[official docs](https://network.pivotal.io/docs/api#how-to-authenticate)
for more details on obtaining a Pivotal Network API token.

It is advised to run the acceptance tests against the Pivotal Network integration
environment endpoint i.e. `HOST='https://pivnet-integration.cfapps.io'`.

Run the tests with the following command:

```
API_TOKEN=my-token \
HOST='https://pivnet-integration.cfapps.io' \
./bin/test_all
```

### Contributing

Please make all pull requests to the `develop` branch, and
[ensure the tests pass locally](https://github.com/pivotal-cf/pivnet-cli#running-the-tests).

### Project management

The CI for this project can be found at https://sunrise.ci.cf-app.com and the
scripts can be found in the
[pivnet-resource-ci repo](https://github.com/pivotal-cf/pivnet-resource-ci).

The roadmap is captured in [Pivotal Tracker](https://www.pivotaltracker.com/projects/1474244).
