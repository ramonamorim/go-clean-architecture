package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ramonamorim/go-clean-architecture/configs"
	"github.com/ramonamorim/go-clean-architecture/internal/event/handler"
	"github.com/ramonamorim/go-clean-architecture/internal/infra/graph"
	"github.com/ramonamorim/go-clean-architecture/internal/infra/grpc/pb"
	"github.com/ramonamorim/go-clean-architecture/internal/infra/grpc/service"
	"github.com/ramonamorim/go-clean-architecture/internal/infra/web/webserver"
	"github.com/ramonamorim/go-clean-architecture/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddRoute("/order", http.MethodPost, webOrderHandler.Create)
	webserver.AddRoute("/order", http.MethodGet, webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase, ListOrdersUseCase: *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	var ch *amqp.Channel
	var err error

	for i := 0; i < 10; i++ {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ, retrying... (%d/10)", i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			panic(err)
		}
		return ch
	}

	panic(fmt.Sprintf("Failed to connect to RabbitMQ after 10 attempts: %v", err))
}
