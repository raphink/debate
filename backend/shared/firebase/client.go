package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var firestoreClient *firestore.Client

// InitFirestore initializes the Firestore client using Application Default Credentials
func InitFirestore(ctx context.Context) error {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	firestoreClient = client
	log.Println("Firestore client initialized successfully")
	return nil
}

// GetClient returns the initialized Firestore client
func GetClient() *firestore.Client {
	return firestoreClient
}

// Close closes the Firestore client connection
func Close() error {
	if firestoreClient != nil {
		return firestoreClient.Close()
	}
	return nil
}
