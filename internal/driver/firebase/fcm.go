package firebase

import (
	"context"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var (
	fcmClient *messaging.Client
	fcmOnce   sync.Once
)

func FcmInit() (error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./firebase-admin-sdk.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	client, err := app.Messaging(context.TODO())
	if err != nil {
		return err
	}

	fcmOnce.Do(func() {
		fcmClient = client
	})
	return nil
}

func GetFcmClient() *messaging.Client {
	return fcmClient
}
