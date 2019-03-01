# qri-datapackage

Command Line tool for integrating Qri with the [Frictionless Data Datapackage format](https://frictionlessdata.io/data-packages/). Currently a work-in-progress proof-of concept type thing.

```
$ qri-datapackage
open knowledge foundation datapackage qri integration

Usage:
   [command]

Available Commands:
  export      write a qri dataset as a datapackage zip archive
  help        Help about any command
  import      import a datapackage into qri

Flags:
      --debug   show debug output
  -h, --help    help for this command

Use " [command] --help" for more information about a command.
```

# Installation
This is a bit of an experiment, so for now we're hoping you're ok with building from source. If others want it we can distribute binaries:
1. You'll need go installed: https://golang.org, installation should include instructions for getting `$GOPATH/bin` on your `$PATH`, be sure to follow those for easy next steps.
2. You'll need the Qri command line client, instructions for setup are here: https://github.com/qri-io/qri#building-from-source. If can type `qri` in a terminal & see instructions, you're good to go.
3. All that's left is to get this package:
```
$ go get -u github.com/qri-io/qri-datapackage
```

Which will place a new binary called `qri-datapackage` at `$GOPATH/bin`, which (ideally) means you can just type `qri-datapackage`, and it'll show a help message.

If you have _any issue at all_, hit up our discord chat: https://discord.gg/etap8Gb, someone should be able to help.

# Usage

For now the flow is pretty limited, import from datapackage, repeat as necessary to create new versions, then export.
```
$ cd path/to/folder/with_datapackage

$ qri-datapackage import datapackage.json
dataset saved: b5/human_services_data@QmSyDX5LYTiwQi861F5NAwdHrrnd1iRGsoEvCyzQMUyZ4W/ipfs/Qmda9B5BQesk8gXMPJnubGThgkny9AiEB7CLHSwkXHY4LA

# prove nothing changed by repeating:
$ qri-datapackage import datapackage.json
error saving: no changes detected
qri error: exit status 1

$ qri-datapackage export me/human_services_data
exported datapackage .zip archive to: human_services_data_datapackage.zip
```