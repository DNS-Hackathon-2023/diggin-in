
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
        result=recursion(root_servers)
        print(result)

def recursion(ns_servers):
        for server in ns_servers:
           dns_results = measure.dig('A test.disjoint.superdns.nl @' + server)
           if(len(dns_results)==1):
                  if("aa" not in dns_results[0]['flags']):
                     ns_servers=extract_ns_servers(dns_results)
                     if(len(ns_servers)>0):
                        return recursion(ns_servers)
                  if("aa" in dns_results[0]['flags']):
                     return dns_results


def extract_ns_servers(dns_results):
    print(dns_results)
    if len(dns_results) == 1:
        authorities = dns_results[0]['authority']
        ns_servers=[]
        for authority in authorities:
            if authority['type'] == "NS":
                ns_records = authority['data']
                print(ns_records)
                ns_servers.append(ns_records)
        return ns_servers      
    return None
