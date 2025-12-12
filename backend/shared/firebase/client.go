package firebase

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

// InitFirestore initializes the Firestore client using Application Default Credentials
func InitFirestore(ctx context.Context) error {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Println("Warning: GCP_PROJECT_ID not set, attempting to use Application Default Credentials")
	}

	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	if databaseID == "" {
		databaseID = "debates" // default database name
	}

	conf := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	// Use DatabaseID option when creating Firestore client
	client, err := app.Firestore(ctx, option.WithDatabaseID(databaseID))
	if err != nil {
		return err
	}

	firestoreClient = client
	log.Printf("Firestore client initialized successfully for project: %s, database: %s", projectID, databaseID)
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
