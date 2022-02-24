# Khan

An interactive CLI management tool for [Hashicorp's Nomad](https://www.nomadproject.io/)

<img width="983" alt="image" src="https://user-images.githubusercontent.com/5923958/155530556-35ccee7a-67eb-4788-845b-6d8f2f871a7b.png">

## Why

Nomad has a great CLI tool with a lot of features, but copy/pasting IDs around
and wanting to see steady updates among multiple objects can be tricky.  There's
also Nomad's full web UI, but this has overhead in setup and can be a pain to
work with through SSH tunnels.

There is room for a middle ground for quick debugging/troubleshooting on the
command line, similar to a tool like [k9s](https://github.com/derailed/k9s).

Enter Khan!

## Installation

Head to the [releases page](https://github.com/Evertras/khan/releases) and
download the binary directly for your OS.

## Configuration

Khan uses the [default Nomad configuration variables](https://www.nomadproject.io/docs/commands#connection-environment-variables).

## Developer Requirements

The following are required as global installs for development:

* Python 3 (standard install for Linux/Mac should be fine)
* Go 1.17+ ([goenv](https://github.com/syndbg/goenv) is recommended)
* Make (probably already installed)

Other tools will be handled automatically by the Makefile.  These tools will be
downloaded locally to the `./bin` folder.  For ease of use, you may want to use
[direnv](https://direnv.net/) with the supplied [.envrc.example](.envrc.example)
file to add this local path to your bin, so you can run these tools as if they
were globally installed.

## Running a Nomad test server

Any Nomad server can be used for testing, but for simplicity and self-contained
development a quick development server can be brought up with
`make nomad-test-server` in another terminal.  This will start a Nomad agent in
dev mode which brings up the server and a single client, itself.

More complicated setups will be added later via Vagrant, but this is enough for
a starting point.
