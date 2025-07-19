package gapi

import (
	"context"
	"time"

	db "github.com/akkahshh24/go-simple-bank-app/db/sqlc"
	"github.com/akkahshh24/go-simple-bank-app/pb"
	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/akkahshh24/go-simple-bank-app/val"
	"github.com/akkahshh24/go-simple-bank-app/worker"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Hash the password before storing it
	// Use the getter function to get the password from the request for beter validation
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		// If hashing fails, return an internal server error
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	// Create the user transaction
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		// If the transaction fails, check if it's a unique violation error
		if db.ErrorCode(err) == db.UniqueViolation {
			// If the username or email already exists, return an already exists error
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		// If there's another error, return an internal server error
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
