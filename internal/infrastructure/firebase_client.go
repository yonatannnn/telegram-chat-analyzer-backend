
package infrastructure

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func SetupFirebase(credsPath string) (*firebase.App, error) {
	return firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(credsPath))
}
