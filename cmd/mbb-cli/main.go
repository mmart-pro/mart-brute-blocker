package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/mmart-pro/mart-brute-blocker/internal/config"
	"github.com/mmart-pro/mart-brute-blocker/internal/grpc/pb"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
)

var (
	configFlag  string
	hostFlag    string
	portFlag    int
	versionFlag bool

	rootCmdP = &cobra.Command{
		Use:   "mbb-cli",
		Short: "Клиент коммандной строки для mbb-service",
		Long:  `Клиент коммандной строки для управления mart brute blocker service`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(*cobra.Command, []string) {
			if versionFlag {
				printVersion()
			}
		},
	}

	allowCmdP = &cobra.Command{
		Use:     "allow ip/bitmask",
		Example: fmt.Sprintf("  %v allow 192.168.1.0/24\n  %v allow 172.16.1.1/32", rootCmdP.Use, rootCmdP.Use),
		Short:   "Добавляет подсеть в белый список",
		Long:    "Команда allow добавляет подсеть в белый список. Возвращает ошибку, если подсеть пересекается с уже существующией в списке.",
		Args:    cobra.ExactArgs(1),
	}

	denyCmdP = &cobra.Command{
		Use:     "deny ip/bitmask",
		Example: fmt.Sprintf("  %v deny 192.168.1.0/24\n  %v deny 172.16.1.1/32", rootCmdP.Use, rootCmdP.Use),
		Short:   "Добавляет подсеть в чёрный список",
		Long:    "Команда deny добавляет подсеть в чёрный список. Возвращает ошибку, если подсеть пересекается с уже существующией в списке.",
		Args:    cobra.ExactArgs(1),
	}

	removeCmdP = &cobra.Command{
		Use:     "remove ip/bitmask",
		Example: fmt.Sprintf("  %v remove 192.168.1.0/24\n  %v remove 172.16.1.1/32", rootCmdP.Use, rootCmdP.Use),
		Short:   "Удаляет подсеть из чёрного или белого списка по точному совпадению.",
		Long: "Команда remove удаляет подсеть из чёрного или белого списка по точному совпадению. " +
			" Возвращает ошибку, если подсеть не найена в списках.",
		Args: cobra.ExactArgs(1),
	}

	checkCmdP = &cobra.Command{
		Use:     "check ip",
		Example: fmt.Sprintf("  %v check 172.16.1.1", rootCmdP.Use),
		Short:   "Проверяет включение ip в один из списков.",
		Long:    "Проверяет включение ip в один из списков. Возвращает: BlackList | WhiteList | NotInList.",
		Args:    cobra.ExactArgs(1),
	}

	clearbucketCmdP = &cobra.Command{
		Use:     "clearbucket ip login",
		Example: fmt.Sprintf("  %v clearbucket 172.16.1.1 ivan@petrov", rootCmdP.Use),
		Short:   "Обнуляет лимит по указанным ip и login",
		Long:    "Обнуляет лимит по указанным ip и login",
		Args:    cobra.ExactArgs(2),
	}
)

func main() {
	rootCmdP.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "версия приложения")
	rootCmdP.PersistentFlags().StringVarP(&configFlag, "config", "c", "", "json-файл конфигурации")
	rootCmdP.PersistentFlags().StringVar(&hostFlag, "host", "", "адрес сервера")
	rootCmdP.PersistentFlags().IntVar(&portFlag, "port", 0, "порт сервера")

	rootCmdP.AddCommand(allowCmdP, denyCmdP, removeCmdP, checkCmdP, clearbucketCmdP)

	for _, c := range rootCmdP.Commands() {
		c.Run = runner
	}

	err := rootCmdP.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func runner(cmd *cobra.Command, args []string) {
	// с какими параметрами работаем
	runtimeConfig := config.CliConfig{}
	if configFlag != "" {
		cfg, err := config.NewCliConfig(configFlag)
		if err != nil {
			log.Fatal(fmt.Errorf("error read config from %s: %w", configFlag, err))
		}
		runtimeConfig = cfg
	}

	if hostFlag != "" {
		runtimeConfig.GrpcConfig.GrpcHost = hostFlag
	}
	if portFlag != 0 {
		runtimeConfig.GrpcConfig.GrpcPort = strconv.Itoa(portFlag)
	}

	addr := net.JoinHostPort(runtimeConfig.GrpcConfig.GrpcHost, runtimeConfig.GrpcConfig.GrpcPort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("couldn't connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMBBServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// run command
	switch cmd {
	case allowCmdP:
		_, err = client.Allow(ctx, &pb.SubnetReq{Subnet: args[0]})
	case denyCmdP:
		_, err = client.Deny(ctx, &pb.SubnetReq{Subnet: args[0]})
	case removeCmdP:
		_, err = client.Remove(ctx, &pb.SubnetReq{Subnet: args[0]})
	case clearbucketCmdP:
		_, err = client.ClearBucket(ctx, &pb.ClearBucketRequest{Ip: args[0], Login: args[0]})
	case checkCmdP:
		var res *pb.ContainsResponse
		res, err = client.Contains(ctx, &pb.IpReq{Ip: args[0]})
		if err == nil {
			fmt.Printf("ip contains in %s\n", model.ListType(res.ListType).String())
		}
	default:
		err = errors.New("unknown command")
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}
