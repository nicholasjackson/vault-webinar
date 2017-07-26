package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	writeCertsFromVault()

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(rw, "Hello World")
	})

	err := http.ListenAndServeTLS(":8433", "cert.pem", "key.pem", nil)

	log.Fatal(err)
}

type vaultCertRequest struct {
	CommonName string `json:"common_name"`
}

type vaultCertResponse struct {
	Data vaultResponseData `json:"data"`
}

type vaultResponseData struct {
	Certificate string `json:"certificate"`
	CAChain     string `json:"ca_chain"`
	PrivateKey  string `json:"private_key"`
}

func writeCertsFromVault() {
	client := http.DefaultClient

	requestData := vaultCertRequest{
		CommonName: "localhost",
	}

	data, _ := json.Marshal(requestData)
	request, _ := http.NewRequest("POST", "http://localhost:8200/v1/pki/issue/localhost", bytes.NewReader(data))
	request.Header.Add("X-Vault-Token", os.Getenv("VAULT_TOKEN"))

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	vr := vaultCertResponse{}
	json.NewDecoder(response.Body).Decode(&vr)

	ioutil.WriteFile("cert.pem", []byte(vr.Data.CAChain), 0644)
	ioutil.WriteFile("key.pem", []byte(vr.Data.PrivateKey), 0644)
}
