# Configuration
All the configuration for the server is done through a library called [Viper](https://github.com/spf13/viper) which allows for the parsing of configuration files and accessing them in a simple and idiomatic way.
Viper supports reading configuration from JSON, TOML, YAML, HCL, envfiles, and Java properties files.
As well as, command line flags and environment variables.
Every field has a default configuration value if the field is not overridden by the configuration file.

## File Structure
The configuration file has three sections: `http`, `logging`, and `database`.
Each section is responsible for a different section of the server.
Below are the configuration keys and their descriptions.

### HTTP
This configures the host and port the server will listen on.
If the port is in use, then an error will be thrown and the server will exit.
It also configures the domain in which the server can be accessed.
This can either be an IP or domain, but must contain the scheme (either http or https).
Finally, you can specify whether to delete all of the uploaded files on the starting of the server.
You could also do this manually as the files are stored in the `uploaded` folder in the directory where the server is running.

### Logging
This configures the logger with the level and output format.
The log levels are `trace`, `debug`, `info`, `warn`, `error`, `fatal`, and `panic`.
The trace level is the most verbose and gives exact insight into what is happening in the server.
It should only be used for debugging as it generates mass amounts of logs.
The panic level is the least verbose and will never happen on this server.
It calls the `panic` function, causing a stack trace and the server to abruptly halt.
The potential formats are `text` and `json`.
The text format is colored and meant to be human readable.
As for JSON, it outputs the same data as text, but in JSON format meant for the parsing by a log aggregator like Prometheus or Logstash.

### Database
This configures the host and port to connect to the database.
If the server is unable to connect to the database, the server will throw an error and will exit.
You must also specify the username, password, and database to use.
If an authentication error occurs or the database does not exist, the the server will error and exit.
The configuration can specify the SSL connection mode to use.
The valid values are `disable`, `allow`, `prefer`, `require`, `verify-ca`, and `verify-full`.
Disable does not use SSL at all, allow will use it if available, prefer will attempt to use it, but won't fail if it cannot use SSL.
As for require, verify CA, and verify full, they will all enforce SSL, but to varying degrees.
Require does no validation on the certificates, verify CA ensures the certificate authority that issued the certificate is valid, and verify full ensures the entire chain is valid.

## Configuration Keys
Below are all the keys and their defaults in the configuration file.
The section is the enclosing field in which the keys exist.

| Section | Key | Type | Description | Default |
|---|---|---|---|---|
| http | host | string | Address to listen on. Use 0.0.0.0 for all addresses | 127.0.0.1 |
| http | port | integer | Port on the server to listen on | 8080 |
| http | domain | string | Domain/IP where the service is accessible | http://127.0.0.1:8080 |
| http | reset_files | boolean | Delete all of the files that have been uploaded | false |
| logging | format | string | Format to log the output in | text |
| logging | level | string | Set the minimum level to log | info |
| database | host | string | Address where the database can be accessed | 127.0.0.1 |
| database | port | integer | Port the database is listening on | 5432 |
| database | ssl | string | Mode to use SSL in | disable |
| database | username | string | User to access the database as | postgres |
| database | password | string | Password associated with the username | postgres |
| database | database | string | Database to write tables to | postgres |
| database | reset | boolean | Delete the tables if they already exist |

## Example
While Viper supports HCL, envfiles, and Java properties files, those configuration languages do not support nested values.
As such only JSON, TOML, and YAML are able to properly configure the server.
To see an example configuration file in YAML, view [config.sample.yaml](/api/config.sample.yaml).
