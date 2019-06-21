# auto-docs

## background

This project aims to be somewhat of a Knowledge Management system
within an automated software development ecosystem.

## configuration

By default, configuration is loaded from ``~/.ad.yaml``, but this file
can be specified using the ``--config`` flag when running the command.

```yaml
# Listen defines the listener string that the server will respond to.
# Generally, this will be a colon followed by the desired port.
Listen: ":8008"

# Name is the pretty name to use through the application. (Not used).
Name: "AutoDocs"

# Git is an object that defines attributes used for interacting
# with git.
Git:

  # URI is the remote location of the repository.
  URI: "git@github.com:CerealBoy/auto-docs.git"

  # Branch denotes the remote branch to use.
  Branch: "master"

  # LocalPath is a location on-disk for auto-docs to manage the repo.
  LocalPath: "/tmp/auto-docs-git"

  # SSHKey is the path to the private key to use for git auth.
  SSHKey: "/home/ad/.ssh/id_rsa"

  # Timeout is the maximum time to wait (in ms) for a pull request to
  # complete. (Not used).
  Timeout: 5000

  # Username is the user to connect to git with.
  Username: "git"

  # Password is used instead of SSH if specified. (Not used).
  Password: ""

  # Period is the length (in s) between pull requests.
  Period: 10
```

## building

``go-bindata`` is required to load binary data into the Go context,
allowing for files to be deflated at runtime. To install this tool,
a make target is provided.

    make bin-prep

To compile the frontend components ready for binary compilation,
another make target is available.

    make install

To build the bindata requirements, separate make targets are
utilised.

    make bin-dist

Once all bindata requirements are available, a series of binaries
can be constructed from 1 makefile command.

    make binaries

A series of binaries will then be available within the
``build/`` directory.

To build the complete set of required binaries, and prepare all
the necessary dependent components, a target is also available.

    make complete-binaries

## development

Running tests will require the frontend to be built, and bindata
to be made available. These are captured in appropriate make
targets, and summed up into a single test.

    make test

For streamlined development, the frontend can be setup to make use
of a watch task.

    make dev-fe

A base target is available to compile the Go binary and execute it.

    make dev-be

