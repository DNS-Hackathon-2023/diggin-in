tzat_servers = [
	'shale.dotat.at.',
	'ns1.mythic-beasts.com.',
	'ns2.mythic-beasts.com.'
]

def loop():
    # Loop through each tzat server
    for tzat_server in tzat_servers:
        # Retrieve DNS results using the 'dig' command with the specified com server
        dns_results = measure.dig('soa tz.dotat.at @' + tzat_server)
        
        # Extract the serial number from the DNS results
        serial = extract_serial(dns_results)
        
        # Retrieve the previously stored serial number for the current tzat server
        found_serial = state.get(tzat_server)
        
        # Compare the current serial number with the previously stored serial number
        if found_serial and found_serial != serial:
            # If the serial numbers don't match, print the new serial number and the tzat server
            result={}
            result['event']="Serial-Change"
            result['new_serial']=serial
            result['old_serial']=found_serial
            result['server']=tzat_server
            collect(result)
        else:
            # If the serial numbers match or no previous serial number exists, update the stored serial number
            state.set(tzat_server, serial)
              
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
