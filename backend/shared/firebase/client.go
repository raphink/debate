package firebase

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

var firestoreClient *firestore.Client

// InitFirestore initializes the Firestore client using Application Default Credentials
func InitFirestore(ctx context.Context) error {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("GCP_PROJECT_ID environment variable is required")
	}

	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	if databaseID == "" {
		databaseID = "debates" // default database name
	}

	// Construct database path: projects/{project}/databases/{database}
	databasePath := fmt.Sprintf("projects/%s/databases/%s", projectID, databaseID)

	// Create Firestore client directly with database path
	client, err := firestore.NewClient(ctx, projectID, firestore.DatabaseID(databaseID))
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}

	firestoreClient = client
	log.Printf("Firestore client initialized successfully for %s", databasePath)
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
