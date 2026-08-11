package grpc

import (
	"context"

	pb "github.com/xerox-1315/TravelShare-Email/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailClient struct {
	client pb.EmailServiceClient
}

func NewEmailClient(addr string) (*EmailClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &EmailClient{client: pb.NewEmailServiceClient(conn)}, nil
}

// отправить код верификации
func (e *EmailClient) SendVerificationCode(email string) (bool, error) {
	resp, err := e.client.SendVerificationCode(context.Background(), &pb.SendCodeRequest{
		Email: email,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

// проверить код верификации
func (e *EmailClient) VerifyCode(email, code string) (bool, error) {
	resp, err := e.client.VerifyCode(context.Background(), &pb.VerifyCodeRequest{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return false, err
	}
	return resp.Verify, nil
}
