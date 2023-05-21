def loop():
    zone="com"
    #zone="tz.dotat.at"
    dns_results = measure.dig('NS '+ zone +' @8.8.8.8')
    auth_servers = extract_ns_servers(dns_results)
    # Loop through each com server
    for auth_server in auth_servers:
        # Retrieve DNS results using the 'dig' command with the specified com server
        dns_results = measure.dig('soa '+ zone +' @' + auth_server)
        
        # Extract the serial number from the DNS results
        serial = extract_serial(dns_results)
        
        # Retrieve the previously stored serial number for the current com server
        found_serial = state.get(auth_server)
        
        # Compare the current serial number with the previously stored serial number
        if found_serial and found_serial != serial:
            # If the serial numbers don't match, print the new serial number and the com server
            result = {
                'event': "Serial-Change",
                'new_serial': serial,
                'old_serial': found_serial,
                'server': auth_server
            }
            collect(result)
        else:
            # If the serial numbers match or no previous serial number exists, update the stored serial number
            state.set(auth_server, serial)
              
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
    # Check if there is only one DNS result
    if len(dns_results) == 1:
        authorities = None
        if "answer" in dns_results[0]:
            authorities = dns_results[0]['answer']
        elif "authority" in dns_results[0]:
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
