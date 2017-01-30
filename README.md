# Dog

[![Build Status](https://travis-ci.org/dogtools/dog.svg?branch=master)](https://travis-ci.org/dogtools/dog)

Dog is a command line application that executes automated tasks. It works in a similar way as GNU Make but it is a more generic task runner, not a build tool. Dog's default script syntax is `sh` but most interpreted languages like BASH, Python or Ruby can also be used.

## Using Dog

List all tasks in current project

    dog

Execute a task

    dog taskname

Execute a task, printing elapsed time and exit status

    dog -i taskname

## What is a Dogfile?

Dogfile is a specification that uses YAML to describe the tasks related to a project. We think that the Spec will be finished (no further breaking changes) by the v1.0.0 version of Dog.

- Read Dog's own [Dogfile.yml][1]
- Read the [Dogfile Spec][2]

## Installing Dog

If you are using macOS you can install Dog using brew:

    brew tap dogtools/dog
    brew install dog

If you have your golang environment set up, you can use:

    go get -u github.com/dogtools/dog

## Other tools

Tools that use the Dogfile Specification are called *dogtools*. Dog is the first dogtool but there are other things that can be implemented in the future: web and desktop UIs, chat bot interfaces, plugins for text editors and IDEs, tools to export Dogfiles to other formats, HTTP API interfaces, even implementations of the cli in other languages!

The root directory of this repository contains the dog package that can be used to create dogtools in Go.

    import "github.com/dogtools/dog"

Check the `examples/` directory to see how it works.

## Contributing

If you want to help, take a look at the open [bugs][3], the list of all [issues][4] and our [Code of Conduct][5].

[1]: https://github.com/dogtools/dog/blob/master/Dogfile.yml
[2]: https://github.com/dogtools/dog/blob/master/DOGFILE_SPEC.md
[3]: https://github.com/dogtools/dog/issues?q=is%3Aissue+is%3Aopen+label%3Abug
[4]: https://github.com/dogtools/dog/issues
[5]: https://github.com/dogtools/dog/blob/master/CODE_OF_CONDUCT.md
