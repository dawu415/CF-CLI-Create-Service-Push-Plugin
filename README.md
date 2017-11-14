#Create-Service-Push

Reads in a services manifest yml file and creates the services listed and pushes 
an application.

This plugin extends cf push. So it will act the same way as cf push, with the exception
that it creates the services first.

#Additional Parameters


#services-manifest.yml sample
```
---
create-services:
- name:   "my-database-service"
  broker: "p-mysql"
  plan:   "1gb"

- name:   "Another-service"
  broker: "p-broker"
  plan:   "shared"
  parameters: "{\"RAM\": 4gb }"
```
