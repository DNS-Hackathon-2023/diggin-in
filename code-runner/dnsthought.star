# This is intended to be a library of DNS Thought data gathering
# As documented in https://dnsthought.nlnetlabs.nl/raw/

def loop():
    dns_results = measure.dig('A secure.d2a3n1.rootcanary.net')
    result = {
        'event': "DNS Thought",
        'result': dns_results
    }
    collect(result)

def dnskey_rsa_md5_ds_sha256():
    # Send a query to the local resolver for A secure.d2a3n1.rootcanary.net
    dns_results = measure.dig('A secure.d2a3n1.rootcanary.net')

    # Retrieve previous result
    prev_res = state.get('secure.d2a3n1.rootcanary.net')

    # Look into the additional section for an A record, if there, the test failed?
    
def dnskey_algo_ed448():
    # Query the control case
    # Query the test case, sending a specific query to a given server
    # Compare the structure of both responses
    # If the key elements of the responses don't match, return a failed status

def dnskey_algo_ed25519():
    return 0

def dnskey_algo_ecdsa_p384_sha384():
    return 0

def dnskey_algo_ecdsa_p256_sha256():
    return 0

def dnskey_algorithm_ecc_gost():
    return 0
