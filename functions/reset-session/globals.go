package reset_session

import "os"

var projectID string

func init() {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		panic("GOOGLE_CLOUD_PROJECT not set in environment")
	}
}
