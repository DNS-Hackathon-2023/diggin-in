#!/usr/bin/env python3
import sys
import base64
import json
import datetime
import subprocess
from dnslib import DNSRecord

valid_q_types = {'A': 1, 'NS': 2, 'CNAME': 5, 'SOA': 6, 'NULL': 10, 'PTR': 12, 'HINFO': 13, 'MX': 15, 'TXT': 16, 'RP': 17, 'AFSDB': 18, 'SIG': 24, 'KEY': 25, 'AAAA': 28, 'LOC': 29, 'SRV': 33, 'NAPTR': 35, 'KX': 36, 'CERT': 37, 'A6': 38, 'DNAME': 39, 'OPT': 41, 'APL': 42, 'DS': 43, 'SSHFP': 44, 'IPSECKEY': 45, 'RRSIG': 46, 'NSEC': 47, 'DNSKEY': 48, 'DHCID': 49, 'NSEC3': 50, 'NSEC3PARAM': 51, 'TLSA': 52, 'HIP': 55, 'CDS': 59, 'CDNSKEY': 60, 'OPENPGPKEY': 61, 'CSYNC': 62, 'ZONEMD': 63, 'SVCB': 64, 'HTTPS': 65, 'SPF': 99, 'EUI48': 108, 'EUI64': 109, 'TKEY': 249, 'TSIG': 250, 'IXFR': 251, 'AXFR': 252, 'ANY': 255, 'URI': 256, 'CAA': 257, 'TA': 32768, 'DLV': 32769}
valid_q_classes = {'IN': 1, 'CS': 2, 'CH': 3, 'Hesiod': 4, 'None': 254, '*': 255}

def parse_cli_arguments():
    args = sys.argv[1:]  # Skip the script name

    domain = None
    q_type = None
    q_class = None
    local_server = None

    for arg in args:
        if arg.startswith('@'):
            local_server = arg[1:]
        elif arg.upper() in valid_q_types.keys():
            q_type = valid_q_types[arg.upper()]
        elif arg.upper() in valid_q_classes.keys():
            q_class = valid_q_classes[arg.upper()]
        else:
            domain = arg
    if domain and not domain.endswith('.'):
        domain += '.'
    elif domain is None:
        domain = '.'
    if q_type is None:
        q_type = 1 # Default A record
    if q_class is None:
        q_class = 1 # Default IN class
    if local_server is None:
        local_server = "9.9.9.9" # yeah I know this is bad, but it's an hackaton
    return domain, q_type, q_class, local_server

def run_command(command):
    process = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, error = process.communicate()
    return output.decode().strip()

domain, q_type, q_class, local_server = parse_cli_arguments()
evtdig_command = ("evtdig --query {} --type {} --class {} {}".format(domain,q_type, q_class, local_server))

evtdig_output = run_command(evtdig_command)


# JSON string
json_string = evtdig_output

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
