# Digging In

A tongue in check word play to define a Domain Specific Language for meta measurements in the DNS (for now)

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
