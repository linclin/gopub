package gojenkins

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

//CredentialsManager is utility to control credential plugin
//Credentials declared by it can be used in jenkins jobs
type CredentialsManager struct {
	J      *Jenkins
	Folder string
}

const baseFolderPrefix = "/job/%s"
const baseCredentialsURL = "%s/credentials/store/%s/domain/%s/"
const createCredentialsURL = baseCredentialsURL + "createCredentials"
const deleteCredentialURL = baseCredentialsURL + "credential/%s/doDelete"
const configCredentialURL = baseCredentialsURL + "credential/%s/config.xml"
const credentialsListURL = baseCredentialsURL + "api/json"

var listQuery = map[string]string{
	"tree": "credentials[id]",
}

//ClassUsernameCredentials is name if java class which implements credentials that store username-password pair
const ClassUsernameCredentials = "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"

type credentialID struct {
	ID string `json:"id"`
}

type credentialIDs struct {
	Credentials []credentialID `json:"credentials"`
}

//UsernameCredentials struct representing credential for storing username-password pair
type UsernameCredentials struct {
	XMLName     xml.Name `xml:"com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"`
	ID          string   `xml:"id"`
	Scope       string   `xml:"scope"`
	Description string   `xml:"description"`
	Username    string   `xml:"username"`
	Password    string   `xml:"password"`
}

//StringCredentials store only secret text
type StringCredentials struct {
	XMLName     xml.Name `xml:"org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"`
	ID          string   `xml:"id"`
	Scope       string   `xml:"scope"`
	Description string   `xml:"description"`
	Secret      string   `xml:"secret"`
}

//SSHCredentials store credentials for ssh keys.
type SSHCredentials struct {
	XMLName          xml.Name    `xml:"com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey"`
	ID               string      `xml:"id"`
	Scope            string      `xml:"scope"`
	Username         string      `xml:"username"`
	Description      string      `xml:"description,omitempty"`
	PrivateKeySource interface{} `xml:"privateKeySource"`
	Passphrase       string      `xml:"passphrase,omitempty"`
}

//DockerServerCredentials store credentials for docker keys.
type DockerServerCredentials struct {
	XMLName             xml.Name `xml:"org.jenkinsci.plugins.docker.commons.credentials.DockerServerCredentials"`
	ID                  string   `xml:"id"`
	Scope               string   `xml:"scope"`
	Username            string   `xml:"username"`
	Description         string   `xml:"description,omitempty"`
	ClientKey           string   `xml:"clientKey"`
	ClientCertificate   string   `xml:"clientCertificate"`
	ServerCaCertificate string   `xml:"serverCaCertificate"`
}

//KeySourceDirectEntryType is used when secret in provided directly as private key value
const KeySourceDirectEntryType = "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey$DirectEntryPrivateKeySource"

//KeySourceOnMasterType is used when private key value is path to file on jenkins master
const KeySourceOnMasterType = "com.cloudbees.jenkins.plugins.sshcredentials.impl.BasicSSHUserPrivateKey$FileOnMasterPrivateKeySource"

//PrivateKey used in SSHCredentials type, type can be either:
//KeySourceDirectEntryType - then value should be text with secret
//KeySourceOnMasterType - then value should be path on master jenkins where secret is stored
type PrivateKey struct {
	Value string `xml:"privateKey"`
	Class string `xml:"class,attr"`
}

type PrivateKeyFile struct {
	Value string `xml:"privateKeyFile"`
	Class string `xml:"class,attr"`
}

func (cm CredentialsManager) fillURL(url string, params ...interface{}) string {
	var args []interface{}
	if cm.Folder != "" {
		args = []interface{}{fmt.Sprintf(baseFolderPrefix, cm.Folder), "folder"}
	} else {
		args = []interface{}{"", "system"}
	}
	return fmt.Sprintf(url, append(args, params...)...)
}

//List ids if credentials stored inside provided domain
func (cm CredentialsManager) List(domain string) ([]string, error) {

	idsResponse := credentialIDs{}
	ids := make([]string, 0)
	err := cm.handleResponse(cm.J.Requester.Get(cm.fillURL(credentialsListURL, domain), &idsResponse, listQuery))
	if err != nil {
		return ids, err
	}

	for _, id := range idsResponse.Credentials {
		ids = append(ids, id.ID)
	}

	return ids, nil
}

//GetSingle searches for credential in given domain with given id, if credential is found
//it will be parsed as xml to creds parameter(creds must be pointer to struct)
func (cm CredentialsManager) GetSingle(domain string, id string, creds interface{}) error {
	str := ""
	err := cm.handleResponse(cm.J.Requester.Get(cm.fillURL(configCredentialURL, domain, id), &str, map[string]string{}))
	if err != nil {
		return err
	}

	return xml.Unmarshal([]byte(str), &creds)
}

//Add credential to given domain, creds must be struct which is parsable to xml
func (cm CredentialsManager) Add(domain string, creds interface{}) error {
	return cm.postCredsXML(cm.fillURL(createCredentialsURL, domain), creds)
}

//Delete credential in given domain with given id
func (cm CredentialsManager) Delete(domain string, id string) error {
	return cm.handleResponse(cm.J.Requester.Post(cm.fillURL(deleteCredentialURL, domain, id), nil, cm.J.Raw, map[string]string{}))
}

//Update credential in given domain with given id, creds must be pointer to struct which is parsable to xml
func (cm CredentialsManager) Update(domain string, id string, creds interface{}) error {
	return cm.postCredsXML(cm.fillURL(configCredentialURL, domain, id), creds)
}

func (cm CredentialsManager) postCredsXML(url string, creds interface{}) error {
	payload, err := xml.Marshal(creds)
	if err != nil {
		return err
	}

	return cm.handleResponse(cm.J.Requester.PostXML(url, string(payload), cm.J.Raw, map[string]string{}))
}

func (cm CredentialsManager) handleResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode == 409 {
		return fmt.Errorf("Resource already exists, conflict status returned")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	return nil
}
