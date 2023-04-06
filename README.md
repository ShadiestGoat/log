# Shady Logs

[![Go Reference](https://pkg.go.dev/badge/github.com/shadiestgoat/log.svg)](https://pkg.go.dev/github.com/shadiestgoat/log)

This is a *highly* opinionated yet expandable logger. The main focus of this logger is easy of use, which is why this does not have logger instances, but rather a global logger (all of which log to the same set of outputs).

To get started, you have to initialize the logger along with the desired outputs:

```go
package main

import "github.com/shadiestgoat/log"

func main() {
    log.Init(
        log.NewLoggerPrint(), 
        log.NewLoggerFile("log.txt"), 
        log.NewLoggerDiscordWebhook("<@376079696489742338>,", "https://canary.discord.com/api/webhooks/{CHANNEL}/{TOKEN}"),
        // And maybe some custom functions?
    )
}
```

After that, you can call global package functions to log

```go
package main

import "github.com/shadiestgoat/log"

func main() {
    log.Init(log.NewLoggerPrint())

    log.Success("Logger has started!")
    log.Warn("I can also do warnings <3")
    log.Error("Aaaa!! Errorss!!!!!!")
    g()
}

func g() {
    log.Fatal("Oh no, I shall now proceed to panic :(")
}
```

There are also util functions for error handling

```go
// ... in a request handler for example
postID := "abc123"

post, err := apiWrapper.FetchPost(postID) // example api wrapper, that could return an error

// If err != nil: this will log "While fetching post 'abc123': {whatever the error message is}", and exit the req handler. Outputs the level 'ERROR'
// If err == nil: Does nothing
if log.ErrorIfErr(err, "fetching post '%s'", postID) {
    return
}
```

There is also `log.FatalIfErr(...)` which acts in the same way way that `log.ErrorIfErr` acts, but, in the end it `panic`s (so no need to do error checking).


```go
// ... in a request handler for example
postID := "abc123"

post, err := apiWrapper.FetchPost(postID) // example api wrapper, that could return an error

// If err != nil: this will log "While fetching post 'abc123': {whatever the error message is}", then crash the programme through panic(). Outputs the level 'ERROR'
// If err == nil: Does nothing
log.FatalIfErr(err)
```
