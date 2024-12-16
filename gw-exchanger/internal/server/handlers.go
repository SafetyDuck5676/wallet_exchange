package grpc

import (
	"context"

	pb "github.com/SafetyDuck5676/grpc_duck/proto-exchange/grpc/pb"
)

func (s *Server) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	rates, err := s.storage.GetExchangeRates()
	if err != nil {
		return nil, err
	}

	return &pb.ExchangeRatesResponse{Rates: rates}, nil
}

func (s *Server) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	rate, err := s.storage.GetExchangeRate(req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, err
	}

	return &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
	}, nil
}
