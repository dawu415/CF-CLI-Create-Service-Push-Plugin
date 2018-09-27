# Create-Service-Push

Reads in a services manifest yml file and creates the services listed and pushes 
an application.

This plugin extends cf push. So it will act the same way as cf push, with the exception that it creates the services first.

# Additional Parameters
The create-service-push cli plugin supports the following optional parameters:

 version  1.0.0 and above
 ------------------------ 
 * `--service-manifest MANIFEST_FILE`: Specify the fullpath and filename of the services creation manifest.  Defaults to services-manifest.yml.

 * `--no-service-manifest`:              Specifies that there is no service creation manifest.

 * `--no-push`:                          Create the services but do not push the application.

 version  1.3.0 and above
 ------------------------ 
 * `--use-env-vars-prefixed-with PREFIX`: Allows services-manifest to use environment variables that have a prefix `PREFIX` for variable substitution. 

 * `--var KEY=VALUE`: Sets a single variable to used to be substituted into the services-manifest. Can be specified multiple times. Does not propagate to cf push unless the `--push-as-subprocess` is specified. 

 * `--vars-file FULLPATH_FILENAME`: Sets a single vars file containing variables to be used for substitution into a services-manifest.  Can be specified multiple times. Does not propagate to cf push unless the `--push-as-subprocess` is specified.

 * `--push-as-subprocess`: Forces the cf push to be called from the OS as a sub-process rather than use the internal CF CLI Plugin Command architecture.  This makes the assumption that there will be a file named `cf` or `cf.exe` that can be found in the current working directory or in the System's PATH environment variable. This was introduced to work around a temporary breaking refactor change in the CF CLI plugin architecture where new features such as `--var` could not be directly used.

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
## Support for user provided services is now available as of v1.1.0

Introduces a new field called `type`, which is used to indicate the type of service
that one wishes to instantiate.  If `type` is not specified, it is assumed that the service
is a standard brokered service from the market place.

`type` takes in the strings below, otherwise, it is assumed to be `brokered`

`brokered`:  A marketplace service broker

`credentials`: A user provided service, holding credentials, that is used for credentials and connection strings etc.  When this type is specified, a `credentials` field, along with arbitary number and type of parameters, is expected and must be provided.

`route`: A user provided route service. When this is specified, the field, `url` with *https* schema must be provided.

`drain`: A user provided log drain service. When this is specified, the field, `url` must be provided.


Sample services-manifest.yml showing support for brokered and user provided services of all types. 
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

# Updating Services
## Support for updateService is available as of 1.2.0

Updating of services is implemented here for only the service `Parameters` and `Tags`. Please see Tags section below for more information on limitations. Updating of services works for brokered, as well as all user-provided service types.

Updating of service plans is not supported for safety reasons where it maybe better
do this in a controlled manner. (Happy to take on requests otherwise :-) )

The flag to use in the services-manifest file is `updateService: <bool>`.
By default, this is set to `False`. 

Example `services-manifest.yml`
```
---
create-services:
- name:   "my-configserver"
  broker: "p-config-server"
  plan:   "standard"
  updateService: true
```

# Tags
## Support for tags is available as of 1.2.0

Tags are available for Brokered services only.  
The CF CLI v6.37 currently does not support tags for user-provided services...yet!

The flag to use in the services-manifest file is `tags: "comma separated <string>"`.
By default, if not `tags` are provided, the service is created without the tags, i.e., don't include the `-t` in the service creation command. 

Example `services-manifest.yml` for tags

```
---
create-services:
- name:   "my-configserver"
  broker: "p-config-server"
  plan:   "standard"
  tags:   "Something, ConfigServer, appname-config-server"
```


# Variable Substitution
## Support for variable substitution is available as of 1.3.0

Variables can be specified via the parameter inputs `--var KEY=VALUE` and `--vars-file FULLPATH_FILENAME`, which sets a single variable or a specifies a file containing a set of variables that can be substituted into a services-manifest, respectively. Note that these flags do not get passed to cf push, unless `--push-as-subprocess` is specified.

In addition to the above, 1.3.0 also allows environment variables to be substituted in. This is done via the flag `--use-env-vars-prefixed-with PREFIX`. This is limited by environment variables that have been created with some prefix, e.g. APPNAME can be a prefix that one sets in the environment variable such as APPNAME_APPNAME_configserverID and will be picked up the plugin.

Sample services-manifest using enviroment variable
```
---
create-services:
- name:   "((environment))-configserver"
  broker: "p-config-server"
  plan:   "standard"
  tags:   "((APPNAME_configserverID)), ConfigServer, appname-config-server"
```