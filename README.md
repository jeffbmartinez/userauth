# User Authentication System

Short name: Userauth

Userauth is responsible for user authentication. This includes the following things:

* Verifying OAuth2 tokens provided by third parties when users are signing in.
* Associating third party unique user IDs with a userauth unique user ID.
    * Looking up the userauth unique user ID from a third party user ID.
* Creating session IDs for users
* Verifying an existing session ID is valid (i.e. it exists and hasn't expired or been otherwise removed) for a particular user.
* Deleting/removing/expiring/invalidating session IDs

 The Storage Service consists of two parts:

* Userauth Service - The server which handles requests to do the above work.
* The storage system - For persisting session IDs as well as linking third party user IDs to userauth user IDs. Currently [redis](http://redis.io/) is used.

## Userauth API

### POST /login/google

Create a new user session. If a previous session for this user already exists, it is expired and this new one is used.

#### Request body

A json string containing the key `idtoken` and the google id token as the value. Something like this:

    { "idtoken": "[the id token goes here]" }

#### Response

There are two relevant portions to the response:

1. The response body is a json string containing one key (`success`) with the value being `true` if the login request succeeded and `false` otherwise.
1. When the login request succeeds, an [httponly](https://www.owasp.org/index.php/HttpOnly) and [secure](https://www.owasp.org/index.php/SecureFlag) cookie will be sent back. This means a `Set-Cookie` header will be present in the response:
    * `sid` (session ID) - Contains a unique session ID string which should be set as a cookie in subsequent requests to the server. The session ID will be checked against the user's session ID.

Example HTTP response:

```
HTTP/1.1 200 OK
[other headers...]
Content-Type: application/json
Set-Cookie: sid=[the cookie value goes here]; HttpOnly; Secure

{ "success" : true }
```

Note that there is no expiration set for the cookies, meaning they are "session" cookies, which means they expire when the browser is closed.

### POST /logout

Logs a user out by expiring their current session ID. If the user wasn't logged in to begin with they remain not logged in.

#### Request body

There is no request body. A user is logged out by expiring their session ID cookie.

#### Response

A 200 OK response is the only relevant thing sent back. Any other response is an indicator that something went wrong.

### POST /verify/session

Verifies a session ID.

* Note: This request is meant to be made from server side code, rather than any client side (browser javascript) code. Since the session ID cookie is set as [httponly](https://www.owasp.org/index.php/HttpOnly), it shouldn't be possible for the javascript code to know what these are in order to make the request.

#### Request Body

A JSON string containing:

* `sid` - The value of a session ID cookie as set by, for example, a call to `/login/*`.

Example:

```
{
  "sid": "[session ID cookie value goes here]",
}
```

#### Response

A JSON string containing:

* `valid` - boolean: true if the session ID is valid, false otherwise.
* `reason` - string: If `value` is `true`, this will be an empty string. If `value` is false, `reason` contains an explanation for why the verification failed. It will contain one of:
    * `notfound` - No record was found of the session ID.
    * `expired` - The session ID is expired.


```
{
  "valid" : false,
  "reason": "expired"
}
```

## Redis storage schema

[Redis](http://redis.io/) is used to persist the userauth data. If you're new to redis, try the [interactive redis tutorial](http://try.redis.io/) to get a feel for how it works.

Certain pieces of information are persisted in order to keep track of user authentication.

The following are stored in redis and used store information about third party IDs, userauth IDs, and session IDs.

### Third party ID Lookup

| Key Name | Value | Value Type | Description |
| --- | --- | --- | --- |
| google-id:{id} | userauth ID | string | Look up userauth ID by google ID |

### Session Information

| Key Name | Value | Value Type | Description |
| --- | --- | --- | --- |
| session:{session-id} | userauth ID | string | Look up a userauth ID by a session ID |
| user:{userauth-id} | session ID | string | Look up a session ID by a userauth ID |

When these entries are created they will be followed up with an [EXPIRE](http://redis.io/commands/expire) call to ensure they expire after a set amount of time. Once that happens the entries will be gone.
