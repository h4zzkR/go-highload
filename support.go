package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"hash/fnv"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "grpc-messenger/proto"

	"github.com/redis/go-redis/v9"
)

func CacheMessage(client *redis.Client, msg *pb.MSResponse) {
	// fmt.Print(msg.Timestamp.AsTime().String())
	key := msg.Message.Name + "\t" + msg.Timestamp.AsTime().Format(time.RFC3339)
	content := msg.Message.Content
	ctx := context.Background()

	_, err := client.HSet(ctx, RedisCacheDBName, []string{key, content}).Result()
	if err != nil {
		log.Fatalf("Error adding %s", content)
	}
}

func UnCacheMessage(client *redis.Client, msg *pb.MSResponse) {
	// score := float64(msg.Timestamp.AsTime().Unix())
	// content := msg.Message.Name + "\t" + msg.Message.Content
	// ctx := context.Background()

	// _, err := client.ZAdd(ctx, "history", redis.Z{Score: score, Member: content}).Result()
	// if err != nil {
	// 	log.Fatalf("Error adding %s", content)
	// }
}

func GetCachedHistory(client *redis.Client) map[string]string {
	ctx := context.Background()
	return client.HGetAll(ctx, RedisCacheDBName).Val()
}

func readUnderlying(lines chan interface{}) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines <- s.Text()
	}
	lines <- s.Err()
}

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

	log.Printf("Interceptor: user %s | token %s", username[0], token[0])

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
