package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bokunodev/hostapd211"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	remote_socket, err := hostapd211.FindRemoteSocket("/run/hostapd", "wlp2s0")
	if err != nil {
		panic(err)
	}

	client, err := hostapd211.NewClient(remote_socket)
	if err != nil {
		panic(err)
	}

	if err := client.Attach(ctx, 5*time.Second, callback); err != nil {
		log.Println("attach error:", err)
	}
}

func callback(ctx context.Context, msg string) error {
	fmt.Println("received:", msg)
	return nil
}
