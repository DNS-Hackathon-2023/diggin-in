import base64
import json
import datetime
from dnslib import DNSRecord

# JSON string
json_string = input()

# Remove the "RESULT " prefix
json_string = json_string.replace('RESULT ', '')

# Parse the JSON string
data = json.loads(json_string)

# Extract the value of the "abuf" field
abuf = data['result']['abuf']
datetime_obj = datetime.datetime.fromtimestamp(data['time'])

# DNS answer buffer
dns_answer_buffer = abuf

# Decoding the DNS answer buffer
decoded_buffer = base64.b64decode(dns_answer_buffer)

# Parsing the DNS message using dnslib
dns_record = DNSRecord.parse(decoded_buffer)

formatted_datetime = datetime_obj.strftime(";; WHEN: %a %b %d %H:%M:%S %Z %Y")


print("""; <<>> eVDiG ConVerter <<>> dummy.com
;; global options: +cmd
;; Got answer:""")
print(dns_record)
print(""";; Query time: {} msec
;; SERVER: {}#{}({}) ({})""".format(int(data['result']['rt']),data['dst_addr'], data['dst_port'], data['dst_addr'], data['proto']))
print(formatted_datetime)
print(""";; MSG SIZE  rcvd: {}
""".format(data['result']['size']))
