# Create-Service-Push

Reads in a services manifest yml file and creates the services listed and pushes 
an application.

This plugin extends cf push. So it will act the same way as cf push, with the exception
that it creates the services first.

# Additional Parameters
The create-service-push cli plugin supports the following optional parameters:

 * --service-manifest <MANIFEST_FILE>: Specify the fullpath and filename of the services creation manifest.  Defaults to services-manifest.yml
 * --no-service-manifest:              Specifies that there is no service creation manifest
 * --no-push:                          Create the services but do not push the application

# services-manifest.yml sample
```
---
create-services:
- name:   "my-database-service"
  broker: "p-mysql"
  plan:   "1gb"

- name:   "Another-service"
  broker: "p-brokerName"
  plan:   "sharedPlan"
  parameters: "{\"RAM\": 4gb }"
```
