package api

import (
	"context"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query/call"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ignite/backend/api/pb"
)

type server struct {
	pb.UnimplementedMapperServer

	db postgres.Adapter
}

func (s server) ChainsFromValidatorAddress(ctx context.Context, req *pb.ChainsFromValidatorAddressRequest) (*pb.ChainsFromValidatorAddressResponse, error) {
	page := req.GetPage()

	// Create a call query to select the launch IDs from a custom database view
	c := call.New("launch_validator", call.WithFields("launch_id"))
	qry := query.
		NewCall(c).
		AppendFilters(
			postgres.NewFilter("address", req.GetAddress()),
		).
		WithPageSize(page.GetSize()).
		AtPage(page.GetNumber())

	// Execute the query
	cr, err := s.db.Query(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}

	defer cr.Close()

	// Read the launch IDs from the query results
	res := pb.ChainsFromValidatorAddressResponse{}

	for cr.Next() {
		var launchID uint64

		if err := cr.Scan(&launchID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to read launch ID: %v", err)
		}

		res.LaunchIDs = append(res.LaunchIDs, launchID)
	}

	return &res, nil
}

func (s server) ChainsFromCoordinator(ctx context.Context, req *pb.ChainsFromCoordinatorRequest) (*pb.ChainsFromCoordinatorResponse, error) {
	page := req.GetPage()

	// Create a call query to select the launch IDs from a custom database view
	c := call.New("launch_chain_created", call.WithFields("launch_id"))
	qry := query.
		NewCall(c).
		AppendFilters(
			postgres.NewFilter("coordinator_id", req.GetCoordinatorID()),
		).
		WithPageSize(page.GetSize()).
		AtPage(page.GetNumber())

	// Execute the query
	cr, err := s.db.Query(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}

	defer cr.Close()

	// Read the launch IDs from the query results
	res := pb.ChainsFromCoordinatorResponse{}

	for cr.Next() {
		var launchID uint64

		if err := cr.Scan(&launchID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to read launch ID: %v", err)
		}

		res.LaunchIDs = append(res.LaunchIDs, launchID)
	}

	return &res, nil
}

func (s server) CampaignsFromCoordinator(ctx context.Context, req *pb.CampaignsFromCoordinatorRequest) (*pb.CampaignsFromCoordinatorResponse, error) {
	page := req.GetPage()

	// Create a call query to select the campaign IDs from a custom database view
	c := call.New("campaign_campaign_created", call.WithFields("campaign_id"))
	qry := query.
		NewCall(c).
		AppendFilters(
			postgres.NewFilter("coordinator_id", req.GetCoordinatorID()),
		).
		WithPageSize(page.GetSize()).
		AtPage(page.GetNumber())

	// Execute the query
	cr, err := s.db.Query(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}

	defer cr.Close()

	// Read the launch IDs from the query results
	res := pb.CampaignsFromCoordinatorResponse{}

	for cr.Next() {
		var campaignID uint64

		if err := cr.Scan(&campaignID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to read launch ID: %v", err)
		}

		res.CampaignIDs = append(res.CampaignIDs, campaignID)
	}

	return &res, nil
}
