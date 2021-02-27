package service

import (
	"net/http"
	"scheduler0/server/managers/credential"
	"scheduler0/server/transformers"
	"scheduler0/utils"
)

type Credential Service

func (credentialService *Credential) CreateNewCredential(credentialTransformer transformers.Credential) (string, *utils.GenericError) {
	credentialManager := credentialTransformer.ToManager()
	return credentialManager.CreateOne(credentialService.Pool)
}

func (credentialService *Credential) FindOneCredentialByUUID(UUID string) (*transformers.Credential, error) {
	credentialDto := transformers.Credential{UUID: UUID}
	credentialManager := credentialDto.ToManager()
	if err := credentialManager.GetOne(credentialService.Pool); err != nil {
		return nil, err
	} else {
		outboundDto := transformers.Credential{}
		outboundDto.FromManager(credentialManager)
		return &outboundDto, nil
	}
}

func (credentialService *Credential) UpdateOneCredential(credentialTransformer transformers.Credential) (*transformers.Credential, error) {
	credentialManager := credentialTransformer.ToManager()
	if _, err := credentialManager.UpdateOne(credentialService.Pool); err != nil {
		return nil, err
	} else {
		outboundDto := transformers.Credential{}
		outboundDto.FromManager(credentialManager)
		return &outboundDto, nil
	}
}

func (credentialService *Credential) DeleteOneCredential(UUID string) (*transformers.Credential, error) {
	credentialDto := transformers.Credential{UUID: UUID}
	credentialManager := credentialDto.ToManager()
	if _, err := credentialManager.DeleteOne(credentialService.Pool); err != nil {
		return nil, err
	} else {
		outboundDto := transformers.Credential{}
		outboundDto.FromManager(credentialManager)
		return &outboundDto, nil
	}
}

func (credentialService *Credential) ListCredentials(offset int, limit int, orderBy string) (*transformers.PaginatedCredential, *utils.GenericError) {
	credentialManager := credential.Manager{}

	total, err := credentialManager.Count(credentialService.Pool)
	if err != nil {
		return nil, err
	}

	if total < 1 {
		return  nil, utils.HTTPGenericError(http.StatusNotFound, "there no credentials")
	}

	if credentialManagers, err := credentialManager.GetAll(credentialService.Pool, offset, limit, orderBy); err != nil {
		return nil, err
	} else {
		credentialTransformers := make([]transformers.Credential, len(credentialManagers))

		for i, credentialManager := range credentialManagers {
			credentialTransformers[i].FromManager(credentialManager)
		}

		return &transformers.PaginatedCredential{
			Data: credentialTransformers,
			Total: total,
			Offset: offset,
			Limit: limit,
		}, nil
	}
}

func (credentialService *Credential) ValidateServerAPIKey(apiKey string, apiSecret string) (bool, *utils.GenericError) {
	credentialManager := credential.Manager{
		ApiKey: apiKey,
	}

	getApIError := credentialManager.GetByAPIKey(credentialService.Pool)
	if getApIError != nil {
		return false, getApIError
	}

	return credentialManager.ApiSecret == apiSecret, nil
}

func (credentialService *Credential) ValidateIOSAPIKey(apiKey string, IOSBundle string) (bool, *utils.GenericError) {
	credentialManager := credential.Manager{
		ApiKey: apiKey,
	}

	getApIError := credentialManager.GetByAPIKey(credentialService.Pool)
	if getApIError != nil {
		return false, getApIError
	}

	return credentialManager.IOSBundleIDRestriction == IOSBundle, nil
}

func (credentialService *Credential) ValidateAndroidAPIKey(apiKey string, androidPackageName string) (bool, *utils.GenericError) {
	credentialManager := credential.Manager{
		ApiKey: apiKey,
	}

	getApIError := credentialManager.GetByAPIKey(credentialService.Pool)
	if getApIError != nil {
		return false, getApIError
	}

	return credentialManager.AndroidPackageNameRestriction == androidPackageName, nil
}

func (credentialService *Credential) ValidateWebAPIKeyHTTPReferrerRestriction(apiKey string, callerUrl string) (bool, *utils.GenericError) {
	credentialManager := credential.Manager{
		ApiKey: apiKey,
	}

	getApIError := credentialManager.GetByAPIKey(credentialService.Pool)
	if getApIError != nil {
		return false, getApIError
	}

	if callerUrl == credentialManager.HTTPReferrerRestriction {
		return true, nil
	}

	return false, utils.HTTPGenericError(http.StatusUnauthorized, "the user is not authorized")
}