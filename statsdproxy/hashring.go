// functions to build and maintain a hashring
package statsdproxy

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
)

const (
	MAXRINGSIZE = 50
)

type HashRingID uint32
type HashRing []StatsDBackend

func NewHashRing() *HashRing {
	ret := make(HashRing, 0, MAXRINGSIZE)
	return &ret
}

// add a new server instance into the hashring
// accepts an instance of StatsDBackend as parameter
// returns the sorted hashring with the newly appended server and error
func (ring *HashRing) Add(backend StatsDBackend) (HashRing, error) {
	if !backend.Alive() {
		err_msg := fmt.Sprintf("Backend %s:%s doesn't seem to be alive.", backend.Host,
			backend.Port)
		return *ring, errors.New(err_msg)
	}
	new_ring := append(*ring, backend)
	sort.Sort(ByHashRingID(new_ring))
	return new_ring, nil
}

// simple function to get a position in a hashring. The logic is ripped from
// libketama and uses MD5 to determine the position
//
// accepts a string to hash
//
// returns a HashRingID and error
func GetHashRingPosition(data string) (HashRingID, error) {
	h := md5.New()
	_, err := io.WriteString(h, data)
	if err != nil {
		log.Printf("Error creating hash for %s", data)
		return 0, err
	}
	digest := h.Sum(nil)

	result1 := uint32(digest[3]) << 24
	result2 := uint32(digest[2]) << 16
	result3 := uint32(digest[1]) << 8
	result4 := uint32(digest[0])

	id := result1 | result2 | result3 | result4
	if DebugMode {
		log.Printf("HashRingID for %s is %d", data, id)
	}

	return HashRingID(id), nil
}

// get the StatsDBackend instance that is responsible for a metric.
// Responsible in this case means either the ID of the metric is lower than
// the Ring ID of the backend. Or the first backend if the metric ID is higher
// than all backend IDs. This case is the ring wrap around part of the hash
// ring.
// accepts a metric name as a string as parameter
// returns a pointer to a StatsDBackend instance and error
func (ring *HashRing) GetBackendForMetric(name string) (*StatsDBackend, error) {
	if len(*ring) == 0 {
		return nil, errors.New("No backends in the hashring.")
	}
	backend := (*ring)[0]

	metric_id, err := GetHashRingPosition(name)
	if err != nil {
		msg := fmt.Sprintf("Unable to get hashring position for %s", name)
		return nil, errors.New(msg)
	}
	if DebugMode {
		log.Printf("Choosing backend from %v", *ring)
	}
	for _, possible_backend := range *ring {
		if possible_backend.Alive() && metric_id < possible_backend.RingID {
			// we only set the backend if it has a higher RingID and is alive
			backend = possible_backend
		}

	}
	if DebugMode {
		log.Printf("Backend for %s is %d", name, backend.Port)
	}

	return &backend, nil
}

// implement the sort interface so StatsDBackend instances are sortable by
// hashring ID
type ByHashRingID HashRing

func (a ByHashRingID) Len() int           { return len(a) }
func (a ByHashRingID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHashRingID) Less(i, j int) bool { return a[i].RingID < a[j].RingID }
