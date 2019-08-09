# gyroid [![CircleCI](https://circleci.com/gh/Coteh/gyroid.svg?style=shield)](https://circleci.com/gh/Coteh/gyroid)

A [Pocket](https://getpocket.com/) client designed to help you quickly refine your Pocket library.

![Preview](Preview.gif "Preview Image")

## Motivation

I have an ever-growing collection of articles saved on Pocket that remain untagged, and thus difficult to track down when I want to read them. Also, I find that I am more likely to read articles on the top of my Pocket list, which are the more recently added items, than ones that are in the middle or closer to the bottom. I would like to have a quick and easy way of scanning through my untagged items and either assign them to appropriate tags or bump them back up to the top of my list.

Luckily, Pocket has a number of great tools available in its [public API](https://getpocket.com/developer/docs/overview) that I can use to achieve these tasks, plus more. I put together this CLI tool to have a simple tool available to accomplish these refinement tasks much more efficiently than on the browser webapp. Quickly assign tags, bump, mark as favourite, and archive articles! All from the comfort of the terminal. Keep your articles organized before you're on-the-go.

## Features

- Mark articles with one or more tags
- Favourite/unfavourite articles
- Bump articles to the top of the list (the same action as unarchiving an article on Pocket)
- Archive an article
- Delete an article
- Open an article on your web browser directly from the CLI
- Add an article to your Pocket list, either by manual input or directly from clipboard

## Installation

The quickest way to get started with `gyroid` is to run the following command to install it to your GOPATH:

```
go get -u github.com/Coteh/gyroid
```

You will need a Pocket API consumer key to use `gyroid`. You can grab one [here](https://getpocket.com/developer/apps/new).

You will then need to create a textfile in `~/.config/gyroid` (where `~` is your home directory) called `consumer_key` and simply paste the consumer key into it.

Now you can run the program from your terminal (assuming GOPATH is in your PATH) as follows:

```
gyroid
```

See this project's [`maskfile.md`](maskfile.md) for more commands to build and run the project. You will need the latest version of [Go](https://golang.org/). (1.11 or newer as this project uses the new versioned module system)

*The commands within the `maskfile.md` can also be executed using [`mask`](https://github.com/jakedeichert/mask). If you don't have mask installed, you can install it from [here](https://crates.io/crates/mask) using `cargo` (requires [Rust toolchain](https://rustup.rs/))*

## Configuration
The configuration file for `gyroid` will be found in `~/.config/gyroid/config.yml`, where `~` represents your home directory. `gyroid` will populate one on startup if not found, or you can define it yourself using the following template:

```yaml
clipboard: false     # whether or not to use clipboard (disabled by default)
```

## Issues
- Number of requests made and number of articles to be returned from each request need to be fine-tuned [#5](../../issues/5)
- Certain edge cases of Pocket API (such as Retrieve request with offset past final articles on list) need to be handled [#10](../../issues/10)
- Check out [Issues](../../issues) to post your own

## Future Additions
- Customize ordering of Pocket articles to be loaded [#8](../../issues/8)
- Quick tag system [#6](../../issues/6)
- Consider switching to [go-pocket](https://github.com/motemen/go-pocket) for Pocket API interaction [#7](../../issues/7)
