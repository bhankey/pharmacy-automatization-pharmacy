package interceptors

import (
	"context"
	"errors"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/apperror"
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const RequestID = "x-request-id"

type ErrorHandlingInterceptor struct {
	log logger.Logger
}

func NewErrorHandlingInterceptor(log logger.Logger) *ErrorHandlingInterceptor {
	return &ErrorHandlingInterceptor{
		log: log,
	}
}

func (i *ErrorHandlingInterceptor) ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			errStatus, ok := status.FromError(err)
			if ok {
				return apperror.NewClientErrorFromGRPC(errStatus) // nolint: wrapcheck, nolintlint
			}

			return err
		}

		return nil
	}
}

func (i *ErrorHandlingInterceptor) ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		log := i.log.WithFields(logrus.Fields{
			"method": info.FullMethod,
		})

		if err != nil { // nolint: nestif, nolintlint
			requestID := ""
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				if headers, ok := md[RequestID]; ok && len(headers) > 0 {
					requestID = headers[0]
				}
			}

			if requestID == "" {
				log.Warnf("failed to get request id")
			}

			var clientError apperror.ClientError
			if errors.As(err, &clientError) {
				log.WithFields(logrus.Fields{
					"error":      clientError.ErrorToLog,
					"message":    clientError.Message,
					"code":       clientError.Code,
					"request_id": requestID,
				}).Warnf("response.client.error")

				return resp, clientError.GetGRPCError() // nolint: wrapcheck, nolintlint
			}

			log.WithFields(logrus.Fields{
				"error":      err,
				"request_id": requestID,
			}).Errorf("response.error")

			return resp, apperror.NewClientError(apperror.Common, err).GetGRPCError() // nolint: wrapcheck, nolintlint
		}

		return resp, nil
	}
}
