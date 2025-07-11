package mockbuilder

import "testing"

type BuilderContainer struct {
	AccountRepoBuilder *mockAccountRepoBuilder
	SessionRepoBuilder *mockSessionRepoBuilder
	HasherBuilder      *mockHasherBuilder
	CacheBuilder       *mockCacheBuilder
	AccountVerifier    *mockAccountVerifierBuilder
	TokenSigner        *mockSignerBuilder
}

func NewBuilderContainer(t *testing.T) *BuilderContainer {
	return &BuilderContainer{
		AccountRepoBuilder: newMockAccountRepoBuilder(t),
		SessionRepoBuilder: newMockSessionRepoBuilder(t),
		HasherBuilder:      newMockHasherBuilder(t),
		CacheBuilder:       newMockCacheBuilder(t),
		AccountVerifier:    newMockAccountVerifierBuilder(t),
		TokenSigner:        NewMockSignerBuilder(t),
	}
}
