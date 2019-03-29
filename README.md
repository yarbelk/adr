# ADR

A tool for managing and writing ADRs for a project.

Warning: this is evolving quickly, its currently is usable in that you can init, create and render
ADRs, but the code is in this 'get it working so I can use ADRs to track my thinking and make it
better' so i can develop it while dogfooding the tool.

Contributions would be nice, eventually, but right now its probably not a good idea because
so much changes and I'm not always in a position to push the code.

Vague Highlevel
---------------

- The command creates the adrs as a machine readable (currently TOML, but will make that
  configurable, though i don't know why anyone would actually care) format.
- This is rendred to markdown, because i'm too lazy to parse markdown to do things like
  update 'Status' and 'Related' fields.  Much easier to manage this as toml and render.
- Rendering is to markdowm, locations are set in config.  config uses viper, and
  I'll document it sometime.
- Its supposed to have a 'global' config, and a local config.
- There will be autocomplete
- It will make you coffee and ask about your day.

right now, it creates things in a DRAFT state, and will let you have
multiple authors, multiple related, and superceeded adrs.  It should update
related and superceeded ones for you.

Rending to markdown is 'manual': `adr render`, or when creating a new ADR, pass
in `-R` to auto render all adrs again.

It will use the `EDITOR` in you environment for new adrs.  Or `VISUAL`.  You have
to save and close the editor for changes to be picked up.

There are status, but i haven't really implemented them and they are
placeholders

TODO for Docs
-------------

actually document the behaviour and ideas.  Some things are placeholder and I haven't
actually thought them through




the `--help` text


```
Write, manage and list your  ADRs for a project.  For example:

initialize a new ADR directory:
        adr init docs/adrs

Add a new ADR:
        adr new Implement a Widget Factory

List the existing ADRs:
        adr list

Usage:
  adr [command]

Available Commands:
  help        Help about any command
  init        initialize a new
  new         New ADR
  render      render all the ADRs to human readable format like markdown

Flags:
      --config string   config file (default is $HOME/.adr.yaml)
      --dir string      where to store the adrs
  -h, --help            help for adr
      --version         version for adr

Use "adr [command] --help" for more information about a command. 
```
