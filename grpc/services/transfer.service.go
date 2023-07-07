package grpc_server

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	auth "github.com/BogoCvetkov/go_mastercalss/grpc/interceptors"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/BogoCvetkov/go_mastercalss/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransferService struct {
	server interfaces.IGServer
	pb.UnimplementedTransferServiceServer
}

// Satisfy IGService interface
func (s *TransferService) PassServerConfig(server interfaces.IGServer) {
	s.server = server
}

func (s *TransferService) RegisterService() {
	pb.RegisterTransferServiceServer(s.server.GetGServer(), s)
}

func (s *TransferService) RegisterServiceOnGateway(c context.Context, mux *runtime.ServeMux) error {
	if err := pb.RegisterTransferServiceHandlerServer(c, mux, s); err != nil {
		return err
	}

	return nil
}

// Methods

func (s *TransferService) CreateTransfer(c context.Context, data *pb.CreateTransferReq) (*pb.CreateTransferRes, error) {
	reqU := auth.GetReqUser(c)

	acc, err := s.validAccount(c, data.FromAccountId, data.Currency)
	if err != nil {
		return nil, err
	}

	if acc.Owner != reqU.ID {
		return nil, status.Errorf(codes.Internal, "From acc does not belong to caller")
	}

	if _, err := s.validAccount(c, data.ToAccountId, data.Currency); err != nil {
		return nil, err
	}

	document := db.CreateTransferParams{
		FromAccountID: data.FromAccountId,
		ToAccountID:   data.ToAccountId,
		Amount:        data.Amount,
	}

	result, err := s.server.GetStore().TransferTrx(c, document)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	transfer := pb.Transfer{
		Id:            result.Transfer.ID,
		FromAccountId: result.FromAccount.ID,
		ToAccountId:   result.Transfer.ToAccountID,
		Amount:        result.Transfer.Amount,
		CreatedAt:     timestamppb.New(result.Transfer.CreatedAt),
	}

	return &pb.CreateTransferRes{Transfer: &transfer}, nil
}

func (ctr *TransferService) validAccount(c context.Context, accountID int64, currency string) (*db.Account, error) {
	account, err := ctr.server.GetStore().GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Account not found")
		}

		return nil, status.Errorf(codes.Internal, "Failed to get account")
	}

	if account.Currency != currency {
		msg := fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		return nil, status.Errorf(codes.Internal, msg)

	}

	return &account, nil
}
