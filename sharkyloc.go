package sharkyloc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	int64Size   = 8
	addressSize = 32
	encodedSize = addressSize + int64Size + int64Size
)

type SharkyLoc struct {
	mu sync.Mutex
	f  io.ReadWriter
	m  map[string]Loc
}

type Loc struct {
	Bucket int64
	Offset int64
}

type locmeta struct {
	Address [addressSize]byte
	Bucket  int64
	Offset  int64
}

func New(f io.ReadWriter) (*SharkyLoc, error) {

	l := &SharkyLoc{
		f: f,
		m: make(map[string]Loc),
	}

	err := l.load()
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *SharkyLoc) Read(addr swarm.Address) (Loc, error) {

	if m, ok := l.m[addr.ByteString()]; ok {
		return m, nil
	}

	return Loc{}, errors.New("not found")
}

func (l *SharkyLoc) Write(addr swarm.Address, bucket, offset int64) error {

	l.mu.Lock()
	defer l.mu.Unlock()

	l.m[addr.ByteString()] = Loc{Bucket: bucket, Offset: offset}

	var a [32]byte
	for i, b := range addr.Bytes() {
		a[i] = b
	}

	return binary.Write(l.f, binary.LittleEndian, locmeta{Address: a, Bucket: bucket, Offset: offset})
}

func (l *SharkyLoc) load() error {

	b, err := io.ReadAll(l.f)
	if err != nil {
		return err
	}

	mMeta := make([]locmeta, len(b)/encodedSize)
	err = binary.Read(bytes.NewBuffer(b), binary.LittleEndian, mMeta)
	if err != nil {
		return err
	}

	for _, m := range mMeta {
		l.m[string(m.Address[:])] = Loc{Offset: m.Offset, Bucket: m.Bucket}
	}

	return nil
}
