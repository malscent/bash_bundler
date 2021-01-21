# Simple Bash Bundler

A simple CLI app for bundling multiple bash scripts together via the source command.

Usage:

```
This simple tool bundles bash files into a single bash file.

Usage:
   sbb {flags}
   sbb <command> {flags}

Commands: 
   bundle                        bundle a bash script
   help                          displays usage informationn
   version                       displays version number

Flags: 
   -e, --entry                   The entrypoint to the bash script to bundle. 
   -h, --help                    displays usage information of the application (default: false)
   -o, --output                  The output file to write to 
   ```