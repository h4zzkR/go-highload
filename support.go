package main

import (
	"crypto/rand"
	"encoding/base64"
	"hash/fnv"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Sith technique to generate UIDs based on usernames
// Why: https://auth0.com/learn/token-based-authentication-made-easy
func MakeToken(username string) string {
	// Generate a random 32-byte array
	tokenBytes := make([]byte, UidLength)
	rand.Read(tokenBytes)

	h := fnv.New32a()
	h.Write([]byte(username))
	usernameHash := h.Sum32()

	tokenBytes = append(tokenBytes, byte(usernameHash>>24))
	tokenBytes = append(tokenBytes, byte(usernameHash>>16))
	tokenBytes = append(tokenBytes, byte(usernameHash>>8))
	tokenBytes = append(tokenBytes, byte(usernameHash))

	token := base64.RawStdEncoding.EncodeToString(tokenBytes)

	return token
}

func StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	md, ok := metadata.FromIncomingContext(ss.Context())

	if !ok {
		return status.Error(codes.InvalidArgument, "metadata not found")
	}

	token, ok := md["authtoken"]
	if !ok || len(token) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization token")
	}

	username, ok := md["username"]
	if !ok || len(username) == 0 {
		return status.Error(codes.Unauthenticated, "missing username")
	}

	ctx := context.WithValue(ss.Context(), "authToken", token[0])
	ctx = context.WithValue(ctx, "username", username[0])

	return handler(srv, &WrappedStream{ss, ctx})
}

// A wrapper for the server stream that includes a modified context
type WrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *WrappedStream) Context() context.Context {
	return w.ctx
}
