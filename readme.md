# hostapd211

A Go client for interacting with the `hostapd` control interface using the `wpa_ctrl` protocol.


## Example usage:

```go
    client := hostapd211.NewClient("/var/run/hostapd/wlan0")
    reply, err := client.Ping(ctx)
```

## Installation

```bash
    go get github.com/bokunodev/hostapd211
```
