## Required Initial Calls

#### Deploy Services
##### Aliases
- @eve-bot deploy services for {{ namespace }} in {{ environment }}

##### Additional Parameters
- dryrun (bool)
- services (array)

service is in the format of {{ artifact_in_db }}:{{ optional version }}

##### Examples
- @eve-bot deploy services for 2020.1 in int
- @eve-bot deploy services for 2020.2 in int services=infocus-cloud-client
- @eve-bot deploy services for 2020.1 in int dryrun=true
- @eve-bot deploy services for 2020.2 in int services=infocus-cloud-client:2020.1.3.4 dryrun=true
- @eve-bot deploy services for 2020.1 in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true

