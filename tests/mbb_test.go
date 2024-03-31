package tests

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/mmart-pro/mart-brute-blocker/internal/grpc/pb"
	"github.com/mmart-pro/mart-brute-blocker/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client  pb.MBBServiceClient
	ctx     = context.Background()
	stepErr error
)

func FeatureContext(s *godog.ScenarioContext) {
	s.Step(`^я хочу добавить подсеть "([^"]*)" в White list$`, iWantToAddIPAddressToWhiteList)
	s.Step(`^я вызываю функцию Allow с подсетью "([^"]*)"$`, iCallFunctionAllowWithIPAddress)
	s.Step(`^я вызываю функцию Allow с ещё одной подсетью "([^"]*)"$`, iCallFunctionAllowWithIPAddress2)
	s.Step(`^IP-адрес "([^"]*)" должен быть добавлен в White list$`, ipAddressShouldBeAddedToWhiteList)
	s.Step(`^я хочу добавить подсеть "([^"]*)" в Black list$`, iWantToAddIPAddressToBlackList)
	s.Step(`^я вызываю функцию Deny с подсетью "([^"]*)"$`, iCallFunctionDenyWithIPAddress)
	s.Step(`^я вызываю функцию Deny с ещё одной подсетью "([^"]*)"$`, iCallFunctionDenyWithIPAddress2)
	s.Step(`^IP-адрес "([^"]*)" должен быть добавлен в Black list$`, ipAddressShouldBeAddedToBlackList)

	s.Step(`^подсеть "([^"]*)" находится в White list$`, subnetIsInWhiteList)
	s.Step(`^подсеть "([^"]*)" находится в Black list$`, subnetIsInBlackList)
	s.Step(`^я вызываю функцию Remove с подсетью "([^"]*)"$`, iCallFunctionRemoveWithSubnet)
	s.Step(`^IP-адрес "([^"]*)" должен быть удален из списков$`, ipaddrShouldBeRemovedFromLists)
	s.Step(`^я вызываю функцию Contains с IP-адресом "([^"]*)"$`, iCallFunctionContainsWithIPAddress)
	s.Step(`^я должен получить ответ, что IP-адрес "([^"]*)" находится в White list$`, iShouldGetResponseThatIPAddressIsInWhiteList)
	s.Step(`^я должен получить ответ, что IP-адрес "([^"]*)" находится в Black list$`, iShouldGetResponseThatIPAddressIsInBlackList)

	s.Step(`^я получаю ошибку$`, iGetAnError)
}

func iGetAnError() error {
	if stepErr != nil {
		return nil //nolint:nilerr
	}
	return errors.New("error expected")
}

func iWantToAddIPAddressToWhiteList(_ string) error {
	return nil
}

func iWantToAddIPAddressToBlackList(_ string) error {
	return nil
}

func iCallFunctionAllowWithIPAddress(ip string) error {
	_, err := client.Allow(ctx, &pb.SubnetReq{Subnet: ip})
	return err
}

func iCallFunctionAllowWithIPAddress2(ip string) error {
	_, stepErr = client.Allow(ctx, &pb.SubnetReq{Subnet: ip})
	return nil
}

func iCallFunctionDenyWithIPAddress(ip string) error {
	_, err := client.Deny(ctx, &pb.SubnetReq{Subnet: ip})
	return err
}

func iCallFunctionDenyWithIPAddress2(ip string) error {
	_, stepErr := client.Deny(ctx, &pb.SubnetReq{Subnet: ip})
	return stepErr
}

func ipAddressShouldBeAddedToWhiteList(ip string) error {
	res, err := client.Contains(ctx, &pb.IpReq{Ip: ip})
	if err != nil {
		return err
	}
	if model.ListType(res.ListType) != model.WhiteList {
		return errors.New("IP address is not in the White list")
	}
	return nil
}

func ipAddressShouldBeAddedToBlackList(ip string) error {
	res, err := client.Contains(ctx, &pb.IpReq{Ip: ip})
	if err != nil {
		return err
	}
	if model.ListType(res.ListType) != model.BlackList {
		return errors.New("IP address is not in the Black list")
	}
	return nil
}

func subnetIsInWhiteList(subnet string) error {
	_, err := client.Allow(ctx, &pb.SubnetReq{Subnet: subnet})
	return err
}

func subnetIsInBlackList(subnet string) error {
	_, err := client.Deny(ctx, &pb.SubnetReq{Subnet: subnet})
	return err
}

func iCallFunctionRemoveWithSubnet(subnet string) error {
	_, err := client.Remove(ctx, &pb.SubnetReq{Subnet: subnet})
	return err
}

func ipaddrShouldBeRemovedFromLists(ip string) error {
	res, err := client.Contains(ctx, &pb.IpReq{Ip: ip})
	if err != nil {
		return err
	}

	if model.ListType(res.ListType) != model.NotInList {
		return errors.New("IP addr is still in the list")
	}

	return nil
}

func iCallFunctionContainsWithIPAddress(ip string) error {
	_, err := client.Contains(ctx, &pb.IpReq{Ip: ip})
	return err
}

func iShouldGetResponseThatIPAddressIsInWhiteList(ip string) error {
	res, err := client.Contains(ctx, &pb.IpReq{Ip: ip})
	if err != nil {
		return err
	}

	if model.ListType(res.ListType) != model.WhiteList {
		return errors.New("IP address is not in the White list")
	}

	return nil
}

func iShouldGetResponseThatIPAddressIsInBlackList(ip string) error {
	res, err := client.Contains(ctx, &pb.IpReq{Ip: ip})
	if err != nil {
		return err
	}

	if model.ListType(res.ListType) != model.BlackList {
		return errors.New("IP address is not in the Black list")
	}

	return nil
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty", // can define default values
}

func init() {
	godog.BindCommandLineFlags("", &opt)
}

func TestMain(_ *testing.M) {
	status := 0

	flag.Parse()
	opt.Paths = flag.Args()

	suite := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: initializeTestSuite,
		ScenarioInitializer:  FeatureContext,
		Options:              &opt,
	}

	status = suite.Run()

	os.Exit(status)
}

func initializeTestSuite(ctx *godog.TestSuiteContext) {
	var conn *grpc.ClientConn

	ctx.BeforeSuite(func() {
		var err error
		conn, err = grpc.Dial("mbb-api:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		// conn, err = grpc.Dial("localhost:15000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("couldn't connect: %v", err)
		}
		client = pb.NewMBBServiceClient(conn)
	})

	ctx.AfterSuite(func() {
		conn.Close()
	})
}
