# uhandles

An escalation of [jsonsubst](https://github.com/arbitar/jsonsubst).  
Small static Go program with no external runtime dependencies that takes a [handlebars](https://handlebarsjs.com) template in from stdin, renders it with data from a JSON file (passed by argument), and outputs the result to stdout.  
This is basically just a simple wrapper around [aymerick/raymond](https://github.com/aymerick/raymond).

## Usage
Very simple. One argument: JSON file.

`./uhandles data.json <input.txt >output.txt`

Or, in a pipeline of some kind...

`cat input.txt | ./uhandles data.json | tee output.txt`

## Build
If you're cool and already have a good working Go build environment on your system, just run `make build`.

I'm not cool though, so I use Docker for this kind of thing, which is why it's the default make target. Just run `make` on a system with Docker running on it, and a build image will be created, `make build` executed within the container where all the prerequisites exist, and the build image will be removed after it's done.

## Notes
No external runtime dependencies, so it should run anywhere. Which means it'll run inside your quick-and-dirty Docker stuff without complaint.
