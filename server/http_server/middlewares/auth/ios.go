package auth

import (
	"net/http"
	"scheduler0/server/service"
	"scheduler0/utils"
)

// IsIOSClient returns true is the request is coming from an ios app
func IsIOSClient(req *http.Request) bool {
	apiKey := req.Header.Get(APIKeyHeader)
	bundleID := req.Header.Get(IOSBundleHeader)
	return  len(apiKey) > 9 && len(bundleID) > 9
}

// IsAuthorizedIOSClient returns true if the credential is authorized ios app
func IsAuthorizedIOSClient(req *http.Request, pool *utils.Pool) (bool, *utils.GenericError) {
	apiKey := req.Header.Get(APIKeyHeader)
	IOSBundleID := req.Header.Get(IOSBundleHeader)

	credentialService := service.Credential{
		Pool: pool,
	}

	return credentialService.ValidateIOSAPIKey(apiKey, IOSBundleID)
}