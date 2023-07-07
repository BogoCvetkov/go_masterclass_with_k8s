package grpc_server

import (
	"context"
	"database/sql"
	"time"

	"github.com/BogoCvetkov/go_mastercalss/auth"
	db "github.com/BogoCvetkov/go_mastercalss/db"
	gen "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/BogoCvetkov/go_mastercalss/pb"
	"github.com/BogoCvetkov/go_mastercalss/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	server interfaces.IGServer
	pb.UnimplementedUserServiceServer
}

// Satisfy IGService interface
func (s *UserService) PassServerConfig(server interfaces.IGServer) {
	s.server = server
}

func (s *UserService) RegisterService() {
	pb.RegisterUserServiceServer(s.server.GetGServer(), s)
}

func (s *UserService) RegisterServiceOnGateway(c context.Context, mux *runtime.ServeMux) error {
	if err := pb.RegisterUserServiceHandlerServer(c, mux, s); err != nil {
		return err
	}

	return nil
}

// Methods

func (s *UserService) CreateUser(c context.Context, data *pb.CreateUserReq) (*pb.CreateUserRes, error) {

	// Encrypt password
	hash, err := util.HashPassword(data.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	document := gen.CreateUserParams{
		Username:       data.Username,
		FullName:       data.FullName,
		Email:          data.Email,
		HashedPassword: hash,
	}

	// Create new user
	user, err := s.server.GetStore().CreateUser(c, document)
	if err != nil {

		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "Email already exists")
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := filterUser(&user)

	return &pb.CreateUserRes{User: result}, nil
}

func (s *UserService) LoginUser(c context.Context, data *pb.LoginUserReq) (*pb.LoginUserRes, error) {

	// Find user
	user, err := s.server.GetStore().GetUser(c, data.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Verify password
	err = util.CheckPassword(data.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Invalid credentials")
	}

	// Prepare access token payload
	p1, err := auth.NewTokenPayload(user.ID, s.server.GetConfig().TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed generating token")
	}

	// Generate access token
	token, err := s.server.GetAuth().GenerateToken(p1)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed generating token")
	}

	// Prepare REFRESH token payload
	p2, err := auth.NewTokenPayload(user.ID, s.server.GetConfig().RTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed generating refresh token")
	}

	// Generate REFRESH token
	rtoken, err := s.server.GetAuth().GenerateToken(p2)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed generating refresh token")
	}

	// Get Metadata
	md := util.ExtractMetaData(c)

	// Store session
	arg := gen.CreateSessionParams{
		ID:           p2.TokenID,
		UserID:       user.ID,
		RefreshToken: rtoken,
		UserAgent:    md.UserAgent,
		ClientIp:     md.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    p2.ExpiresAt.Time,
	}

	_, err = s.server.GetStore().CreateSession(c, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	filtered := filterUser(&user)

	return &pb.LoginUserRes{Token: token, RefreshToken: rtoken, User: filtered}, nil

}

func (s *UserService) RefreshToken(c context.Context, data *pb.RefreshTokenReq) (*pb.RefreshTokenRes, error) {
	// Verify refresh token
	payload, err := s.server.GetAuth().VerifyToken(data.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	// Check user session data
	session, err := s.server.GetStore().GetSession(c, payload.TokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		return nil, status.Errorf(codes.Internal, "authorization failed")
	}

	if payload.UserID != session.UserID || data.Token != session.RefreshToken {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if session.IsBlocked {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, status.Errorf(codes.Unauthenticated, "expired token")

	}

	// Prepare access token payload
	p, err := auth.NewTokenPayload(payload.UserID, s.server.GetConfig().TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed generating token")
	}

	// Generate access token
	token, err := s.server.GetAuth().GenerateToken(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed generating token")
	}

	return &pb.RefreshTokenRes{Token: token}, nil
}

func filterUser(u *gen.User) *pb.User {

	result := &pb.User{
		Id:        u.ID,
		CreatedAt: timestamppb.New(u.CreatedAt),
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
	}

	return result

}
