package internalgrpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mmart-pro/mart-brute-blocker/internal/errors"
	"github.com/mmart-pro/mart-brute-blocker/internal/grpc/pb"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
)

type Service struct {
	pb.UnimplementedMBBServiceServer
	mbbService MbbService
}

type simpleSubnetRequest func(context.Context, model.Subnet) error

func NewService(mbbService MbbService) *Service {
	return &Service{
		mbbService: mbbService,
	}
}

func (s *Service) Health(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *Service) Allow(ctx context.Context, req *pb.SubnetReq) (*emptypb.Empty, error) {
	return s.callEmptyResponse(ctx, s.mbbService.Allow, req)
}

func (s *Service) Deny(ctx context.Context, req *pb.SubnetReq) (*emptypb.Empty, error) {
	return s.callEmptyResponse(ctx, s.mbbService.Deny, req)
}

func (s *Service) Remove(ctx context.Context, req *pb.SubnetReq) (*emptypb.Empty, error) {
	return s.callEmptyResponse(ctx, s.mbbService.Remove, req)
}

func (s *Service) Exists(ctx context.Context, req *pb.SubnetReq) (*pb.ExistsResponse, error) {
	r := model.NewSubnet(req.Subnet)
	if r == nil {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidSubnet.Error())
	}
	res, err := s.mbbService.Exists(ctx, *r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ExistsResponse{ListType: int32(res)}, nil
}

func (s *Service) Contains(ctx context.Context, req *pb.IpReq) (*pb.ContainsResponse, error) {
	ip := model.NewIPAddr(req.Ip)
	if ip == nil {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidIP.Error())
	}
	res, err := s.mbbService.Contains(ctx, *ip)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ContainsResponse{ListType: int32(res)}, nil
}

func (s *Service) ClearBucket(ctx context.Context, req *pb.ClearBucketRequest) (*emptypb.Empty, error) {
	ip := model.NewIPAddr(req.Ip)
	if ip == nil {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidIP.Error())
	}
	if req.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidLogin.Error())
	}
	if err := s.mbbService.ClearBucket(ctx, *ip, req.Login); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	ip := model.NewIPAddr(req.Ip)
	if ip == nil {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidIP.Error())
	}
	if req.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidLogin.Error())
	}
	if req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidPassword.Error())
	}
	check, err := s.mbbService.Check(ctx, *ip, req.Login, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CheckResponse{Allow: check}, nil
}

func (s *Service) callEmptyResponse(ctx context.Context, method simpleSubnetRequest, req *pb.SubnetReq) (*emptypb.Empty, error) {
	r := model.NewSubnet(req.Subnet)
	if r == nil {
		return nil, status.Errorf(codes.InvalidArgument, errors.ErrInvalidSubnet.Error())
	}
	if err := method(ctx, *r); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
