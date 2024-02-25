package strhash_test

import (
	"context"
	"errors"
	"fmt"
	"hashcalc/internal/strhash"
	"hashkeeper/pkg/hashlog"
	"sync"
	"testing"
)

// black box tests

func TestInitUndef(t *testing.T) {
	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.ModeUndef})

	if err == nil {
		t.Fatal("error can`t be nil")
	}

	if hashMaker != nil {
		t.Fatal("hash maker must be nil")
	}
}

func TestInit224(t *testing.T) {
	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.Mode224, BuffersCapacity: 1})
	if err != nil {
		t.Fatal(err)
	}

	in := strhash.InItem{
		Data: []byte("line28"),
	}

	hashMaker.InChan() <- in
	close(hashMaker.InChan())

	err = hashMaker.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	out := <-hashMaker.OutChan()

	if err := compareLength(out, 28); err != nil {
		t.Fatal(err)
	}

	if err := compareHash(out, "cb00f7c0a0ebda6468c9d85a3cc61c7a3f36d879a350fa34dbcd0be9"); err != nil {
		t.Fatal(err)
	}
}

func TestInit256(t *testing.T) {
	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.Mode256, BuffersCapacity: 1})
	if err != nil {
		t.Fatal(err)
	}

	in := strhash.InItem{
		Data: []byte("line32"),
	}

	hashMaker.InChan() <- in
	close(hashMaker.InChan())

	err = hashMaker.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	out := <-hashMaker.OutChan()

	if err := compareLength(out, 32); err != nil {
		t.Fatal(err)
	}

	if err := compareHash(out, "c77bbb8ddc0b29a759fb17c789dbdb5466e9499c3b1a94923e7745a1174c44eb"); err != nil {
		t.Fatal(err)
	}
}

func TestInit384(t *testing.T) {
	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.Mode384, BuffersCapacity: 1})
	if err != nil {
		t.Fatal(err)
	}

	in := strhash.InItem{
		Data: []byte("line48"),
	}

	hashMaker.InChan() <- in
	close(hashMaker.InChan())

	err = hashMaker.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	out := <-hashMaker.OutChan()

	if err := compareLength(out, 48); err != nil {
		t.Fatal(err)
	}

	if err := compareHash(out, "d98c16109a8c5e725e77c6682f62c77b50119d084126430e78196a0944c3d263a8574296c26966a043b61d924f4b9505"); err != nil {
		t.Fatal(err)
	}
}

func TestInit512(t *testing.T) {
	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.Mode512, BuffersCapacity: 1})
	if err != nil {
		t.Fatal(err)
	}

	in := strhash.InItem{
		Data: []byte("line64"),
	}

	hashMaker.InChan() <- in
	close(hashMaker.InChan())

	err = hashMaker.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	out := <-hashMaker.OutChan()

	if err := compareLength(out, 64); err != nil {
		t.Fatal(err)
	}

	if err := compareHash(out, "07d6f2e3fc5d2a75649d236d8fa4436dab71bf9cab679840ac618f0360bab73aa4bc344feb6807d484d37ac59d7afdad96213fcf9ee41dcabd5b2bc12fbaa501"); err != nil {
		t.Fatal(err)
	}
}

func TestMultiple(t *testing.T) {
	const buffersCapacity = 1000

	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.Mode256, BuffersCapacity: buffersCapacity})
	if err != nil {
		t.Fatal(err)
	}

	for _, in := range makeExampleIn(buffersCapacity) {
		hashMaker.InChan() <- in
	}
	close(hashMaker.InChan())

	err = hashMaker.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	var outs []strhash.OutItem

	for out := range hashMaker.OutChan() {
		outs = append(outs, out)
	}

	if len(outs) != buffersCapacity {
		t.Fatalf("wrong length of output slice, expected: %d, received: %d", buffersCapacity, len(outs))
	}
}

func TestCancel(t *testing.T) {
	const buffersCapacity = 1000

	hashMaker, err := strhash.NewSha3StrHashMaker(strhash.Sha3StrHashCfg{Mode: strhash.Mode256, BuffersCapacity: buffersCapacity})
	if err != nil {
		t.Fatal(err)
	}

	for _, in := range makeExampleIn(buffersCapacity) {
		hashMaker.InChan() <- in
	}
	close(hashMaker.InChan())

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	var runErr error

	go func() {
		runErr = hashMaker.Run(ctx)
		wg.Done()
	}()

	cancel()
	wg.Wait()
	if runErr == nil {
		t.Fatal("error can`t be nil")
	}

	if !errors.Is(runErr, context.Canceled) {
		t.Fatal("error must be context canceled")
	}
}

// Helpers

func compareLength(out strhash.OutItem, expectedLen int) error {
	if len(out.Hash) != expectedLen {
		return hashlog.WithStackErrorf("wrong length of output, expected: %d, received: %d", expectedLen, len(out.Hash))
	}
	return nil
}

func compareHash(out strhash.OutItem, expectedStr string) error {
	curStr := fmt.Sprintf("%x", out.Hash)
	if curStr != expectedStr {
		return hashlog.WithStackErrorf("wrong hash of output, expected: %s, received: %s", expectedStr, curStr)
	}
	return nil
}

func makeExampleIn(length int) []strhash.InItem {
	in := make([]strhash.InItem, length)
	for i := 0; i < length; i++ {
		in[i].Data = []byte(fmt.Sprintf("line%d", i))
	}
	return in
}
