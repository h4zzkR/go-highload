package main

import (
	pb "grpc-messenger/proto"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	HostAdrr string
	password string

	userAuthTkns     map[string]string
	userMessagePipes map[string]chan *pb.MSResponse
	tell             chan *pb.MSResponse

	userTknsMtx         sync.RWMutex
	userMessagePipesMtx sync.RWMutex

	pb.UnimplementedMessengerServer
}

func NewServer(hostAdrr, password string) *server {
	return &server{
		HostAdrr: hostAdrr,
		password: password,

		tell: make(chan *pb.MSResponse),

		userMessagePipes: make(map[string]chan *pb.MSResponse),
		userAuthTkns:     make(map[string]string),
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting server on %s", s.HostAdrr)

	server := grpc.NewServer(grpc.StreamInterceptor(StreamInterceptor))
	pb.RegisterMessengerServer(server, s)

	lsner, err := net.Listen("tcp", s.HostAdrr)
	if err != nil {
		log.Fatal("Cannot bind on host adress.")
		return err
	}

	go s.tellRoutine(ctx)

	// Call Serve() on the server to do a blocking wait until the process is killed or Stop() is called.
	go func() {
		server.Serve(lsner)
		cancel()
	}()

	// Wait in Run thread till shutdown
	<-ctx.Done()

	server.GracefulStop()
	return nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Password != s.password {
		return nil, status.Error(codes.Unauthenticated, "password is wrong")
	}

	token, err := s.newUserSession(req.Username)

	if err != nil {
		return nil, err
	}

	log.Printf("%s has entered session	 with %s token", req.Username, token)

	return &pb.LoginResponse{Token: token}, nil
}

/*
*
 */
func (s *server) MessageStream(srv pb.Messenger_MessageStreamServer) error {
	token, ok := srv.Context().Value("authToken").(string)
	username, ok2 := srv.Context().Value("username").(string)

	if !ok || !ok2 {
		return status.Error(codes.Unauthenticated, "invalid token or username metadata")
	}

	if err := s.authorizeRequest(username, token); err != nil {
		return err
	}

	log.Printf("User %s started stream", username)

	defer srv.Context().Done()

	// Get pipe where will put messages
	thisClientPipe := s.connectPipe(token)

	// Run processor for this client
	go s.fromPipeSender(srv, thisClientPipe)

	for {
		// Got new message:
		request, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("User %s sent message", username)

		s.tell <- &pb.MSResponse{
			Timestamp: timestamppb.Now(),
			Message: &pb.MSResponse_Message{
				Name:    username,
				Content: request.Message,
			}}

	}

	// return stream.Context().Err()
}

func (s *server) newUserSession(username string) (string, error) {
	s.userTknsMtx.RLock()

	if _, found := s.userAuthTkns[username]; found {
		return "", status.Error(codes.AlreadyExists, "user with that username is already exists in this chat")
	}

	authToken := MakeToken(username)
	s.userAuthTkns[username] = authToken

	s.userTknsMtx.RUnlock()

	return authToken, nil
}

func (s *server) authorizeRequest(username, authToken string) error {
	s.userTknsMtx.RLock()
	if realToken, found := s.userAuthTkns[username]; !found {
		return status.Error(codes.NotFound, "user with this username isn't authorized")
	} else if realToken != authToken {
		return status.Error(codes.PermissionDenied, "wrong token for this username")
	}

	return nil
}

// This thing sends messages from broadcast pipe (s.tell) to all clients
func (s *server) tellRoutine(ctx context.Context) {
	for {
		response := <-s.tell
		s.userMessagePipesMtx.RLock()

		for _, clientPipe := range s.userMessagePipes {
			select {
			case clientPipe <- response:
			default:
				// ignore message bcz clientPipe is full
			}
		}

		s.userMessagePipesMtx.RUnlock()
	}
}

// Holds pipe and srv for each client, sends new messages to them
func (s *server) fromPipeSender(srv pb.Messenger_MessageStreamServer, pipe chan *pb.MSResponse) {
	for {
		select {
		case <-srv.Context().Done():
			return
		case response := <-pipe:
			if status, ok := status.FromError(srv.Send(response)); ok {
				code := status.Code()
				if code != codes.OK {
					log.Printf("sending the message to the client failed with code %s", code)
					return
				}
			}

		}
	}
}

func (s *server) connectPipe(token string) chan *pb.MSResponse {
	s.userMessagePipesMtx.Lock()
	pipe, found := s.userMessagePipes[token]

	if !found {
		pipe = make(chan *pb.MSResponse, ClientStreamCapacity)
		s.userMessagePipes[token] = pipe
	}

	s.userMessagePipesMtx.Unlock()
	return pipe
}
