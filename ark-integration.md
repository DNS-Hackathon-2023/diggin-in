### Ark Integration
Ark uses scamper as dns query software.

Due to limited time constrains we did not implement the scamper wrapper.

To do so we will need to:
- Run scamper daemon ```scamper -P <port>```

Attach to the TCP port on the local machine and send the following commands:
```
attach format json
host -s nameserver -r domain
```

It will output something like:
```json
{
  "type":"host", "version":"0.1", "src":"172.23.130.14", "dst":"8.8.8.8",
  "userid":0, "start":{"sec":1684599535,"usec":392069}, "wait":5000,
  "retries":0, "stop":"DONE", "qname":"google.com", 
  "qclass":"IN", "qtype":"A", "qcount":1, 
  "queries":[{"id":8, "ancount":1, "nscount":0, "arcount":0,
      "tx":{"sec":1684599535,"usec":392108}, "rx":{"sec":1684599535,"usec":441987}, 
      "an":[{"class":"IN", "type":"A", "ttl":300, "name":"google.com", "address":"142.250.179.206"}]}
      ]
}
```

This should be converted to ```jc --dig``` like format.
