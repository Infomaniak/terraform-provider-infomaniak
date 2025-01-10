package mock

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
)

func genId() string {
	var b = make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &out)
	_, err = enc.Write(b)
	if err != nil {
		panic(err)
	}

	return out.String()
}

func genKubeconfig() string {
	var b = make([]byte, 1024)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &out)
	_, err = enc.Write(b)
	if err != nil {
		panic(err)
	}

	return out.String()
}
