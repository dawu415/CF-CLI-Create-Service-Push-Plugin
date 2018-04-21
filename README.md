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

# User-Provided Services
Support for user provided services is now available as of v1.1.0

Introduces a new field called `type`, which is used to indicate the type of service
that one wishes to instantiate.  If `type` is not specified, it is assumed that the service
is a standard brokered service from the market place.

`type` takes in the strings below, otherwise, it is assumed to be `brokered`

`brokered`:  A marketplace service broker

`credentials`: A user provided service, holding credentials, that is used for credentials and connection strings etc.  When this type is specified, a `credentials` field, along with arbitary number and type of parameters, is expected and must be provided.

`route`: A user provided route service. When this is specified, the field, 'url' with https schema must be provided.

`drain`: A user provided log drain service. When this is specified, the field, 'url' must be provided.


Sample service-manifest.yml
```
---
create-services:
- name:   "my-database-service"
  broker: "p-mysql"
  plan:   "1gb"
  
- name:   "Credentials-UPS"
  type:   "credentials"
  credentials:
    host: "https://abc.mydatabase.com/abcd"
    username: david
    password: 12.23@123password
    
- name:   "Route-UPS"
  type:   "route"
  url:    "https://www.google.com"
  
- name:   "LogDrain-UPS"
  type:   "drain"
  url:    "syslog-tls://server.myapp.com:1020"
  ```
