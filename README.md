# Dog

Dog is a command line application that executes tasks. It works in a similar way as GNU Make or ruby's Rake but it is a more generic task runner, not a build tool. It can be used as a layer on top of your Makefile or your shell scripts. Dog's default script syntax is `sh` but it also supports BASH, Python or Ruby so you can write your tasks in any language.

## What is a Dogfile?

Dogfile is a specification that uses YAML to describe the tasks related to a project. We think that the Spec will be finished (no further breaking changes) by the v1.0 version of Dog.

- Read the [Dogfile Spec](https://github.com/dogtools/dog/blob/master/DOGFILE_SPEC.md)
- Read [Dog's own Dogfile](https://github.com/dogtools/dog/blob/master/Dogfile.yml)

## Other tools

Our name for tools that use Dogfiles is *dogtools*. Dog is the first dogtool but there is a lot more that can be done: web and desktop UIs, chat bot interfaces, plugins for text editors and IDEs, tools to export Dogfiles to other formats, HTTP API interfaces, even implementations of the cli in other languages! To simplify the process of creating dogtools we are implementing parts of Dog as Go packages so you can import them in your project (see [parser](https://github.com/dogtools/dog/tree/master/parser), [types](https://github.com/dogtools/dog/tree/master/types) and [execute](https://github.com/dogtools/dog/tree/master/execute)). Let us know if you have any uncovered need one of these packages.

## Contributing

At this moment we are focused on implementing the basics that will allow us to publish v0.1. This project is organized using GitHub [Issues](https://github.com/dogtools/dog/issues) and [Pull Requests](https://github.com/dogtools/dog/pulls).

If you want to help, take a look at:

- Open [bugs](https://github.com/dogtools/dog/issues?q=is%3Aissue+is%3Aopen+label%3Abug)
- Lacking features for [v0.2.0](https://github.com/dogtools/dog/milestone/2)
- Lacking features for [v0.3.0](https://github.com/dogtools/dog/milestone/3)
- Our [Code of Conduct](https://github.com/dogtools/dog/blob/master/CODE_OF_CONDUCT.md)

In case you are not interested in improving Dog but on building your own tool on top of the Dogfile Spec, please help us discussing it:

- Dogfile Spec [open discussion](https://github.com/dogtools/dog/issues?q=is%3Aissue+is%3Aopen+label%3A%22dogfile+spec%22)
