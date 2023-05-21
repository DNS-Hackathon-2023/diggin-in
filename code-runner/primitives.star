# Put here primitives allowing parsing or handling results

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
                    serial = soa_parts[2]
                    return serial

    # Return None if the serial number extraction fails
    return None

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

