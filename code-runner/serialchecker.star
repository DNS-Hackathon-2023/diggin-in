root_servers = [
    "198.41.0.4",
    "199.9.14.201",
    "192.33.4.12",
    "199.7.91.13",
    "192.203.230.10",
    "192.5.5.241",
    "192.112.36.4",
    "198.97.190.53",
    "192.36.148.17",
    "192.58.128.30",
    "193.0.14.129",
    "199.7.83.42",
    "202.12.27.33"
]

load("primitives.star", "extract_serial")

def loop():
    # Loop through each root server
    for root_server in root_servers:
        # Retrieve DNS results using the 'dig' command with the specified root server
        dns_results = measure.dig('soa . @' + root_server)
        
        # Extract the serial number from the DNS results
        serial = extract_serial(dns_results)
        
        # Retrieve the previously stored serial number for the current root server
        found_serial = state.get(root_server)
        
        # Compare the current serial number with the previously stored serial number
        if found_serial and found_serial != serial:
            # If the serial numbers don't match, print the new serial number and the root server
            print(serial + " @ " + root_server)
        else:
            # If the serial numbers match or no previous serial number exists, update the stored serial number
            state.set(root_server, serial)
              
