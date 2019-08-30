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

$ pivnet release --product-slug=elastic-runtime --release-version=1.8.8 --format json \
  | jq '{"id": .id, "release_date": .release_date, "release_type": .release_type}'

{
  "id": 2555,
  "release_date": "2016-10-13",
  "release_type": "Security Release"
}
```

# Batch Command Examples

The Pivnet UI has few places to edit content in batches.  Batch processing is relegated to the [Pivnet Resource](https://github.com/pivotal-cf/pivnet-resource) and Pivnet CLI.

Many commands require release versions.  You could use the following pipeline to product a list of release versions:

```sh
$ pivnet --format=json releases --product-slug=MY_PRODUCT_SLUG  | jq '.[].version'
"1.10.5"
"1.10.4"
"1.10.3"
"1.10.2"
"1.10.1"
"1.10.0"
...
```

If you wanted to, for example, a remove a user group from all releases, you could use a pipeline similar to the following:

```sh
$ pivnet --format=json releases --product-slug=MY_PRODUCT_SLUG | jq '.[].version' | xargs -I{} pivnet remove-user-group --product-slug=MY_PRODUCT_SLUG --release-version={} --user-group-id=USER_GROUP_ID_TO_REMOVE
```

Similarly, you could, for example, add an Elastic Runtime 2.0.0 release dependency to all "1.10.*" releases for a product with a pipeline similar to the following:

```sh
$ pivnet --format=json releases --product-slug=MY_PRODUCT_SLUG | jq -r '.[].version' | grep '^1\.10\.' | xargs -I{} pivnet add-release-dependency --product-slug=MY_PRODUCT_SLUG --release-version={} --dependent-product-slug=elastic-runtime --dependent-release-version=2.0.0
```

