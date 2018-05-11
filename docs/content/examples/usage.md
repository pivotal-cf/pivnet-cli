# Usage

Using the Pivnet CLI requires a valid `Pivotal Network API token` or `UAA Refresh Token`.

Refer to the
[official docs](https://network.pivotal.io/docs/api#how-to-authenticate)
for more details on obtaining a Pivotal Network API token.

Example usage:

```sh
$ pivnet login --api-token='my-api-token'
$ pivnet products

+-----+------------------------------------------------------+--------------------------------+
| ID  |                         SLUG                         |              NAME              |
+-----+------------------------------------------------------+--------------------------------+
|  60 | elastic-runtime                                      | Pivotal Cloud Foundry Elastic  |
|     |                                                      | Runtime                        |
+-----+------------------------------------------------------+--------------------------------+

$ pivnet r -p elastic-runtime -r 1.8.8 --format json \
  | jq '{"id": .id, "release_date": .release_date, "release_type": .release_type}'

{
  "id": 2555,
  "release_date": "2016-10-13",
  "release_type": "Security Release"
}
```
