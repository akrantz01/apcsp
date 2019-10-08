# JSON Web Tokens
In order to authenticate users with the API routes, [JSON Web Tokens (JWT)](https://jwt.io) are issued after a user has logged in.
They contain the data on who the token is for and when it is valid.
The tokens are all signed to ensure that the data can be trusted and not be tampered with in transit or while being stored.

## Structure of a Token
The following is from the [introduction page](https://jwt.io/introduction/) on the JWT website. <br>
There are three parts that make up a JWT: `header`, `payload`, and `signature`.
Each section has been [base64](https://en.wikipedia.org/wiki/Base64) encoded with URL safe characters and separated by periods (`.`).
The end result, after concatenated, will look something like this:
```
Format:
header.payload.signature

Example:
eyJhbGciOiJIUzUxMiIsImtpZCI6IjVkOTEwMjc4M2NjYzAxM2EyNzFmMzQwOSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Njk4NzA4NDAsImlhdCI6MTU2OTc4NDQ0MCwibmJmIjoxNTY5Nzg0NDQwLCJzdWIiOiI1ZDhmMDU2NTQzZTM1MTViNzM1OTQ4NDAiLCJyb2xlIjoib3duZXIiLCJ0ZWFtIjoiNWQ4ZjA1NjU0M2UzNTE1YjczNTk0ODNmIn0.wyGS_7EuxLIYQH87dGVD2-QkcitxTax2b5i9BDqpUCJBUqX_USy6NyjB36o3b-5iLVEnhxwHBniM1QE6TKUetw
```
Here you can clearly see each segment, but not the data it contains.
You can copy and paste it into the [JWT debugger](https://jwt.io/#debugger-io) to see the data that it contains.

### The Header
The header contains the data to be able to identify the type of the token and the algorithm used to sign the token.
The type will have the key `typ` and will always be `JWT` as it will always be a JWT.
The algorithm is specified by the key `alg` which contains a 2 character code and the number of bits that the algorithm uses.
For example, the standard HMAC SHA256 algorithm will be represented as `HS256`.
To see the other algorithms that can be used, see [The Signature](#the-signature) section.
Once constructed, the JSON is base64 URL encoded.
Base64 URL encoding is the same as standard base64 encoding, but it uses a dash (`-`) and an underscore (`_`) instead to make sure it is URL safe.
<br>
Example:
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```
In this case, the type is a JWT, and the algorithm is HMAC SHA256.

### The Payload
The payload contains fields called claims.
Claims are statements about the user or resource, and additional data.
There are three different types of claims: `registered`, `public`, and `private`.
Registered claims are a set of predefined that are not mandatory but are highly suggested.
Registered claims include issuer (`iss`), subject (`sub`), audience (`aud`), expiration time (`exp`), not before time (`nbf`), issued at time (`iat`), and JWT id (`jti`).
The expiration, not before, and issued at times are all in seconds since January 1st, 1970 or the [Unix Epoch](https://en.wikipedia.org/wiki/Unix_time).
Public claims are defined in the [IANA JSON Web Token Registry](https://www.iana.org/assignments/jwt/jwt.xhtml) so as to avoid collisions.
Private claims are agreed upon between the parties that will be exchanging the information and are neither `registered` nor `public`.
Just as the header JSON, the payload JSON is base64 url encoded after it is constructed.
<br>
Example:
```json
{
  "exp": 1569870840,
  "iat": 1569784440,
  "nbf": 1569784440,
  "sub": "5d8f056543e3515b73594840",
  "role": "owner",
  "team": "5d8f056543e3515b7359483f"
}
```
In this case, the expiration time is 3 days after its issuance, and the issued at and not before times are the same as it is valid instantly after it is issued.
The subject is the id of the user that the token pertains to.
The role is the role that the specified user has and the team is the team that the user is on.
Both role and team are private claims as they are not defined in the IANA JWT Registry, and subject, expiration, not before, and issued at are registered claims.

### The Signature
The signature is the part of the token that allows for the verification of the data.
It is composed of the base64 url encoding of the header combined with the base64 url encoding of the payload and the combined with the secret.
```
Signing/Verification:
HMACSHA<bits>(
    base64UrlEncode(header) + "." + base64Url(payload),
    secret
)

Signing:
<RSA or ECDSA or RSA-PSS>SHA<bits>(
    base64UrlEncode(header) + "." + base64Url(payload),
    private_key
)

Verification:
<RSA or ECDSA or RSA-PSS>SHA<bits>(
    base64UrlEncode(header) + "." + base64Url(payload),
    public_key
)
```
The signature is used to ensure that the message was not changed while in transit.
In the case of a private key being used to sign the token, it can also verify the sender of the JWT.
<br><br>
Below is a list of the possible algorithms that can be used:
- Hash-based Message Authentication Code using Secure Hashing Algorithm (HMACSHA) (`HS256`, `HS384`, `HS512`)
- RSA Signature Scheme with Appendix using PKCS1 v1.5 (RSASSA-PKCS1-v1_5) (`RS256`, `RS384`, `RS512`)
- Elliptic Curve Digital Signature Algorithm (ECDSA) (`ES256`, `ES384`, `ES512`)
- RSA Signature Scheme with Appendix using Probabilistic Signature Scheme (RSASSA-PSS) (`PS256`, `PS384`)
The number after the 2 character code specifies the number of bits that the algorithm uses.
After each signature is calculated, it is hashed with SHA using the specified number of bits.

## How We Use It
We use JWTs to ensure that, when a user is authenticating with the API, they are who they say they are.
In our schema, we include the claims subject, issued at, not before, and expiration.
As per the [specification](https://tools.ietf.org/html/rfc7515#section-4.1.4), an additional parameter `kid`, or key id, can be specified in the case of dynamic secret generation.
The key id refers to the signing key in the database that corresponds with the generated token.
This is done in order to decrease the chance of fraudulent tokens being generated.
<br><br>
On every API request, except for the login and registration routes, an authorization token is required.
It is validated by middleware that ensures the route is not login or registration.
The token is passed in the headers portion of the request.
If the token is valid then the request will proceed, but if it is invalid, the status code will be `401 Unauthorized` with the reason why.
<br>
Example:
```
Authorization: xxxxxxxx.yyyyyyyy.zzzzzzzz
```
