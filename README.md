<p align="center">
  <a href="https://goreportcard.com/report/github.com/dimmerz92/sesh">
    <img src="https://goreportcard.com/badge/github.com/dimmerz92/sesh" alt="Go Report Card" />
  </a>
  <a href="https://pkg.go.dev/github.com/dimmerz92/sesh">
    <img src="https://pkg.go.dev/badge/github.com/dimmerz92/sesh" alt="Go Reference" />
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT License" />
  </a>
  <a href="https://github.com/dimmerz92/sesh">
    <img src="https://img.shields.io/badge/Go%20Modules-Yes-green.svg" alt="Go Module"/>
  </a>
</p>

# Sesh

Sesh provides a fast and simple session store with for Go back ends, powered by [BadgerDB](https://github.com/hypermodeinc/badger).

**Key features:**
- Powered by BadgerDB, an LSM key/value database (very fast!)
- In memory or persistent storage.
- Optional cookie functionality.
- Context setting middleware.
- Session configurability.
- Context configurability.
- Cookie configurability.

# Installation

Add to your project:

```bash
go get github.com/dimmerz92/sesh@latest
```

# Examples

## Without Cookies

```go
import (
    "fmt"

    "github.com/dimmerz92/sesh"
)

type MyCustomSessionData struct {
    Hello         string
    MeaningOfLife int
}

var SessionStore *sesh.SessionStore

func DoStuff() {
    // Create a new store for the sessions (here we ignore the error for the example, but make sure to check it!).
    // This will attempt to read (or create if not exists) the session_data directory in the root of your project.
    // Attributes within the DefaultConfig can be changed using the WithXXX chain methods.
    // Alternatively, a custom Config can be constructed with Config{}.
    SessionStore, _ = sesh.NewSessionStore(sesh.DefaultConfig())

    // Create a new session using the New method of the SessionStore and provide your data.
    // New can handle any kind of data from primitives through to maps and custom structs.
    // Again here we ignore the returned error for the example, but please handle this.
    // New will return a string representation of a uuid which is the session ID.
    sessionId, _ := SessionStore.New(MyCustomSessionData{Hello: "world", MeaningOfLife: 42})

    // Retrieve session data by creating a data pointer and supplying it to the Get method of the SessionStore.
    // This workflow is similar to decoding/unmarshaling for json data, so shouldn't be too unfamiliar.
    var data MyCustomSessionData
    if err := SessionStore.Get(sessionId, &data); err == nil {
        fmt.Printf("Hello: %s, Meaning of life: %d", data.Hello, data.MeaningOfLife)
    }

    // Delete a specific session by simply providing the session ID.
    // It is possible but highly unlikely that the Delete method will error, but it should be checked.
    SessionStore.Delete(sessionId)

    SessionStore.Close()
}
```

Here we first initialise a session store with the default config. The default config has some chain methods that can be helpful for changing one or two default parameters. These can be changed with the `WithXXX` methods. See [config.go](https://github.com/dimmerz92/sesh/blob/master/config.go) to see the configurability of the config.

After this, we added a new session and provided our own custom data struct. Any type of data can be added to a session (from primitives to maps and structs), it is recommended that you define a schema that works best for your project and stick with that.

Then we retrieve the session data that we just created. This process is similar to unmarshaling of json data. A pointer needs to be provided to the function for it to unpack the session data in to.

Finally, we deleted the session using the session ID in the delete method. This is highly unlikely to error, but it is possible. An error is most likely going to be due to a database connection issue, and should be handled accordingly.

Don't forget to close the session store when you're done with it. If you opted for persistent sessions (default), then they will be saved in the default session_data folder. It's recommended to only close once, especially if you have a long lived app (i.e., a HTTP server), keeping it open isn't harmful.

## With Cookies

```go
// Similar to the example without cookies, these functions are a wrapper of those functions that set, update, and invalidate cookies as well.
// These functions work in conjunction with the SessionChecker middleware, which checks for cookies and sets the session data on the request context.
func main() {
    sessionStore, _ := sesh.NewSessionStore(sesh.DefaultConfig())

    mux := http.NewServeMux()

    mux.HandleFunc("GET /", MyHandler)

    server := &http.Server{
        Addr:    ":8000",
        // Here we wrap the mux in the session checker middleware. We could alternatively only wrap handlers we want checked.
        // The session checker requires a data type so it knows how to handle the unpacking of the session data.
        // This does restrict all sessions to only containing the same data type, so design your schema accordingly.
        // The session checker middleware will check for the presence of a session cookie, get the session data and add it to the request context.
        Handler: sesh.SessionChecker[MyCustomSessionData](mux, sessionStore),
    }
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Create a new session with a cookie.
    // Can return errors, ignored here, but please check.
    sessionId, _ := SessionStore.NewWithCookie(w, MyCustomSessionData{Hello: "world", MeaningOfLife: 42})

    // Retrieve the data for a session.
    // Checks for the presence of a cookie and extracts the data into your provided pointer.
    // If a valid session is found, the cookie expiry is updated if configured (default).
    // Can return errors, ignored here, but please check.
    var data MySessionData
    SessionStore.GetWithCookie(w, r, &data)

    // Delete a session.
    // Checks for the presence of a cookie, deletes the session and invalidates the cookie if it exists.
    // Similar to the non-cookie function, this can error but it is unlikely. Please check.
    SessionStore.DeleteWithCookie(w, r)
}
```

Here we wrapped our server mux in the session checker middleware. This middleware will check all requests for the given handler or mux for the presence of session cookies, get the session from the provided session store and set the data on the request context. The middleware requires a type so it knows how to extract the data. This limits all sessions to the same datatype, so design your session schema accordingly.

Following this, we registered a handler that creates a new session with the previously defined custom data struct, gets the session, and deletes it. Through all of these calls, the session ID is not needed as the functions check for the presence of a session cookie on the request.

# License

MIT licensed. See the LICENSE file for details.
