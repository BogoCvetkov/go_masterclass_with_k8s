package grpc_server

import (
	"context"
	"database/sql"

	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	gen "github.com/BogoCvetkov/go_mastercalss/db/generated"
	auth "github.com/BogoCvetkov/go_mastercalss/grpc/interceptors"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/BogoCvetkov/go_mastercalss/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountService struct {
	server interfaces.IGServer
	pb.UnimplementedAccountServiceServer
}

// Satisfy IGService interface
func (s *AccountService) PassServerConfig(server interfaces.IGServer) {
	s.server = server
}

func (s *AccountService) RegisterService() {
	pb.RegisterAccountServiceServer(s.server.GetGServer(), s)
}

func (s *AccountService) RegisterServiceOnGateway(c context.Context, mux *runtime.ServeMux) error {
	if err := pb.RegisterAccountServiceHandlerServer(c, mux, s); err != nil {
		return err
	}

	return nil
}

// Methods

func (s *AccountService) CreateAccount(c context.Context, data *pb.CreateAccountReq) (*pb.CreateAccountRes, error) {
	reqU := auth.GetReqUser(c)

	if data.Currency == "" {
		return nil, status.Errorf(codes.Internal, "Currency not provided")
	}

	document := gen.CreateAccountParams{
		Owner:    reqU.ID,
		Currency: data.Currency,
		Balance:  0,
	}

	acc, err := s.server.GetStore().CreateAccount(c, document)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	account := pb.Account{
		Id:        acc.ID,
		Owner:     acc.Owner,
		Balance:   acc.Balance,
		Currency:  acc.Currency,
		CreatedAt: timestamppb.New(acc.CreatedAt),
	}

	return &pb.CreateAccountRes{Account: &account}, nil
}

func (s *AccountService) GetAccount(c context.Context, data *pb.GetAccountReq) (*pb.GetAccountRes, error) {
	reqU := auth.GetReqUser(c)

	id := data.GetId()

	params := db.GetAccountByOwnerParams{
		ID:    int64(id),
		Owner: reqU.ID,
	}

	acc, err := s.server.GetStore().GetAccountByOwner(c, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.Internal, "Account not found")

		}

		return nil, status.Errorf(codes.Internal, "Failed to get account")
	}

	account := pb.Account{
		Id:        acc.ID,
		Owner:     acc.Owner,
		Balance:   acc.Balance,
		Currency:  acc.Currency,
		CreatedAt: timestamppb.New(acc.CreatedAt),
	}

	return &pb.GetAccountRes{Account: &account}, nil

}

func (s *AccountService) ListAccounts(c context.Context, data *pb.ListAccountReq) (*pb.ListAccountRes, error) {

	reqU := auth.GetReqUser(c)

	// Default
	if data.Limit == 0 {
		data = &pb.ListAccountReq{
			Limit: 100,
			Page:  1,
		}
	}

	query := db.ListAccountsParams{
		Owner:  reqU.ID,
		Limit:  data.Limit,
		Offset: (data.Page - 1) * data.Limit,
	}

	accs, err := s.server.GetStore().ListAccounts(c, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.Internal, "No accounts found")

		}

		return nil, status.Errorf(codes.Internal, "Failed to get account")
	}

	accounts := []*pb.Account{}
	for _, a := range accs {
		account := pb.Account{
			Id:        a.ID,
			Owner:     a.Owner,
			Balance:   a.Balance,
			Currency:  a.Currency,
			CreatedAt: timestamppb.New(a.CreatedAt),
		}
		accounts = append(accounts, &account)
	}

	return &pb.ListAccountRes{Accounts: accounts}, nil

}
