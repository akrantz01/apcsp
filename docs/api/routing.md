# Routing
We route requests to their respective functions by the request path and the method being used.
This allows us to have a server do multiple things, rather than just having one ultra-complex route.
Routing by both path and method gives further specificity as each HTTP verb has its own intended use.
The API is designed in a [RESTful](https://restfulapi.net/) manner where there are routes to manage a set of a resource and to manage a single resource.

## Path-Based
Path-based routing is the most basic form of routing where a server receives some path following the domain, like `/api/users/1`, and acts on the data passed.
In order to do path routing, we use a library called [Mux](https://github.com/gorilla/mux) which has many routing capabilities though we only use the path routing.
The library builds on top of the standard `net/http` library by adding path routing and path parameters.
Path parameters are variables that are specified in the path instead of the body or query parameters.
For example, a path template could be `/users/{id}/posts`, where a possible request could be `/users/2/posts`.
In this case, the the variable would be `id` as it is surrounded by two curly braces `{}` and is able to change based on the given request.

We use path-based routing as the first level of specification.
In order to tell the server to create, read, update, or delete a resource, the proper path specifying the resource must be given.
Below are each of the paths for the four resources we have:
```
# User resource
/api/users
/api/users/{user}

# Chat resource
/api/chats
/api/chats/{chat}

# Message resource
/api/chats/{chat}/messages
/api/chats/{chat}/messages/{message}

# File resource
/api/files/{file}
```
In each of these, except for the files, there are two levels of specificity.
The first one pertaining to the entire group of the resource, and the second pertaining to a single resource.
The single resource is specified through path parameters which allows for the easy denotation of the resource.

## Method-Based
Method-based routing is just as simple as path-based routing.
While it is offered by Mux, we use a custom solution in order to implement it in the way we want.
Our solution uses a `switch` statement that assigns a certain function to a HTTP verb.
There are `GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `CONNECT`, `OPTIONS`, `TRACE`, and `PATCH` as defined in [MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods).
These are the possible options, but we only use `GET`, `POST`, `PUT`, and `DELETE`.

Each verb is used for a different purpose: `GET` for retrieving a resource, `POST` for creating, `PUT` for updating, and `DELETE` for deleting.
In general, `GET` and `POST` will be used together when referring to a group of resources, rather than a specific resource.
And for `GET`, `PUT`, and `DELETE`, they will be used together when referring to a specific resource.
As you can see with `GET`, it can be used in both general and single contexts because you can describe a single item and a list of items.

The switch statement works by checking the method of the request and calling a specific function based on the method.
In pseudo-code, the switch statement looks like this:
```
switch (request method) {
   case "GET":
       read_resource()
       break

   // NOTE: POST, PUT, and DELETE will not normally appear together as explained above
   case "POST":
       create_resource()
       break

   case "PUT":
       update_resource()
       break

   case "DELETE":
       delete_resource()
       break

   default:
       send_method_not_allowed()
}
```
For an implemented example, see [`users/exported.go`](/api/users/exported.go), [`chats/exported.go`](/api/chats/exported.go), [`files/exported.go`](/api/files/exported.go), and [`messages/exported.go`](/api/messages/exported.go).
