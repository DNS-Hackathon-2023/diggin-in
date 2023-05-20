def propagation(domain, ref_serial, max_attempts=10, sleep_time=15):
    # Assuming "domain" is the right name
    # Return a list of dictionaries where the key is the nameserver and the value is the time elapsed between the start of probing until the nameserver returns a matching serial

    # Find all authoritative nameservers for a given domain
    # Call dig function to discover all NS records.
    # TODO: Use the name servers or the addresses?
    auth_ns = measure.dig(domain)

    # XXX: Assuming this is available
    t0 = time.time()

    # Regularly find the serial, capture how long did it take from the initial
    res = []
    tested_ns = []
    for i in range(0, max_attempts):
        # Query all nameservers
        for ns in auth_ns:
            # No purpose on testing a nameserver that worked
            if ns not in tested_ns:
                soa_value = measure.dig(domain, ns, 'SOA')

                if soa_value['serial'] == ref_serial:
                    elapsed = time.time() - t0
                    res.append({'nameserver': ns, 'domain': domain, 'elapsed': elapsed})
                    tested_ns.append(ns)

        if len(auth_ns) == len(tested_ns):
            break

        # Not sure if I can sleep here
        time.sleep(sleep_time)

    return res
      
