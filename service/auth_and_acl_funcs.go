package service

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	gomeetContext "github.com/gomeet/gomeet/utils/context"
	"github.com/gomeet/gomeet/utils/jwt"
	"github.com/gomeet/gomeet/utils/log"
)

func (s *profileServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	lFields := log.Fields{"fullMethodName": fullMethodName}
	log.Debug(ctx, "AuthFuncOverride call", lFields)

	newCtx, err := gomeetContext.ParseJWTFromContext(ctx, s.jwtSecret)
	if err != nil {
		if s.jwtSecret == "" ||
			fullMethodName == "/grpc.gomeetexamples.profile.Profile/Version" ||
			fullMethodName == "/grpc.gomeetexamples.profile.Profile/ServicesStatus" {
			return ctx, nil
		}

		log.Error(ctx, "Authentification failed", err, log.Fields{})

		return nil, err
	}

	return newCtx, nil
}

func (s *profileServer) AclFuncOverride(ctx context.Context, fullMethodName string, msg proto.Message) error {
	lFields := log.Fields{"fullMethodName": fullMethodName}
	log.Debug(ctx, "AclFuncOverride call", lFields)

	// return an error `errors.New("Error message")` to prevent the user from accessing this request
	if s.jwtSecret == "" ||
		fullMethodName == "/grpc.gomeetexamples.profile.Profile/Version" ||
		fullMethodName == "/grpc.gomeetexamples.profile.Profile/ServicesStatus" {
		return nil
	}

	jwtClaims, ok := ctx.Value("jwt").(jwt.Claims)
	if !ok {
		return errors.New("Invalid jwt")
	}

	lFields["jwtClaims"] = jwtClaims
	log.Debug(ctx, "AclFuncOverride call - allowed", lFields)

	// here the user is allowed from accessing this request
	return nil
}
