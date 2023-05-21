# Security considerations

  - CPU usage
  - memory usage
  - test execution time
  - scheduling timeframe
  - number of simultaneous tests running
  - malicious code injection
  - process isolation
  - user access
  - cache poisoning of resolver on host network (from inside)
  - internal network DNS configuration exposure

  - global restrictions? (orchestration)
  - list of known hosts? (white-/black-listing)

# Mitigations

Restrictions set on host should be calculated (at a healthy a safety marigin) with regards to the host capabilities. 

- CPU usage and memory usage should be restricted on individual tests, as well as on the total resources consumed.
- execution time and scheduling timeframe should be restricted on individual tests to prevent congestion.
- number of simultaneous tests during any given timeframe shuld be restricted with the maximum consumed resources per test in mind.

- host environment and application should be hardened and tests isolated to prevent cross contamination of data, as well as execution of code outside the sandbox.
- access to host should be restricted.
- host network admin should take apropriate measures to prevent cache poisoning of the internal ressolver from maliciously constructed queries and exfiltration of internal DNS network configuration (using for example PTR queries)

- global restrictions to prevent DDOS is not in scope at this time, but may be relevant if fleet based test suites are deployed. This requires some sort of orchestration.
- a public list of known hosts may be either helpful or harmful, as it can be used by network admins to generate white-/black-lists as well as by malicious parties to find vulnerable implementations. Also out of scope at this time.
