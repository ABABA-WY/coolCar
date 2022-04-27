package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"os"
	"strings"
)

type tokenVerifier interface {
	Verify(token string) (string, error)
}
type interceptor struct {
	verifier tokenVerifier
}

// Interceptor creates an auth interceptor
func Interceptor(publicFile string) (grpc.UnaryServerInterceptor, error) {
	file, err := os.Open(publicFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key from publicFile:%v\n", err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot readall from file:%v\n", err)
	}
	pKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key:%v\n", err)
	}
	i := &interceptor{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pKey,
		},
	}
	return i.HandleRequest, nil
}

// HandleRequest req：请求，info：相关信息，handler：接下来要做的事
func (i *interceptor) HandleRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//获取token
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "tokenFromContext error")
	}
	aid, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid:%d", err)
	}
	return handler(ContextWithAccountID(ctx, aid), req)
}

func tokenFromContext(c context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "")
	}
	tkn := ""
	for _, v := range md["authorization"] {
		if strings.HasPrefix(v, "Bearer ") {
			tkn = v[len("Bearer "):]
		}
	}
	if tkn == "" {
		return "", status.Errorf(codes.Unauthenticated, "token is nil")
	}
	return tkn, nil
}

type accountIDKey struct {
}

// ContextWithAccountID creates a context with a given account id
func ContextWithAccountID(c context.Context, aid string) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}

// AccountIDFromContext gets account id from context
//returns Unauthenticated error if no account id is available
func AccountIDFromContext(c context.Context) (string, error) {
	v := c.Value(accountIDKey{})
	aid, ok := v.(string)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "has not account id key value")
	}
	return aid, nil
}
