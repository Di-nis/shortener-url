package grpcserver_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/mocks"
	"github.com/Di-nis/shortener-url/internal/models"
	pb "github.com/Di-nis/shortener-url/internal/proto"
	grpcserver "github.com/Di-nis/shortener-url/internal/server/grpc"
	"github.com/Di-nis/shortener-url/internal/toolkit"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	UUID = "01KA3YRQCWTNAJEGR5Z30PH6VT"

	urlOriginal1 = "https://www.khl.ru/"
	urlShort1    = "lJJpJV7h"
	urlOriginal2 = "https://www.dynamo.ru/"
	urlShort2    = "kiFL71uv"

	urlIn1 = models.URLJSON{
		UUID:     UUID,
		Original: urlOriginal1,
	}

	urlIn2 = models.URLJSON{
		UUID:     UUID,
		Original: urlOriginal2,
	}

	urlOut1 = models.URLBase{
		UUID:        UUID,
		URLID:       "1",
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}

	urlOut2 = models.URLBase{
		UUID:        UUID,
		URLID:       "2",
		Original:    urlOriginal2,
		Short:       urlShort2,
		DeletedFlag: false,
	}

	urlsOut = []models.URLBase{urlOut1, urlOut2}
)

var cfg *config.Config

func setEnv() {
	var err error

	err = os.Setenv("SERVER_ADDRESS", "localhost:8080")
	if err != nil {
		log.Fatalf("set env SERVER_ADDRESS failed: %v", err)
	}

	err = os.Setenv("BASE_URL", "http://localhost:8080")
	if err != nil {
		log.Fatalf("set env BASE_URL failed: %v", err)
	}

	err = os.Setenv("ENABLE_GRPC", "true")
	if err != nil {
		log.Fatalf("set env ENABLE_GRPC failed: %v", err)
	}
}

func TestMain(m *testing.M) {
	setEnv()

	cfg = config.NewConfig()
	cfg.Load()

	os.Exit(m.Run())
}

func newTestService(t *testing.T) (*grpcserver.ShortenerServiceServer, *mocks.MockURLUseCase) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockUseCase := mocks.NewMockURLUseCase(ctrl)

	s := grpcserver.NewShortenerServiceServer(mockUseCase, cfg)
	return s, mockUseCase
}

func TestShortenerServiceServer_ShortenURL(t *testing.T) {
	fullUrlShort1 := toolkit.AddBaseURLToResponse(cfg.BaseURL, urlShort1)

	tests := []struct {
		name    string
		mock    func(*mocks.MockURLUseCase)
		in      *pb.URLShortenRequest
		want    *pb.URLShortenResponse
		wantErr error
	}{
		{
			name: "URL создан успешно",
			mock: func(mock *mocks.MockURLUseCase) {
				mock.EXPECT().CreateURLOrdinary(gomock.Any(), urlIn1).Return(urlOut1, nil)
			},
			in: pb.URLShortenRequest_builder{
				Url: &urlOriginal1,
			}.Build(),
			want: pb.URLShortenResponse_builder{
				Result: &fullUrlShort1,
			}.Build(),
			wantErr: nil,
		},
		{
			name: "URL уже существует",
			mock: func(mock *mocks.MockURLUseCase) {
				mock.EXPECT().CreateURLOrdinary(gomock.Any(), urlIn2).Return(urlOut2, constants.ErrorURLAlreadyExist)
			},
			in: pb.URLShortenRequest_builder{
				Url: &urlOriginal2,
			}.Build(),
			want:    nil,
			wantErr: status.Errorf(codes.AlreadyExists, `URL %s already exist`, urlOriginal2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockUseCase := newTestService(t)
			tt.mock(mockUseCase)

			ctx := context.WithValue(context.Background(), constants.UserIDKey, UUID)

			got, gotErr := s.ShortenURL(ctx, tt.in)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortenURL() = %v, want %v", got, tt.want)
			}

			st, ok := status.FromError(gotErr)
			if !ok {
				t.Fatalf("expected grpc status error")
			}
			if st.Code() != status.Code(tt.wantErr) {
				t.Errorf("code = %v, want %v", st.Code(), status.Code(tt.wantErr))
			}
		})
	}
}

func TestShortenerServiceServer_ExpandURL(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(*mocks.MockURLUseCase)
		in      *pb.URLExpandRequest
		want    *pb.URLExpandResponse
		wantErr error
	}{
		{
			name: "Оригинальный URL получен",
			mock: func(mock *mocks.MockURLUseCase) {
				mock.EXPECT().GetOriginalURL(gomock.Any(), urlShort1).Return(urlOriginal1, nil)
			},
			in: pb.URLExpandRequest_builder{
				Id: &urlShort1,
			}.Build(),
			want: pb.URLExpandResponse_builder{
				Result: &urlOriginal1,
			}.Build(),
			wantErr: nil,
		},
		{
			name: "URL не найден",
			mock: func(mock *mocks.MockURLUseCase) {
				mock.EXPECT().GetOriginalURL(gomock.Any(), urlShort1).Return("", constants.ErrorURLNotExist)
			},
			in: pb.URLExpandRequest_builder{
				Id: &urlShort1,
			}.Build(),
			want:    nil,
			wantErr: status.Errorf(codes.NotFound, `URL %s not found`, urlShort1),
		},
		{
			name: "URL ранее был удален",
			mock: func(mock *mocks.MockURLUseCase) {
				mock.EXPECT().GetOriginalURL(gomock.Any(), urlShort1).Return("", constants.ErrorURLAlreadyDeleted)
			},
			in: pb.URLExpandRequest_builder{
				Id: &urlShort1,
			}.Build(),
			want:    nil,
			wantErr: status.Errorf(codes.NotFound, `URL %s already deleted`, urlShort1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockUseCase := newTestService(t)
			tt.mock(mockUseCase)

			ctx := context.WithValue(context.Background(), constants.UserIDKey, UUID)

			got, gotErr := s.ExpandURL(ctx, tt.in)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpandURL() = %v, want %v", got, tt.want)
			}

			st, ok := status.FromError(gotErr)
			if !ok {
				t.Fatalf("expected grpc status error")
			}
			if st.Code() != status.Code(tt.wantErr) {
				t.Errorf("code = %v, want %v", st.Code(), status.Code(tt.wantErr))
			}
		})
	}
}

func TestShortenerServiceServer_ListUserURLs(t *testing.T) {
	fullUrlShort1 := toolkit.AddBaseURLToResponse(cfg.BaseURL, urlShort1)
	fullUrlShort2 := toolkit.AddBaseURLToResponse(cfg.BaseURL, urlShort2)

	tests := []struct {
		name    string
		mock    func(*mocks.MockURLUseCase)
		in      *emptypb.Empty
		want    *pb.UserURLsResponse
		wantErr error
	}{
		{
			name: "Список URL успешно получен",
			mock: func(mock *mocks.MockURLUseCase) {
				mock.EXPECT().GetAllURLs(gomock.Any(), UUID).Return(urlsOut, nil)
			},
			in: &emptypb.Empty{},
			want: pb.UserURLsResponse_builder{
				Url: []*pb.URLData{
					pb.URLData_builder{ShortUrl: &fullUrlShort1, OriginalUrl: &urlOriginal1}.Build(),
					pb.URLData_builder{ShortUrl: &fullUrlShort2, OriginalUrl: &urlOriginal2}.Build(),
				},
			}.Build(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockUseCase := newTestService(t)
			tt.mock(mockUseCase)

			ctx := context.WithValue(context.Background(), constants.UserIDKey, UUID)

			got, gotErr := s.ListUserURLs(ctx, tt.in)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUserURLs() = %v, want %v", got, tt.want)
			}

			st, ok := status.FromError(gotErr)
			if !ok {
				t.Fatalf("expected grpc status error")
			}
			if st.Code() != status.Code(tt.wantErr) {
				t.Errorf("code = %v, want %v", st.Code(), status.Code(tt.wantErr))
			}

		})
	}
}
