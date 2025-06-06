[![Netlify Status][netifly-img]][netifly]

# Jaeger website

This repo houses all the assets used to build the Jaeger website, available at https://jaegertracing.io.

The site is built with [Hugo](https://gohugo.io/) and hosted by [Netlify](https://www.netlify.com/).

## Setup

Install the "extended" Hugo binary from [hugo/releases](https://github.com/gohugoio/hugo/releases) (use one of the `hugo_extended_*` binaries) or
use a package manager if it is available for your operating system.

>  The "extended" version of Hugo supports [Sass](https://sass-lang.org), which is necessary to build the site locally.

The currently used version of Hugo is defined in the [`netlify.toml`](./netlify.toml) configuration file.

Install the active LTS version of Node.js, then run the following command from the directory of this repo's clone:

```bash
npm install
```

## Running the site locally

If you want to develop the site locally, you can run a single command (assuming that you've run the [setup](#setup)):

```bash
make develop
```

This will start up a local server on localhost port 1313. When you make changes to either the content of the website (in [`content`](content)) *or* to the Sass and JavaScript assets of the page (in [`themes/jaeger-docs/assets`](themes/jaeger-docs/assets)), the browser will automatically update to reflect those changes (usually in under one second).

## Publishing the site

The site is published automatically by [Netlify](https://www.netlify.com/) whenever changes are merged to the `main` branch. The site cannot be published in an ad-hoc way (e.g. through a `make` command or script in the repo).

## Contributing to the site

We strongly encourage you to contribute to this site! For more information, see the [contributing](CONTRIBUTING.md) guide.

## Diagrams

Diagrams included in the documentation are created in the shared [Google Slides document][slides], which supports export to SVG. If you need to make changes to the diagrams as part of a PR, please copy the diagram into a new slide deck and include a shared link to it in the PR along with the exported SVG file. The maintainers will update the main deck with the new version upon merging the PR.

## Publishing new Jaeger version

Please refer to [RELEASE.md](./RELEASE.md) for instructions on how to release new version of documentation.

## Admonitions

There are five admonition types available for the Jaeger docs:

Admonition type | Color
:---------------|:-----
`info` | blue
`success` | green
`danger` | red
`warning` | yellow
`requirement` | purple

Here's an example:

```markdown
{{< danger >}}
We do not recommend that you do this!
{{< /danger >}}
```

You can also add titles:

```markdown
{{< success title="New feature" >}}
Jaeger now supports a new thing that you definitely want.
{{< /success >}}
```

## Link checking

You can check internal links by running `make check-internal-links` and all links, including external links, by running `make check-all-links`.

## Adding new pages and redirects

When new pages are added to the documentation, please add a corresponding entry to [themes/jaeger-docs/layouts/index.redirects](./themes/jaeger-docs/layouts/index.redirects).

## Generating Roadmap page

To generate the `content/roadmap.md` document, run `make fetch-roadmap`.

This script fetches issues from the [GitHub project board](https://github.com/orgs/jaegertracing/projects/4/views/1?layout=table), extracts the required information, and generates the roadmap document. Make sure to set the `GITHUB_TOKEN` environment variable with your GitHub API token before running the script, or save the token in `~/.github_token` file (protect the file so only you can read it: `chmod 0600 <file>`). Personal tokens can be created at https://github.com/settings/tokens/.

## Updating Medium Blog Feed

The homepage displays the latest blog posts from the [Jaeger Medium blog](https://medium.com/jaegertracing).
To avoid network calls during builds and to ensure fast, reliable local development, the Medium RSS feed is downloaded and stored as a static XML file.

### Prerequisites

Ensure you have [`curl`](https://curl.se/) installed on your system to download the RSS feed.

Most Linux and macOS systems already have `curl` pre-installed.
You can verify installation by running:

```bash
curl --version
```

### To update the Medium blog feed
Run the following command to fetch and save the feed locally as XML:

```bash
make fetch-blog-feed
```

This will download the RSS feed from Medium and save it to:
```bash
assets/data/medium.xml
```

After updating, commit the changes:
```bash
git add assets/data/medium.xml
git commit -m "chore: update Medium blog feed"
```

### How it's used
The Hugo site reads and parses data/medium.xml using resources.Get and transform.Unmarshal.
This converts the XML into structured data at build time, allowing full control over the content without relying on live network requests.

## License

[Apache 2.0 License](./LICENSE).

[slides]: https://docs.google.com/presentation/d/1JuurkQn03z0BbOEAViJBEE_WWMj6JQUML-uJm7zizvI/
[netifly-img]: https://api.netlify.com/api/v1/badges/d2b1a1ea-f454-4ba8-990c-cc469c959556/deploy-status
[netifly]: https://app.netlify.com/sites/jaegertracing/deploys
[htmltest]: https://github.com/wjdp/htmltest
