## Required Initial Calls

#### Deploy Services
##### Alias
- @eve-bot deploy {{ namespace }} in {{ environment }}

##### Additional Parameters
- dryrun (bool)
- services (array)

service is in the format of {{ artifact_in_db }}:{{ optional version }}

##### Examples
- @eve-bot deploy current in int
- @eve-bot deploy previous in int services=infocus-cloud-client
- @eve-bot deploy future in int dryrun=true
- @eve-bot deploy future in qa 
- @eve-bot deploy future in int services=infocus-cloud-client:2020.1.3.4 dryrun=true
- @eve-bot deploy current in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true

#### Migrate Customers
##### Alias
- @eve-bot migrate {{ namespace }} in {{ environment }}

##### Additional Parameters
- dryrun (bool)
- databases (array)

##### Examples
- @eve-bot migrate current in int
- @eve-bot migrate future in int databases=infocus,cloud-support

##### Additional Notes
We won't support executing a specific version of a migration framework for now.