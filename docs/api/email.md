# Email
We use email to verify that it is an actual person signing up, reset passwords, and send basic notifications.
This is done using the [SMTP protocol](https://en.wikipedia.org/wiki/Simple_Mail_Transfer_Protocol) with a library called [gomail](https://godoc.org/gopkg.in/gomail.v2).
In order to send mail asynchronously from multiple functions at a time, we use [channels](https://tour.golang.org/concurrency/2).
These allow us to generate a message in a HTTP request and then send it whenever it is possible.

## Configuration
The configuration is done with the `email` block in the configuration file.
The `host` and `port` are self-explanatory, but it is worth noting that the value for the port is dependent on whether SSL is enabled or not.
If SSL is enabled, then the port should be `587`, and if not, it should be `25`.
In theory, any port could be used, but these are the standard ports.
The `username` and `password` correspond to the user you are trying to authenticate as for the server.
The `sender` is the email address that all email will be sent from.
This should be an email on a domain that you own. 
As stated before, `ssl` is a boolean key that specifies whether to use SSL in the connection or not.

## Availability
In order to ensure that the sending of an email is asynchronous, an anonymous function running in a goroutine is used.
This routine initializes a connection and then waits for messages to be received.
After a message has been received, the connection is re-opened if it was closed.
The connection is closed 4 seconds after opening to ensure that timeout errors don't occur with the SMTP server.
In general, 30 seconds would be the maximum connection time, but for AWS Simple Email Service, the timeout seems to be 4 seconds.
After the connection is re-opened, the message is written to the connection.
Assuming no errors occur, then the goroutine continues and waits for another message.
If, for whatever reason, an error occurs, then it will be logged with the specific error and time.

### Steps
A more concise description of the steps:
1. Message is sent through goroutine
1. Message is received
1. The connection re-opened, if it has been closed
1. The message is sent
1. Begins waiting for messages again
