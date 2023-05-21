## Ripe Atlas

Ripe Atlas uses evtdig as dns query software.

We wrote two wrapper, one (parser.py) converts the dig CLI into evtdig CLI.

The second wrapper (converter.py) converts evtdig output to dig output.

### Caveats and Bugs

- evtdig does not support explicit qtype and qname, it requires them in binary (decimal format). 
This is currently handled by the parser.py script.
- evtdig does not parse AUTHORITY section. Hence, this is not handled by the converster script.
