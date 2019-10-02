# The Database
The database stores all user, chat, message, file, and token information in a persistent, non-volatile manner.
This is done so that the data can persist between the starting and stopping of server processes.
We use a [PostgreSQL](https://www.postgresql.org) database to store all the data.
In order to interact with the database, we use an Object-Relational Mapping, or ORM, called [GORM](https://gorm.io).
This allows us to be able to interact with the database easily.

## GORM
GORM can be found on [GitHub](https://github.com/jinzhu/gorm) and the general documentation on its [website](https://gorm.io/docs/).
As for the method and struct documentation, it can be found on [GoDoc](https://godoc.org/github.com/jinzhu/gorm).
It is a full featured ORM with support for associations, transactions, auto-migrations, and loading into structures.

### Connecting to a Database
Though we use PostgreSQL, it supports MySQL, PostgreSQL, Microsoft SQL Server, and SQLite.
In order to connect to each database, you must specify a connection string and the database type.
The connection string is specified based on the specification from the database:
- MySQL
    - [URL Format](https://dev.mysql.com/doc/connector-j/8.0/en/connector-j-reference-jdbc-url-format.html)
    - Remove the protocol portion of the connection string
- PostgreSQL
    - [URL Format](https://www.postgresql.org/docs/10/libpq-connect.html#id-1.7.3.8.3.5)
    - [Parameter Keywords](https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-PARAMKEYWORDS)
- SQLite
    - Specify a path where the database is currently or should be located
    - Alternatively, `:memory:` can be used instead of a file to use an in-memory database
      - __Note__: This is volatile as it is in memory
- Microsoft SQL Server
    - [URL Format](https://docs.microsoft.com/en-us/sql/connect/jdbc/building-the-connection-url?view=sql-server-2017)
    - Remove the `jdbc:` from the connection URL
    
### Associations
With all relational databases, there is the ability to create relationships between the data (per the name *relational*).
The supported associations are: belongs to, has one, has many, and many to many.
Belongs to and has one are essentially the same where one piece of data has another or belongs to another piece of data.
An example of this is when a user has a credit card associated with them or a credit card has a user.
The has many relationship is where one piece of data has many other pieces of data.
An example of this is when a user has many different posts.
Finally, the many to many relationship is where many pieces of data can have many other pieces of data and vice-versa.
This would be similar to where posts can have many tags and tags can have many posts.
When using this relationship, another table is implicitly created that contains two columns with the ids of the data being referenced. 

### CRUD Interface
The database has full CRUD, or create read update delete, abilities allowing for full interactions with the database.
Using defined structures, the developer is easily able to create a row in a table.
Simply by filling out the structure, and calling `db.Create(&row)` the new row can be inserted.
It also supports querying with full filtering abilities, including `WHERE`, `LIKE`, and `ORDER`.
There is also the ability to implicitly update a row buy calling `db.Save(&row)` with the specified information changed in the structure.
Finally, there is deletion of rows by calling `db.Delete(&row)`.

## Schema
Below is the schema definition for the database.
On each structure, the following fields are included automatically: `id`, `created_at`, `updated_at`, and `deleted_at`.
These fields allow for the abstraction of id sequence creation, timestamps of modifications and soft-deletion of data. 
<br><br>
Every table is created on startup of the server if it does not already exist.
Each of the tables is JSON serializable with certain fields excluded to ensure sensitive data is not given to the user.
<br><br>
In the `Name` column, if the contents are `implicit name` in italics, then it means that the field is only accessible during runtime.
In the `JSON Field Name` column, if the contents are `omitted` in italics, then it means that the field will not be included when serialized to JSON.

### Users
This table stores all user information and references the chats that a user has.
The user information includes their name, email, username, and hashed password.
The password is hashed using [Argon2i](https://en.wikipedia.org/wiki/Argon2) which is the currently recommended algorithm.
The relationship between the users and chats tables is a many to many relationship.

| Name | Type | Description | JSON Field Name |
|---|---|---|---|
| name | string | Recognizable name of the user | name |
| email | string | Email to contact the user at | email |
| username | string | Login name of the user | username |
| password | string | Password to identify the user | _omitted_ |
| _implicit name_ | many to many reference to chats | The chats the user is in | _omitted_ |

### Tokens
This table stores the signing key of the token and user it is for.
There is a belongs to relationship where the token belongs to a user.

| Name | Type | Description | JSON Field Name |
|---|---|---|---|
| signing_key | string | Base64 encoded 128-bit key the JWT is signed with | _omitted_ |
| user_id | unsigned integer | ID of the user the token is for | _omitted_ |

### Chats
This table stores the name and non-sequential id of the chat.
It also contains relationships between the users and chats, and the messages in the chat.

| Name | Type | Description | JSON Field Name |
|---|---|---|---|
| display_name | string | Human readable name of the chat | name |
| uuid | string | Non-sequential id of the chat for the API | uuid |
| _implicit name_ | many to many reference to users | The users in the chat | users |
| _implicit name_ | has many reference to messages | The messages in the chat | messages |

### Messages
This table stores message information, and references the user that sent it and any potential file associated with it.
The message information includes its type, the message itself, and the timestamp when it was sent.
There is a belongs to relationship where the message belongs to a user.
In addition, there is a has one relationship to a potential file it has.

| Name | Type | Description | JSON Field Name |
|---|---|---|---|
| chat_id | unsigned integer | ID of the chat the message was sent in | _omitted_ |
| sender_id | unsigned integer | ID of the user that sent the message | _omitted_ |
| type | unsigned integer | Content type of the message (0: text, 1: image, 2: type) | type |
| message | string | Text contained in the message | message |
| file_id | unsigned integer | ID of the file associated with the message | _omitted_ |
| _implicit name_ | has one reference to the file | The file (potentially) associated with the message | file |
| timestamp | 64-bit integer | When the message was sent in Unix time | timestamp |

### Files
This table stores file information and the chat it is apart of.
The file information includes its path on disk, the original file name (if it is a file), its non-sequential id, and whether it has been uploaded or not.

| Name | Type | Description | JSON Field Name |
|---|---|---|---|
| path | string | Where the file can be found on disk | _omitted_ |
| filename | string | The original name of the file | filename |
| uuid | string | Non-sequential id of the file for the API | uuid |
| used | boolean | Whether the file has already been uploaded | used |
| chat_id | unsigned integer | Chat the file is associated with | _omitted_ |
