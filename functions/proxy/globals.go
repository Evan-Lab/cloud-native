package proxy

import "os"

func init() {
	projectID = os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		panic("GCP_PROJECT_ID not set in environment")
	}
}
