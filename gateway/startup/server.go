package startup

import (
	"context"
	"fmt"
	"gateway/infrastructure/api"
	"gateway/startup/config"
	userGateway "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/gateway"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgo "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	userGateway.UnimplementedUserGatewayServer
	tracer otgo.Tracer
	closer io.Closer
	Config *config.Config
}

const name = "gateway"

func NewServer(config *config.Config) (*Server, error) {
	tracer, closer := tracer.Init(name)
	otgo.SetGlobalTracer(tracer)
	server := &Server{
		tracer: tracer,
		closer: closer,
		Config: config,
	}

	return server, nil
}

func (server *Server) GetTracer() otgo.Tracer {
	return server.tracer
}

func (server *Server) GetCloser() io.Closer {
	return server.closer
}

func (server *Server) StartServer(userGatewayS *api.UserGatewayStruct) {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", server.Config.GrpcPort))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Attach the Greeter service to the server
	userGateway.RegisterUserGatewayServer(s, userGatewayS)
	// Serve gRPC server
	log.Println(fmt.Sprintf("Serving gRPC on localhost:%s", server.Config.GrpcPort))
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%s", server.Config.GrpcPort),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(otgo.GlobalTracer()),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_opentracing.StreamClientInterceptor(
				grpc_opentracing.WithTracer(otgo.GlobalTracer()),
			),
		),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = userGateway.RegisterUserGatewayHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.Config.HttpPort),
		Handler: tracer.TracingWrapper(gwmux),
	}

	log.Println(fmt.Sprintf("Serving gRPC-Gateway on http://localhost:%s", server.Config.HttpPort))
	log.Fatalln(gwServer.ListenAndServe())
}

func (server *Server) Start() {
	userGateway := server.initHandlers()
	server.StartServer(userGateway)
}

func (server *Server) initHandlers() *api.UserGatewayStruct {
	return api.NewUserGateway(server.Config)
}
