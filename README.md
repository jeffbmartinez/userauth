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

### POST /verify/token/google

* The token should be sent over a secure connection from the user's browser to the service. This means using https for this resource. The [Let's Encrypt](https://letsencrypt.org/) project is a great way to get https enabled for your service.

Verifies a google ID token. Read up on google ID tokens in the following parts of google's documentation:

* [Integrating Google Sign-In into your web app](https://developers.google.com/identity/sign-in/web/sign-in)
* [Get Profile Information](https://developers.google.com/identity/sign-in/web/people)
* [Authenticate with a backend server](https://developers.google.com/identity/sign-in/web/backend-auth)

#### Request body

A json string containing the key `idtoken` and the google id token as the value. Something like this:

    { "idtoken": "[the id token goes here]" }

#### Response

JSON response consisting of the key `valid` which returns either `true` or `false` as a boolean. Example:

    { "valid" : true }

### POST /login/google

Create a new user session. If a previous session for this user already exists, it is expired and this new one is used.

#### Request body

A json string containing the key `idtoken` and the google id token as the value. Something like this:

    { "idtoken": "[the id token goes here]" }

#### Response

There are two relevant portions to the response:

1. The response body is a json string containing one key (`success`) with the value being `true` if the login request succeeded and `false` otherwise.
1. When the login request succeeds, two [httponly](https://www.owasp.org/index.php/HttpOnly) session cookies will be sent back. This means two `Set-Cookie` headers will be present in the response:
    * `sid` (session ID) - Contains a unique session ID string which should be set as a cookie in subsequent requests to the server. The session ID will be checked against the user's session ID.
    * `uid` (user ID) - Unique ID for the user. Along with the session ID, the user ID should be sent in a cookie in subsequent requests.

Example HTTP response:

```
HTTP/1.1 200 OK
[other headers...]
Content-Type: application/json
Set-Cookie: sid=[the session ID goes here]; HttpOnly
Set-Cookie: uid=[the user ID goes here]; HttpOnly

{ "success" : true }
```

Note that there is no expiration set for the cookies, meaning they are "session" cookies, which means they expire when the browser is closed.

### POST /logout

Logs a user out by invalidating their current session ID. This method relies on the session ID and user ID (as originally retrieved by, for example, the `/login/google` request) existing in the request cookies. The same keys of `sid` and `uid` are used.

There is no request body.

#### Request body

Any request body is simply ignored. The session ID and user ID to invalidate (log out from) must be present in the cookies, which would have been set by the original call to `/login/*`.

#### Response

A 200 OK response is the only relevant thing sent back. Any other response is an indicator that something went wrong.

### POST /verify/session

Verifies a session ID and user ID pair.

These checks are made:

* The session ID belongs to the user ID.
* The session ID has not expired.

* Note: This request is meant to be made from server side code, rather than any client side (browser javascript) code. Since the session ID and user ID cookies are set as [httponly](https://www.owasp.org/index.php/HttpOnly), it shouldn't be possible for the javascript code to know what these are in order to make the request.

#### Request Body

A JSON string containing:

* `sid` - The session ID.
* `uid` - The user ID.

Example:

```
{
  "sid": "[session ID goes here]",
  "uid": "[user ID goes here]"
}
```

#### Response

A JSON string containing:

* `valid` - boolean: true if the session ID and user ID are valid, false other wise.
* `reason` - string: If `value` is `true`, this will be an empty string. If `value` is false, `reason` contains an explanation for why the verification failed. It will contain one of:
    * `notfound` - No record was found of the session ID, user ID, or both.
    * `expired` - The session ID is expired.
    * `mismatch` - The session ID doesn't belong to the user ID


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
