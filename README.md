# goslack
A Slack client library for Golang, updated with new features as I need them

## Installation

    go get github.com/doozr/goslack

## What's included so far?

* Connecting to Slack and the Real Time API
* Getting a list of all users
* Getting a list of all non-Archived public channels
* Reading and Posting messages via the RTM

## Example

```go
connection, err := goslack.New(token)
if err != nil {
    log.Fatal(err)
}

for {
    event := <-connection.RealTime

    if event.Type == "message" {
        message, err := event.RtmMessage()
        // do something with message
    }

    if event.Type == "user_change" {
        userChange, err := event.RtmUserChange()
        // do something with userChange
    }
}
```
