package main

import (
	"bufio"
	"fmt"
	pb "grpc-messenger/proto"
	"log"
	"os"
	"time"

	// "net"
	// "sync"
	// "time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type client struct {
	HostAdrr  string
	password  string
	username  string
	authToken string

	pb.MessengerClient
}

func NewClient(hostAdrr, username, password string) *client {
	return &client{
		HostAdrr:  hostAdrr,
		password:  password,
		username:  username,
		authToken: "",
	}
}

func (c *client) Run(ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting client request on %s", c.HostAdrr)

	conn, err := grpc.DialContext(ctx, c.HostAdrr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Client dialing failure: ", err)
		return err
	}

	defer conn.Close()
	c.MessengerClient = pb.NewMessengerClient(conn)

	if c.authToken, err = c.loginClient(ctx); err != nil {
		return err
	}

	log.Printf("Client %s logged in with %s token", c.username, c.authToken)

	err = c.Messenging(ctx)

	return err
}

func (c *client) Messenging(ctx context.Context) error {
	md := metadata.Pairs("authtoken", c.authToken, "username", c.username)
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	clientProxy, err := c.MessengerClient.MessageStream(ctx)
	if err != nil {
		return err
	}

	go c.receiveRoutine(clientProxy)
	err = c.sendRoutine(clientProxy)

	clientProxy.CloseSend()

	return err
}

// Read from terminal and send it via gRPC
func (c *client) sendRoutine(proxy pb.Messenger_MessageStreamClient) error {
	scanner := bufio.NewScanner(os.Stdin)

	c.printPrompt()

	for {
		select {
		case <-proxy.Context().Done():
			return nil
		default:
			ok := scanner.Scan()
			if !ok {
				err := scanner.Err()
				log.Printf("Scanner failure: %v", err)
				return err
			} else {
				err := proxy.Send(&pb.MSRequest{Message: scanner.Text()})
				if err != nil {
					log.Printf("Message sending failure: %v", err)
					return err
				}
			}
		}

	}
}

func (c *client) receiveRoutine(proxy pb.Messenger_MessageStreamClient) {
	for {
		in, err := proxy.Recv()
		if err != nil {
			return
		}

		c.printMessage(in)
		c.printPrompt()
		// TODO error handling
	}
}

func (c *client) printMessage(msg *pb.MSResponse) {
	ts := msg.Timestamp.AsTime().In(time.Local).Format("02-Jan-2006 15:04")
	fmt.Printf("\r%s [%s says:] %s\n", ts, msg.Message.Name, msg.Message.Content)
}

func (c *client) printPrompt() {
	// fmt.Printf("%s, write a message: ", c.username)
}

func (c *client) loginClient(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ClientConnectionTimeout)
	defer cancel()

	result, err := c.MessengerClient.Login(ctx, &pb.LoginRequest{
		Password: c.password,
		Username: c.username,
	})

	if err != nil {
		log.Fatal("client login failure: ", err)
	}

	return result.Token, err
}
