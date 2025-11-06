package mock

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/base64"
	"math/rand/v2"
)

func genId() int64 {
	return rand.Int64()
}

func genKubeconfig() string {
	var b = make([]byte, 1024)
	_, err := cryptorand.Read(b)
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
