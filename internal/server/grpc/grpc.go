package grpcserver

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	"github.com/Di-nis/shortener-url/internal/authn"
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/models"
	pb "github.com/Di-nis/shortener-url/internal/proto"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/toolkit"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// URLCreator - интерфейс, включащий методы по созданию URL.
type URLCreator interface {
	CreateURLOrdinary(context.Context, any) (models.URLBase, error)
}

// URLReader - интерфейс, включащий методы по получению URL.
type URLReader interface {
	GetOriginalURL(context.Context, string) (string, error)
	GetAllURLs(context.Context, string) ([]models.URLBase, error)
}

// URLUseCase - объединенный интерфейс.
type URLUseCase interface {
	URLCreator
	URLReader
}

// ShortenerServiceServer поддерживает все необходимые методы сервера.
type ShortenerServiceServer struct {
	pb.UnimplementedShortenerServiceServer
	URLCreator URLCreator
	URLReader  URLReader
	Config     *config.Config
}

// NewShortenerServiceServer - создание нового сервера.
func NewShortenerServiceServer(useCase URLUseCase, config *config.Config) *ShortenerServiceServer {
	return &ShortenerServiceServer{
		URLCreator: useCase,
		URLReader:  useCase,
		Config:     config,
	}
}

// ShortenURL - создание короткого URL.
func (s *ShortenerServiceServer) ShortenURL(ctx context.Context, in *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	var response pb.URLShortenResponse

	userID := ctx.Value(constants.UserIDKey).(string)
	urlOriginal := in.GetUrl()
	urlIn := models.URLJSON{
		UUID:     userID,
		Original: urlOriginal,
	}

	urlOut, err := s.URLCreator.CreateURLOrdinary(ctx, urlIn)
	if err != nil {
		if errors.Is(err, constants.ErrorURLAlreadyExist) {
			return nil, status.Errorf(codes.AlreadyExists, `URL %s already exist`, urlOriginal)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	urlOut.Short = toolkit.AddBaseURLToResponse(s.Config.BaseURL, urlOut.Short)

	response.SetResult(urlOut.Short)

	return &response, nil
}

// ExpandURL - получение оригинального URL.
func (s *ShortenerServiceServer) ExpandURL(ctx context.Context, in *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	var response pb.URLExpandResponse

	urlOriginal, err := s.URLReader.GetOriginalURL(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, constants.ErrorURLNotExist) {
			return nil, status.Errorf(codes.NotFound, `URL %s not found`, in.GetId())
		}
		if errors.Is(err, constants.ErrorURLAlreadyDeleted) {
			return nil, status.Errorf(codes.NotFound, `URL %s already deleted`, in.GetId())
		}
		return nil, status.Error(codes.Unavailable, "server unavailable")
	}

	response.SetResult(urlOriginal)

	return &response, nil
}

// ListUserURLs - получение списка всех коротких URL пользователя.
func (s *ShortenerServiceServer) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	var response pb.UserURLsResponse

	userID := ctx.Value(constants.UserIDKey).(string)

	urls, err := s.URLReader.GetAllURLs(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	var urlsOut []*pb.URLData
	for _, url := range urls {
		shortURL := toolkit.AddBaseURLToResponse(s.Config.BaseURL, url.Short)
		urlOut := pb.URLData_builder{
			ShortUrl:    &shortURL,
			OriginalUrl: &url.Original,
		}.Build()
		urlsOut = append(urlsOut, urlOut)
	}

	response.SetUrl(urlsOut)

	return &response, nil
}

// Run - запуск gRPC-сервера.
func Run(ctx context.Context, config *config.Config, repo usecase.URLRepository, svc *service.Service) error {
	listen, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		logger.Sugar.Errorf("failed initializing listener, error - %w", err)
		os.Exit(1)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(authn.Interceptor))

	useCase := usecase.NewURLUseCase(repo, svc)
	pb.RegisterShortenerServiceServer(server, NewShortenerServiceServer(useCase, config))

	logger.Sugar.Info("gRPC-server has started")

	go func() {
		if err := server.Serve(listen); err != nil {
			logger.Sugar.Errorf("gRPC-server failed, error - %w", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		server.GracefulStop()
		if err = repo.Close(); err != nil {
			logger.Sugar.Errorf("failed closing database: %w", err)
		}
	}()
	<-shutDownCtx.Done()

	return nil
}
