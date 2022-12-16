# uhandles
("micro-handles": for small, portable handlebars)

An escalation of [jsonsubst](https://github.com/arbitar/jsonsubst).  
Small static Go program with no external runtime dependencies that takes [handlebars](https://handlebarsjs.com) template(s), renders it with reference to JSON data, and outputs the results as directed.  
This is basically just a wrapper around [aymerick/raymond](https://github.com/aymerick/raymond).

## Basic Usage
In its simplest form, one argument:

`./uhandles data.json <input.txt >output.txt`

Or, in a pipeline of some kind...

`cat input.txt | ./uhandles data.json | tee output.txt`

**Example:**  
*data.json*:
```
{
	"object": {
		"property": "value",
		"list": ["one", "two", "three"]
	},
	"objects": [
		{ "lower": "a", "upper": "A" },
		{ "lower": "b", "upper": "B" }
	]
}
```
*template.tmpl.txt*:
```
Hi! This is object.property's value: {{ object.property }}

Here's what's in an array:
{{#each object.list}}
- {{this}}
{{/each}}

Here's what's in an array of objects:
{{#each objects}}
{{#with this}}
- The upper-case version of "{{lower}}" is "{{upper}}"
{{/with}}
{{/each}}
```

After running:  
`./uhandles data.json <template.tmpl.txt`

The expected output on stdout would be:
```
Hi! This is object.property's value: value

Here's what's in an array:
- one
- two
- three

Here's what's in an array of objects:
- The upper-case version of "a" is "A"
- The upper-case version of "b" is "B"
```

For more information on the template format, please see the official [Handlebars documentation](https://handlebarsjs.com/guide/).

For more information on the flavor of handlebars being used here, please see the project for which this is essentially just a wrapper for, [aymerick/raymond](https://github.com/aymerick/raymond).

## Advanced Usage
If you don't want to the template coming from stdin, or the output to stdout, or if you want the template from a file and the JSON from stdin, or any other number of weird things... there is a more complicated flag-based usage:

`./uhandles --data data.json --template input.txt --output output.txt`

You can use `-` for any of these file paths to specify stdin/stdout as appropriate. It will complain if you specify `-` for both data and template, so don't get too clever.

| Flag                            | Description |
| ------------------------------- | - |
| `-d\|--data <file\|->`          | JSON file path to use to render the template, or "-" for stdin |
| `-t\|--template <file\|dir\|->` | Input Handlebars template path, or directory, or "-" for stdin |
| `-o\|--output <file\|dir\|->`   | Output file path, or directory, or "-" for stdout |
| `--tmpl-token <token>`          | Token to use to detect template files. Default: `".tmpl"`. Is assembled into a glob like `*<token>*`, and output files have instances of `<token>` removed.

## Environment Variables
uhandles will look for the following environment variables and use them in lieu of their associated flags:

| Env                   | Flag               |
| --------------------- | ------------------ |
| `UHANDLES_DATA`       | `-d`, `--data`     |
| `UHANDLES_TEMPLATE`   | `-t`, `--template` |
| `UHANDLES_OUTPUT`     | `-o`, `--output`   |
| `UHANDLES_TMPL_TOKEN` | `--tmpl-token`     |

## Directory Of Templates Rendering
When a directory is specified for `--template`, it looks within for any `*.tmpl*` files within it. Each of these files is interpreted as a template and separately rendered with the data from the JSON file. If `--output` specifies a directory, a corresponding file to each template (with the same permissions and the tmpl-token removed from the filename) will be created in the output directory. Otherwise, if `--output` specifies a single file or stdout, concatenation will be performed as per the next section.

**Example:**  
Input directory:  
```
./input/
  input_1.tmpl.txt
	input_2.txt
	input_3.txt.tmpl
```

With an output directory empty, but present, at `./output/`, then run with:  
`./uhandles -d data.json -t ./input/ -o ./output/`

This results in an output directory containing the rendered templates:  
```
./output/
  input_1.txt
  input_3.txt
```

Note that `input_2.txt` is ignored, because it does not contain the template token.

If the output directory already contains files, they will be left in place, unless they are to be written to by the current operation. If so, they are overwritten entirely and their mode is changed to match the input files.

## Concatenation
When there are multiple template files but only one output file (a file path or stdout), the multiple template files are rendered in directory listing order and concatenated together with single newlines between them into the single output file.

**Example:**  
Input directory:
```
./input/
  input_1.tmpl.txt
  input_2.txt
  input_3.txt.tmpl
```

Run with:  
`./uhandles -d data.json -t ./input/ -o ./output.txt`

Results in a single file, output.txt, containing:
```
<rendered outputs of input_1.tmpl.txt>
<rendered outputs of input_3.txt.tmpl>
```

Note the newline between the two. Also, `input_2.txt` is ignored, because it does not contain the template token. Output could also have been omitted or `-` to direct the output to stdout.

## Build
If you're cool and already have a good working Go build environment on your system, just run `make build`.

I'm not cool though, and I'm pretty lazy, so I use Docker for this kind of thing, which is why it's the default make target. Just run `make` on a system with Docker running on it, and a build image will be created, `make build` executed within the container where all the prerequisites exist, and the build image will be removed after it's done.

If you're iterating, you can `make docker-prep` to create the image, `make docker-build` repeatedly to run the build, and `make docker-clean` to remove the image.

## Notes
No external runtime dependencies, so it should run anywhere. Which means it'll run inside your quick-and-dirty Docker stuff without complaint.
