package credential

import (
	"fmt"

	"github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker-credential-helpers/credentials"
)

type CredentialStore struct {
	nativeStoreWin                client.ProgramFunc
	nativeStoreLinuxPass          client.ProgramFunc
	nativeStoreLinuxSecretService client.ProgramFunc
	nativeStoreMacOSX             client.ProgramFunc
}

func NewCredentialStore() *CredentialStore {
	return &CredentialStore{
		nativeStoreWin:                client.NewShellProgramFunc("docker-credential-wincred.exe"),
		nativeStoreLinuxPass:          client.NewShellProgramFunc("docker-credential-pass"),
		nativeStoreLinuxSecretService: client.NewShellProgramFunc("docker-credential-secretservice"),
		nativeStoreMacOSX:             client.NewShellProgramFunc("docker-credential-osxkeychain"),
	}
}

func (cs *CredentialStore) tryStore(c *credentials.Credentials) error {
	nativeStores := []client.ProgramFunc{cs.nativeStoreWin, cs.nativeStoreLinuxPass, cs.nativeStoreLinuxSecretService, cs.nativeStoreMacOSX}
	for _, nativeStore := range nativeStores {
		err := client.Store(nativeStore, c)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("No fitting native credential manage found")
}

func (cs *CredentialStore) Store(c *credentials.Credentials) error {
	err := cs.tryStore(c)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CredentialStore) tryGet(url string) (*credentials.Credentials, error) {
	nativeStores := []client.ProgramFunc{cs.nativeStoreWin, cs.nativeStoreLinuxPass, cs.nativeStoreLinuxSecretService, cs.nativeStoreMacOSX}
	for _, nativeStore := range nativeStores {
		c, err := client.Get(nativeStore, url)
		if err == nil {
			return c, nil
		}
	}
	return nil, fmt.Errorf("No fitting native credential manage found")
}

func (cs *CredentialStore) Get(url string) (*credentials.Credentials, error) {
	c, err := cs.tryGet(url)
	return c, err
}
