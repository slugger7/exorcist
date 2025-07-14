# Exorcist

This project is a linking or a workflow project for the the exorcist project.

It makes use of git submodules for the [server](https://github.com/slugger7/exorcist-server) and [web](https://github.com/slugger7/exorcist-web).

In VS code you can then manage all three projects git individually but have common things living in this repository like [vscode folder](./.vscode)

## Getting started

- `git clone --recurse-submodules git@github.com:slugger7/exorcist.git`
- `git clone --recurse-submodules https://github.com/slugger7/exorcist.git`

If the project has been cloned without the `--recurse-submodules` option simply run
`git submodule update --init --recursive` to initialise the submodules

