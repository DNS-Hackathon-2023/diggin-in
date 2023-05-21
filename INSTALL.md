# How to run this code locally

## Prerequisites

This code-runner proof of concept has the following dependencies
* go runtime 1.20
* dig, command line utility to run DNS queries
* jc, to parse dig output into JSON

Make sure to install all dependencies before trying

## Installation

Go in the `code-runner` directory and compile the probe code
```
go build
```
It should generate a `remote-probe` binary

## Running an example
The following example 
```
./remote-probe -out oo.json -probeid mac_probe -script hej.star -wait 1
```

Will run the script `hej.star` that queries one domain and manipulates the key/value storage, every single second. Results are saved in the file `oo.json` in JSN format, one per line.

