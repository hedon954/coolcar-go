package shared_auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	shared_token "coolcar/shared/auth/token"
	shared_id "coolcar/shared/id"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	HEADER_AUTHORIZATION = "authorization"
	TOKEN_PREFIX         = "Bearer "
)

type tokenVerifier interface {
	Verify(token string) (string, error)
}

type interceptor struct {
	verifier tokenVerifier
}

// Interceptor createa a grpc auth interceptro
func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	// read public key
	f, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot not public key file: %v", err)
	}
	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key file: %v", err)
	}

	// parse public key
	pk, err := jwt.ParseRSAPublicKeyFromPEM(fBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	i := &interceptor{
		verifier: &shared_token.JWTVerifier{
			PublicKey: pk,
		},
	}
	return i.HandleReq, nil
}

// func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// get token from context
	token, err := tokenFromContext(ctx)
	if err != nil {
		log.Fatalf("cannot get token: %v", req)
		return nil, status.Error(codes.Unauthenticated, "")
	}

	accountID, err := i.verifier.Verify(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
	}
	log.Printf("get accountID: %v", accountID)
	return handler(ContextWithAccountID(ctx, shared_id.AccountID(accountID)), req)
}

// tokenFromContext get token from context
func tokenFromContext(c context.Context) (string, error) {
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}

	// get token from header authorization
	token := ""
	for _, v := range m[HEADER_AUTHORIZATION] {
		if strings.HasPrefix(v, TOKEN_PREFIX) {
			token = v[len(TOKEN_PREFIX):]
		}
	}

	log.Printf("token: %v", token)

	if token == "" {
		return "", status.Error(codes.Unauthenticated, "")
	}

	return token, nil
}

type accountIDKey struct {
}

var aidKey = accountIDKey{}

// ContextWithAccountID returns a context with accountID
func ContextWithAccountID(c context.Context, accountID shared_id.AccountID) context.Context {
	return context.WithValue(c, aidKey, accountID)
}

// AccountIDFromContext returns account from income context
// returns unauthenticated error if no accountID in context
func AccountIDFromContext(c context.Context) (shared_id.AccountID, error) {
	v := c.Value(aidKey)
	accountID, ok := v.(shared_id.AccountID)
	if !ok {
		log.Printf("cannot get accountID: %v", c)
		return "", status.Error(codes.Unauthenticated, "")
	}
	return shared_id.AccountID(accountID), nil
}
