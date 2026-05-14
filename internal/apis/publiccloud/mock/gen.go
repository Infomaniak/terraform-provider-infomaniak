package mock

import "math/rand/v2"

func genId() int64 {
	id := rand.Int64()
	if id < 0 {
		id = -id
	}
	if id == 0 {
		id = 1
	}
	return id
}
