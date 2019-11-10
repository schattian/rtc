package git

import "github.com/sebach1/git-crud/integrity"

type credentials map[integrity.SchemaName]string

func (creds *credentials) Encrypt(schName integrity.SchemaName, cred string) {
}

func (creds *credentials) Decrypt(schName integrity.SchemaName, cred string) {

}
