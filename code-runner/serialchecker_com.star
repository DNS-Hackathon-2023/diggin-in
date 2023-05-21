com_servers = [
    '192.12.94.30',
    '192.33.14.30',
    '192.48.79.30',
    '192.55.83.30',
    '192.43.172.30',
    '192.35.51.30',
    '192.5.6.30',
    '192.42.93.30',
    '192.54.112.30',
    '192.41.162.30',
    '192.52.178.30',
    '192.26.92.30',
    '192.31.80.30'
]

def loop():
    # Loop through each com server
    for com_server in com_servers:
        # Retrieve DNS results using the 'dig' command with the specified com server
        dns_results = measure.dig('soa com @' + com_server)
        
        # Extract the serial number from the DNS results
        serial = extract_serial(dns_results)
        
        # Retrieve the previously stored serial number for the current com server
        found_serial = state.get(com_server)
        
        # Compare the current serial number with the previously stored serial number
        if found_serial and found_serial != serial:
            # If the serial numbers don't match, print the new serial number and the com server
            result = {
                'event': "Serial-Change",
                'new_serial': serial,
                'old_serial': found_serial,
                'server': com_server
            }
            collect(result)
        else:
            # If the serial numbers match or no previous serial number exists, update the stored serial number
            state.set(com_server, serial)
              
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
