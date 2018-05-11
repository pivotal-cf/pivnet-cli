### Development

For local development, use the following Docker command to build docs and preview on localhost:8000:

`$ docker run --rm -it -p 8000:8000 -v "${PWD}:/docs" squidfunk/mkdocs-material:2.7.2`


For building the site in the `site` folder, use the following Docker command:

`$ docker run --rm -it -v "${PWD}:/docs" squidfunk/mkdocs-material:2.7.2 build`


For publishing the site in the `site` folder, use the following command (prerequisite: `npm install -g gh-pages`):

`gh-pages -d site --message 'Auto-generated commit [#157322604]'`
