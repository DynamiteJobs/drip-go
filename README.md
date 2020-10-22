# drip-go
Wrapper for the Drip API https://www.getdrip.com/docs/rest-api

[godoc](https://pkg.go.dev/github.com/dynamite-jobs/drip-go)

# Example
```go
dripClient, err := drip.New("DRIP_API_KEY", "DRIP_ACCOUNT_ID")
if err != nil {
    ...
}
req := &drip.ListSubscribersReq{}
resp, err := dripClient.ListSubscribers(req)
if err != nil {
    ...
}
if len(resp.Errors) > 0 {
    ... // resp.Errors has Drip API errors. https://www.getdrip.com/docs/rest-api#errors
}

// resp.Subscribers has list of subscribers
```

Look at test for more examples.

# Contributions
Please feel free to add a pull request for any feature not available. Pull requests should include test.
