package mock

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"os"
	"path"
	"terraform-provider-infomaniak/internal/apis/publiccloud"
	"time"
)

// PublicCloudObject is the constraint satisfied by every model persisted in
// the mock cache.
type PublicCloudObject interface {
	Key() string
	*publiccloud.PublicCloud | *publiccloud.Project | *publiccloud.User
}

var (
	mockedApiStatePath = path.Join(os.TempDir(), "terraform-provider-infomaniak-publiccloud")
	mockedApiState     = make(map[string][]byte)

	ErrKeyNotFound  = errors.New("key not found")
	ErrDuplicateKey = errors.New("duplicate key found")
)

func getFromCache[K PublicCloudObject](key string) (K, error) {
	var zero K
	obj, found := mockedApiState[key]
	if !found {
		return zero, ErrKeyNotFound
	}

	buff := bytes.NewBuffer(obj)
	var result K
	if err := gob.NewDecoder(buff).Decode(&result); err != nil {
		return zero, err
	}
	return result, nil
}

func addToCache[K PublicCloudObject](obj K) error {
	key := obj.Key()
	if _, found := mockedApiState[key]; found {
		return ErrDuplicateKey
	}

	var buff bytes.Buffer
	if err := gob.NewEncoder(&buff).Encode(obj); err != nil {
		return err
	}

	mockedApiState[key] = buff.Bytes()
	saveCache()
	return nil
}

func updateCache[K PublicCloudObject](obj K) error {
	key := obj.Key()
	if _, found := mockedApiState[key]; !found {
		return ErrKeyNotFound
	}

	var buff bytes.Buffer
	if err := gob.NewEncoder(&buff).Encode(obj); err != nil {
		return err
	}

	mockedApiState[key] = buff.Bytes()
	saveCache()
	return nil
}

func removeFromCache[K PublicCloudObject](obj K) error {
	key := obj.Key()
	if _, found := mockedApiState[key]; !found {
		return ErrKeyNotFound
	}

	delete(mockedApiState, key)
	saveCache()
	return nil
}

func init() {
	gob.Register(&publiccloud.PublicCloud{})
	gob.Register(&publiccloud.Project{})
	gob.Register(&publiccloud.User{})

	stat, err := os.Stat(mockedApiStatePath)
	if err == nil {
		if time.Since(stat.ModTime()) > 24*time.Hour {
			os.Remove(mockedApiStatePath)
			return
		}
	}

	bdy, err := os.ReadFile(mockedApiStatePath)
	if err == nil {
		if err := json.Unmarshal(bdy, &mockedApiState); err != nil {
			os.Remove(mockedApiStatePath)
		}
		return
	}

	if _, err := os.Create(mockedApiStatePath); err != nil {
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

// ResetCache wipes the in-memory cache. Tests call it from TestMain to keep
// runs deterministic.
func ResetCache() {
	mockedApiState = make(map[string][]byte)
}
