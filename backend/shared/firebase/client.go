package firebase

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var firestoreClient *firestore.Client

// InitFirestore initializes the Firestore client using Application Default Credentials
func InitFirestore(ctx context.Context) error {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Println("Warning: GCP_PROJECT_ID not set, attempting to use Application Default Credentials")
	}

	conf := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	firestoreClient = client
	log.Printf("Firestore client initialized successfully for project: %s", projectID)
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
