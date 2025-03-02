package mock

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"terraform-provider-infomaniak/internal/apis/kaas"
	"time"
)

type KaasObject interface {
	Key() string
	*kaas.Kaas | *kaas.InstancePool
}

var (
	mockedApiStatePath = path.Join(os.TempDir(), "terraform-provider-infomaniak-kaas")
	mockedApiState     = make(map[string][]byte)
)

func getFromCache[K KaasObject](key string) (K, error) {
	obj, found := mockedApiState[key]
	if !found {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	var buff = bytes.NewBuffer(obj)
	var result K
	err := gob.NewDecoder(buff).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func addToCache[K KaasObject](obj K) error {
	key := obj.Key()
	_, found := mockedApiState[key]
	if found {
		return fmt.Errorf("key '%s' already present, you can't create two of the same, think mark, think", key)
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

func updateCache[K KaasObject](obj K) error {
	key := obj.Key()
	cachedObject, found := mockedApiState[key]
	if !found {
		return fmt.Errorf("key '%s' not found", key)
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

func removeFromCache[K KaasObject](obj K) error {
	key := obj.Key()
	_, found := mockedApiState[key]
	if !found {
		return fmt.Errorf("key '%s' not found", key)
	}

	delete(mockedApiState, key)
	saveCache()
	return nil
}

func init() {
	// Gob register
	gob.Register(&kaas.Kaas{})
	gob.Register(&kaas.InstancePool{})

	// Check cache age
	stat, err := os.Stat(mockedApiStatePath)
	if err == nil {
		// DeleteKaas cache if old
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

	// Create Kaas tmp file for caching
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
	os.WriteFile(mockedApiStatePath, data, 0666)
}

func ResetCache() {
	mockedApiState = make(map[string][]byte)
}
