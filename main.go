package main

import (
	"log"
	"net"
	"net/http"

	"github.com/akkahshh24/go-simple-bank-app/api"
	db "github.com/akkahshh24/go-simple-bank-app/db/sqlc"
	"github.com/akkahshh24/go-simple-bank-app/gapi"
	"github.com/akkahshh24/go-simple-bank-app/grpc"
	"github.com/akkahshh24/go-simple-bank-app/util"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(connPool)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	// Create a new grpc server
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	// Use field names as given in the proto file for JSON serialization/deserialization
	// camelCase to snake_case conversion
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
