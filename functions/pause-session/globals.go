package pause_session

import "os"

var projectID string
var databaseName string

func init() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		panic("GOOGLE_CLOUD_PROJECT not set in environment")
	}

	databaseName = os.Getenv("FIRESTORE_DB")
	if databaseName == "" {
		panic("FIRESTORE_DB not set in environment")
	}
}
