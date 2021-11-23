package sharkyloc_test

import (
	"bytes"
	sharkyloc "sharky"
	"testing"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/swarm/test"
)

func TestWriteRead(t *testing.T) {

	tc := []struct {
		addr   swarm.Address
		bucket int64
		offset int64
	}{
		{
			addr:   test.RandomAddress(),
			bucket: 0,
			offset: 0,
		},
		{
			addr:   test.RandomAddress(),
			bucket: 0,
			offset: 1,
		},
		{
			addr:   test.RandomAddress(),
			bucket: 1,
			offset: 0,
		},
		{
			addr:   test.RandomAddress(),
			bucket: 1,
			offset: 1,
		},
	}

	// write to new sharkyloc
	var buffer bytes.Buffer
	sw, err := sharkyloc.New(&buffer)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tc {
		err := sw.Write(tc.addr, tc.bucket, tc.offset)
		if err != nil {
			t.Fatal(err)
		}
	}

	// read from new sharyloc
	sr, err := sharkyloc.New(&buffer)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tc {
		m, err := sr.Read(tc.addr)
		if err != nil {
			t.Fatal(err)
		}

		if m.Bucket != tc.bucket {
			t.Fatalf("bucket: want %d got %d", tc.bucket, m.Bucket)
		}

		if m.Offset != tc.offset {
			t.Fatalf("offset: want %d got %d", tc.offset, m.Offset)
		}
	}
}
