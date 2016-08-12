# Dogfile Spec

*NOTE: This document is a draft and will probably change in the future. Most of its content is still open to discussion.*

[Dog](https://github.com/dogtools/dog) is the first tool that uses Dogfiles and is developed at the same time as the Dogfile Spec itself.

## File Format

Dogfiles are [YAML](http://yaml.org/) files that describe the execution of automated tasks. The root object of a Dogfile is an array of map objects (we call them Tasks). This is an example of a Dogfile with two simple Tasks:

```yml
- task: hello
  description: Say Hello
  run: echo hello

- task: bye
  description: Say Good Bye
  run: echo bye
```

Multiple Dogfiles in the same directory are processed together as a single entity. Although the name `Dogfile.yml` is recommended, any file with a name that starts with `Dogfile` and follows this specification is a valid Dogfile.

## Task definition

The task map accepts the following directives. Please note that directives marked with an asterisk are not implemented in Dog yet and their definition and behaviour will possibly change in the future.

### task

Name of the task. A string that may include lowercase characters (a-z), integers (0-9) and hyphens (-).

```yml
- task: mytask
```

### description

Description of the task. Tasks that avoid this directive are not shown in the task list.

```yml
  description: This task does some cool stuff
```

### run

The code that will be executed.

```yml
  run: echo 'hello'
```

Multiline scripts are supported.

```yml
  run: |
    echo "This is the Dogfile in your current directory:"

    for dogfile in `ls -1 Dogfile*`; do
      cat $dogfile
    done
```

### exec

When this directive is not defined, the default executor is `sh` on UNIX-like operating systems and `cmd` on Windows (not tested yet).

Additional executors are supported if they are present in the system. The following example uses the Ruby executor to print 'Hello World'.

```yml
  task: hello-ruby
  description: Hello World from Ruby
  exec: ruby
  run: |
    hello = "Hello World"
    puts hello
```

The following list of executors are known to work:

- sh
- bash
- python
- ruby
- perl

### pre

Pre-hooks execute other tasks before starting the current one.

```yml
  pre: test
```

Multiple consecutive tasks can be executed as pre-hooks. The tasks defined in the following array will be executed in order, one by one.

```yml
  pre:
    - get-dependencies
    - compile
    - package
```

### post

Post-hooks are analog to pre-hooks but they are executed after current task finishes its execution.

```yml
  post: clean
```

Arrays are also accepted for multi task post-hooks.

```yml
  post:
    - compress
    - upload
```

### workdir

Sets the working directory for the task. Relative paths are considered relative to the location of the Dogfile.

```yml
  workdir: ./app/
```

### tags*

When listing tasks, the ones with the same tag will be shown together. This directive is optional but useful on projects including lots of tasks.

```yml
  tags: dev
```

Multiple tags are allowed.

```yml
  tags:
    - build
    - dev
```

### env

Default values for environment variables can be provided in the Dogfile. They can be modified at execution time.

```yml
  env: ANIMAL=Dog
```

Arrays are also supported.

```yml
  env:
   - CITY=Barcelona
   - ANIMAL=Dog
```

### params*

Additional parameters can be provided to the task that will be executed. All parameters are required at runtime.

```yml
- task: who-am-i
  description: Print my location and who am I
  params:
    # Required parameter without default value
    - name: city

    # Parameter with default value
    - name: planet
      default: Earth

    # Parameter with an array of allowed choices
    - name: animal
      choices:
        - dog
        - cat
        - human

    # Parameter with regular expression validation
    - name: age
      regex: ^\d+$

  run: echo "Hello, I'm in the city of $1, planet $2. I am a $3 and I'm $4 years old"
```

The *regex* option and the *choices* option are mutually exclusive.

### register*

Registers store the STDOUT of executed tasks as environment variables so other tasks can get their value later if they are part of the same task-chain execution.

```yml
  task: get-dog-version
  run: dog --version | awk '{print $3}'
  register: DOG_VERSION

  task: print-dog-version
  description: Print Dog version
  pre: get-dog-version
  run: echo "I am running Dog $DOG_VERSION"
```

Dogfiles don't have global variables, use registers instead.

### Non standard directives*

Tools using Dogfiles and having special requirements can define their own directives. The only requirement for a non standard directive is that its name starts with `x_`. These directives are optional and can be safely ignored by other tools.

```yml
- task: clear-cache
  description: Clear the cache
  x_path: /task/clear-cache
  x_tls_required: true
  run: ./scripts/cache-clear.sh
```

(*) Not implemented yet
