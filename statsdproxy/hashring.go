// functions to build and maintain a hashring
package statsdproxy

import (
  "io"
  "crypto/md5"
)

const (
  MAXRINGSIZE = 50
)

type HashRingID uint32
type HashRing []StatsDBackend

func New(size int) *HashRing {
  ret := make(HashRing, size, MAXRINGSIZE)
  return &ret
}

func (ring *HashRing) Add(backend StatsDBackend) (HashRing, error) {
  new_ring := append(*ring, backend)
  return new_ring, nil
}

// simple function to get a position in a hashring. The logic is ripped from
// libketama and uses MD5 to determine the position
func GetHashRingPosition(data string) (HashRingID, error) {
    h := md5.New()
    io.WriteString(h, data)
    digest := h.Sum(nil)

    result1 := uint32(digest[3]) << 24
    result2 := uint32(digest[2]) << 16
    result3 := uint32(digest[1]) << 8
    result4 := uint32(digest[0])

    return HashRingID(result1 | result2 | result3 | result4), nil
}

func GetBackendForMetric(name string) (*StatsDBackend, error) {
  return new(StatsDBackend), nil
}

