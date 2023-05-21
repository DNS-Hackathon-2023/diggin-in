domain = "com."
max_attempts = 10
ref_serial = "1684653963"

def loop():
    # Find all authoritative nameservers for a given domain
    # Call dig function to discover all NS records.
    # TODO: Use the name servers or the addresses?
    auth_ns_query = measure.dig('ns com')

    # Extract the authoritative nameservers
    auth_ns = extract_ns_servers(auth_ns_query)

    # Save auth_ns in the key/value storage for the next run?
    state.set("auth_ns", auth_ns)

    # Retrieve the list of tested_ns
    tested_ns = state.get("tested_ns")

    # Query all nameservers
    for ns in auth_ns:
        # No purpose on testing a nameserver that worked
        if ns not in tested_ns:
            soa_value = measure.dig("SOA com @"+ns)

                if soa_value['serial'] == ref_serial:
                    elapsed = time.time() - t0
                    res.append({'nameserver': ns, 'domain': domain, 'elapsed': elapsed})
                    tested_ns.append(ns)

        if len(auth_ns) == len(tested_ns):
            break

        # Not sure if I can sleep here
        time.sleep(sleep_time)

    return res

# This code comes from "recursive.star", need to check if we can do imports
def extract_ns_servers(dns_results):
    # print(dns_results)

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

# This code comes from serialchecker.star
def extract_serial(dns_results):
    # Check if there is only one DNS result
    if len(dns_results) == 1:
        answer = dns_results[0]['answer']

        # Check if there is only one answer record
        if len(answer) == 1:
            soa = answer[0]

            # Check if the record type is "SOA" (Start of Authority)
            if soa['type'] == "SOA":
                soa_parts = soa['data'].split()

                # Check if the SOA record has the expected number of parts
                if len(soa_parts) == 7:
                    serial = soa_parts[3]
                    return serial

    # Return None if the serial number extraction fails
    return None
