package api

import (
	"context"

	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
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

	// Create a query to select the launch IDs from a custom database view
	qry := query.New(
		"launch_validator",
		query.Fields("launch_id"),
		query.WithFilters(
			postgres.NewFilter("address", req.GetAddress()),
		),
		query.WithPageSize(page.GetSize()),
		query.AtPage(page.GetNumber()),
		query.SortByFields(query.SortOrderAsc, "launch_id"),
	)

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

	// Create a query to select the launch IDs from a custom database view
	qry := query.New(
		"launch_chain_created",
		query.Fields("launch_id"),
		query.WithFilters(
			postgres.NewFilter("coordinator_id", req.GetCoordinatorID()),
		),
		query.WithPageSize(page.GetSize()),
		query.AtPage(page.GetNumber()),
		query.SortByFields(query.SortOrderAsc, "launch_id"),
	)

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

	// Create a query to select the campaign IDs from a custom database view
	qry := query.New(
		"campaign_campaign_created",
		query.Fields("campaign_id"),
		query.WithFilters(
			postgres.NewFilter("coordinator_id", req.GetCoordinatorID()),
		),
		query.WithPageSize(page.GetSize()),
		query.AtPage(page.GetNumber()),
		query.SortByFields(query.SortOrderAsc, "campaign_id"),
	)

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
