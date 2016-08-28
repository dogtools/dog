# Dog

[![Build Status](https://travis-ci.org/dogtools/dog.svg?branch=master)](https://travis-ci.org/dogtools/dog)
[![Join the chat](https://badges.gitter.im/dogtools/dog.svg)](https://gitter.im/dogtools/dog)

Dog is a command line application that executes automated tasks. It works in a similar way as GNU Make but it is a more generic task runner, not a build tool. Dog's default script syntax is `sh` but most interpreted languages like BASH, Python, Ruby or Perl can also be used.

## Installing Dog

If you are on macOS you can install Dog using brew:

```
brew tap dogtools/dog
brew install dog
```

If you have your golang environment set up, you can use:

```
go get github.com/dogtools/dog
```

## Using Dog

List all tasks in current project

    dog

Execute a task

    dog taskname

Execute a task, printing elapsed time and status code

    dog -i taskname

## What is a Dogfile?

Dogfile is a specification that uses YAML to describe the tasks related to a project. We think that the Spec will be finished (no further breaking changes) by the v1.0.0 version of Dog.

- Read Dog's own [Dogfile.yml][1]
- Read the [Dogfile Spec][2]

## Other tools

Tools that use Dogfiles are called *dogtools*. Dog is the first dogtool but there are other things that can implemented in the future: web and desktop UIs, chat bot interfaces, plugins for text editors and IDEs, tools to export Dogfiles to other formats, HTTP API interfaces, even implementations of the cli in other languages! To simplify the process of creating dogtools we are implementing parts of Dog as Go packages so they can be used in other projects (see [parser][3], [types][4] and [execute][5]). Let us know if you have any uncovered need on any of these packages.

## Contributing

If you want to help, take a look at:

- Open [bugs][6]
- Lacking features for [v0.4.0][7]
- Our [Code of Conduct][8]

In case you are not interested in improving Dog but on building your own tool on top of the Dogfile Spec, please help us discussing it:

- Dogfile Spec [open discussion][9]

[1]: https://github.com/dogtools/dog/blob/master/Dogfile.yml
[2]: https://github.com/dogtools/dog/blob/master/DOGFILE_SPEC.md
[3]: https://github.com/dogtools/dog/tree/master/parser
[4]: https://github.com/dogtools/dog/tree/master/types
[5]: https://github.com/dogtools/dog/tree/master/execute
[6]: https://github.com/dogtools/dog/issues?q=is%3Aissue+is%3Aopen+label%3Abug
[7]: https://github.com/dogtools/dog/milestone/4
[8]: https://github.com/dogtools/dog/blob/master/CODE_OF_CONDUCT.md
[9]: https://github.com/dogtools/dog/issues?q=is%3Aissue+is%3Aopen+label%3A%22dogfile+spec%22
