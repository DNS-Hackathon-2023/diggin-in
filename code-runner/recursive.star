root_servers = [
    "a.root-servers.net",
    "b.root-servers.net",
    "c.root-servers.net",
    "d.root-servers.net",
    "e.root-servers.net",
    "f.root-servers.net",
    "g.root-servers.net",
    "h.root-servers.net",
    "j.root-servers.net",
    "k.root-servers.net",
    "l.root-servers.net",
    "m.root-servers.net"
]

def loop():
    # Initialize ns_servers with root_servers
    ns_servers = root_servers
    
    # Bound the recursive resolution to 10 iteration
    for _ in range(10):
        # Query DNS for A record of "google.com" using the first nameserver in ns_servers (starting from root)
        dns_results = measure.dig('A google.com @' + ns_servers[0])
        
        # Check if there is only one DNS result (single dig output)
        if len(dns_results) == 1:
            # Check if answer is not authoriative
            if "aa" not in dns_results[0]['flags']:
                # Extract the nameservers referral from the authoriative section
                ns_servers = extract_ns_servers(dns_results)
                
                # Continue to the next iteration if nameservers are found
                if len(ns_servers) > 0:
                    continue
            
            # If answer is authoriative provide it back
            if "aa" in dns_results[0]['flags']:
                print(dns_results)  # Print the DNS result if the flag is present
    
    return None  # Return None if the loop completes without finding a result


def extract_ns_servers(dns_results):
    print(dns_results)
    
    # Check if there is only one DNS result
    if len(dns_results) == 1:
        authorities = dns_results[0]['authority']
        ns_servers = []
        
        # Iterate over the authority records
        for authority in authorities:
            # Check if the record type is "NS" (Nameserver)
            if authority['type'] == "NS":
                ns_records = authority['data']
                ns_servers.append(ns_records)
        
        return ns_servers  # Return the extracted nameservers
    
    return None  # Return None if the extraction fails

