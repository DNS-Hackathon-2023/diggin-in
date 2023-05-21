## Ripe Atlas Integration

Ripe Atlas uses evtdig as dns query software.

We wrote two wrapper, one (parser.py) converts the dig CLI into evtdig CLI.

The second wrapper (converter.py) converts evtdig output to dig output.
### Example of CLI

```console
test@test:~/RIPE/atlasswprobe-5080-1/usr/local/atlas/bb-13.3/bin$ 
python3 parser.py NS IN google.com @8.8.8.8 | xargs ./evtdig -R | python3 converter.py | jc --dig
```
The args:
```
NS IN google.com @8.8.8.8
```
follows the same format of dig and can be placed in random order or omitted.
### Example of the output parsed by ```jc --dig```
```json
[
  {
    "id": 45298,
    "opcode": "QUERY",
    "status": "NOERROR",
    "flags": [
      "qr",
      "rd",
      "ra"
    ],
    "query_num": 1,
    "answer_num": 4,
    "authority_num": 0,
    "additional_num": 0,
    "question": {
      "name": "google.com.",
      "class": "IN",
      "type": "NS"
    },
    "answer": [
      {
        "name": "google.com.",
        "class": "IN",
        "type": "NS",
        "ttl": 20888,
        "data": "ns1.google.com."
      },
      {
        "name": "google.com.",
        "class": "IN",
        "type": "NS",
        "ttl": 20888,
        "data": "ns3.google.com."
      },
      {
        "name": "google.com.",
        "class": "IN",
        "type": "NS",
        "ttl": 20888,
        "data": "ns4.google.com."
      },
      {
        "name": "google.com.",
        "class": "IN",
        "type": "NS",
        "ttl": 20888,
        "data": "ns2.google.com."
      }
    ],
    "query_time": 11,
    "server": "8.8.8.8#53(8.8.8.8) (UDP)",
    "when": "Sun May 21 10:05:55  2023",
    "rcvd": 100
  }
]
```


### Caveats and Bugs

- evtdig does not support explicit qtype and qname, it requires them in binary (decimal format). 
This is currently handled by the parser.py script.
- evtdig does not parse AUTHORITY section. Hence, this is not handled by the converster script.
