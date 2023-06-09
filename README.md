# DIG-alicious

A tongue in cheek word play to define a Domain Specific Language for meta measurements in the DNS (for now)

## prototype not ready for production

Nothing valued is here. What is here was dangerous and repulsive to us. This message is a warning about danger.

## Introduction

We all have worked with RIPE Atlas and other measurements platforms, open and semi-closed, and ended up doing the same boilerplate of setting up core measurements on a regular basis with a set of probes, until the monitoring period ends or the expected situation materializes.

We'd like to define a language that will describe how to do some high level measurements, and implement highly-efficient providers of those measurements in existing measurement networks

## Use cases

### Monitor the propagation of one zone across all their authoritative nameservers

### Monitor serial number changes

### atlas-shruggd (example: do dns resolution from a different root-server system)

### dnsthought 

does a lot of dns server capabilities measurements. For instance measure DNSSEC algorithm compatibility by comparing a broken zone to a correct one for a given DNSSEC algorithm

### dns minimisation research

query an instrumented server from a set of resolvers to figure out if it does qname min (via TXT fetch on an instrumented server). Then do you further experiments only for the set of resolvers that do qname min

### find all possible resolution paths for a domain

Raffaele use-case: useful for domains under attack

### Define an experiment to define where an authorative nameserver operator might need to place a new anycast node

From Vicky Risk

See also existing research by SIDN https://github.com/SIDN/BGP-Anycast-Tuner

### Distributed monitoring for specific domain record

From Vicky Risk, to quickly identify if a nameserver address has been changed or not


# Scope
* you need to measure from many vantage points
* DNS measurement primitives are available
* results storage is available
* needs to run in a constrained environment
* security is a major concern (measurement nodes typically are "somebody elses computer")
* NOT BE A BOTNET

# Users
Small set of people who will write this type of scripts
Much bigger set of people who will benefit from the availability of these types of scripts for DNS measurements

# Solution
A domain specific language that can run on the measurement platform to glue all of this together

# What we tried

## Starlark

We found the Starlark language which is a subset of Python. It addresses a few things:

* sandboxed: we can restrict the library functions so scripts written in Starlark can only do what we permit

* like Python, so we hope it is more accessible to people writing scripts than alternatives such as Tcl, Lua, Guile, ...

* Starlark implementations exist in Java and Go (from Google), and Rust (from Facebook); it was designed for the Bazel and Buck build systems


# By-catch
You could use something similar to this but without the measurement code, ie. just define the data reduction step

# Caveats
The DNS thought use cases don't work because their measurement server returns an invalid LCLS Extended option that dig 9.18 generates an error that jc can't parse

```
dig A secure.d2a3n1.rootcanary.net | jc --dig
jc:  Error - dig parser could not parse the input data.
             If this is the correct parser, try setting the locale to C
                 (LC_ALL=C).
             For details use the -d or -dd option. Use "jc -h --dig" for help.
```
