package grpcserver_test

import (
	"context"
	"testing"

	"github.com/Di-nis/shortener-url/internal/config"
	pb "github.com/Di-nis/shortener-url/internal/proto"
	grpcserver "github.com/Di-nis/shortener-url/internal/server/grpc"
	"github.com/Di-nis/shortener-url/internal/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestShortenerServiceServer_ShortenURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		useCase *usecase.URLUseCase
		config  *config.Config
		// Named input parameters for target function.
		in      *pb.URLShortenRequest
		want    *pb.URLShortenResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpcserver.NewShortenerServiceServer(tt.useCase, tt.config)
			got, gotErr := s.ShortenURL(context.Background(), tt.in)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ShortenURL() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ShortenURL() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ShortenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenerServiceServer_ExpandURL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		useCase *usecase.URLUseCase
		config  *config.Config
		// Named input parameters for target function.
		in      *pb.URLExpandRequest
		want    *pb.URLExpandResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpcserver.NewShortenerServiceServer(tt.useCase, tt.config)
			got, gotErr := s.ExpandURL(context.Background(), tt.in)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExpandURL() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExpandURL() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ExpandURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenerServiceServer_ListUserURLs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		useCase *usecase.URLUseCase
		config  *config.Config
		// Named input parameters for target function.
		in      *emptypb.Empty
		want    *pb.UserURLsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpcserver.NewShortenerServiceServer(tt.useCase, tt.config)
			got, gotErr := s.ListUserURLs(context.Background(), tt.in)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ListUserURLs() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ListUserURLs() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ListUserURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
