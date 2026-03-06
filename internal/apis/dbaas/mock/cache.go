package mock

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"os"
	"path"
	"terraform-provider-infomaniak/internal/apis/dbaas"
	"time"
)

type DBaasObject interface {
	Key() string
	*dbaas.DBaaS | *dbaas.DBaasBackupSchedule
}

var (
	mockedApiStatePath = path.Join(os.TempDir(), "terraform-provider-infomaniak-dbaas")
	mockedApiState     = make(map[string][]byte)

	ErrKeyNotFound  = errors.New("key not found")
	ErrDuplicateKey = errors.New("duplicate key found")
)

func getFromCache[K DBaasObject](key string) (K, error) {
	obj, found := mockedApiState[key]
	if !found {
		return nil, ErrKeyNotFound
	}

	var buff = bytes.NewBuffer(obj)
	var result K
	err := gob.NewDecoder(buff).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrKeyNotFound
	}

	return result, nil
}

func addToCache[K DBaasObject](obj K) error {
	key := obj.Key()
	_, found := mockedApiState[key]
	if found {
		return ErrDuplicateKey
	}

	var buff bytes.Buffer
	err := gob.NewEncoder(&buff).Encode(obj)
	if err != nil {
		return err
	}

	mockedApiState[key] = buff.Bytes()
	saveCache()
	return nil
}

func updateCache[K DBaasObject](obj K) error {
	key := obj.Key()
	cachedObject, found := mockedApiState[key]
	if !found {
		return ErrKeyNotFound
	}

	var buff = bytes.NewBuffer(cachedObject)
	var result K
	err := gob.NewDecoder(buff).Decode(&result)
	if err != nil {
		return err
	}

	var newBuff bytes.Buffer
	err = gob.NewEncoder(&newBuff).Encode(obj)
	if err != nil {
		return err
	}

	mockedApiState[key] = newBuff.Bytes()
	saveCache()
	return nil
}

func removeFromCache[K DBaasObject](obj K) error {
	key := obj.Key()
	_, found := mockedApiState[key]
	if !found {
		return ErrKeyNotFound
	}

	delete(mockedApiState, key)
	saveCache()
	return nil
}

func init() {
	// Gob register
	gob.Register(&dbaas.DBaaS{})

	// Check cache age
	stat, err := os.Stat(mockedApiStatePath)
	if err == nil {
		// Delete DBaaS cache if old
		if time.Since(stat.ModTime()) > 24*time.Hour {
			os.Remove(mockedApiStatePath)
			return
		}
	}

	// Try to get cache
	bdy, err := os.ReadFile(mockedApiStatePath)
	if err == nil {
		// Cache found
		err := json.Unmarshal(bdy, &mockedApiState)
		if err != nil {
			os.Remove(mockedApiStatePath)
		}
		return
	}

	// Create DBaaS tmp file for caching
	_, err = os.Create(mockedApiStatePath)
	if err != nil {
		panic(err)
	}
}

func saveCache() {
	data, err := json.Marshal(mockedApiState)
	if err != nil {
		return
	}
	//nolint:errcheck
	os.WriteFile(mockedApiStatePath, data, 0666)
}

func ResetCache() {
	mockedApiState = make(map[string][]byte)
}
