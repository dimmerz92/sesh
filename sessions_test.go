package sesh

import (
	"testing"
	"time"
)

type testStruct struct {
	Hello string
	Ice   string
}

var diskStore *sessionStore
var memStore *sessionStore

var diskStoreSesh string
var memStoreSesh string

var testData = testStruct{
	Hello: "world",
	Ice:   "cream",
}

func TestNewSessionStore(t *testing.T) {
	// test disk store
	dStore, err := NewSessionStore(DefaultConfig())
	if err != nil {
		t.Fatalf("failed disk store: %v", err)
	}
	dStore.Close()

	// test mem store
	mStore, err := NewSessionStore(DefaultConfig().WithInMemory(true))
	if err != nil {
		t.Fatalf("failed mem store: %v", err)
	}
	mStore.Close()
}

func TestAddToStore(t *testing.T) {
	// disk add
	dStore, err := NewSessionStore(DefaultConfig().WithSessionLength(time.Second * 2))
	if err != nil {
		t.Fatalf("failed disk store: %v", err)
	}
	diskStore = dStore

	diskStoreSesh, err = diskStore.New(testData)
	if err != nil {
		t.Fatalf("failed disk new session: %v", err)
	}

	// mem add
	mStore, err := NewSessionStore(DefaultConfig().WithInMemory(true).WithSessionLength(time.Second * 2))
	if err != nil {
		t.Fatalf("failed mem store: %v", err)
	}
	memStore = mStore

	memStoreSesh, err = memStore.New(testData)
	if err != nil {
		t.Fatalf("failed mem new session: %v", err)
	}
}

func TestGetFromStore(t *testing.T) {
	// disk get
	var diskData testStruct
	err := diskStore.Get(diskStoreSesh, &diskData)
	if err != nil {
		t.Fatalf("failed disk get: %v", err)
	} else if diskData != testData {
		t.Fatalf("diskData != testData\n%v\n%v", diskData, testData)
	}
	t.Logf("diskData successfully retrieved: %v", diskData)

	// mem get
	var memData testStruct
	err = memStore.Get(memStoreSesh, &memData)
	if err != nil {
		t.Fatalf("failed mem get: %v", err)
	} else if memData != testData {
		t.Fatalf("memData != testData\n%v\n%v", memData, testData)
	}
	t.Logf("memData successfully retrieved: %v", memData)
}

func TestDeleteFromStore(t *testing.T) {
	// disk delete
	err := diskStore.Delete(diskStoreSesh)
	if err != nil {
		t.Fatalf("failed disk delete: %v", err)
	}

	var diskData testStruct
	err = diskStore.Get(diskStoreSesh, &diskData)
	if err == nil {
		t.Fatalf("failed disk delete, not deleted: %v", err)
	}

	// mem delete
	err = memStore.Delete(memStoreSesh)
	if err != nil {
		t.Fatalf("failed mem delete: %v", err)
	}

	var memData testStruct
	err = memStore.Get(memStoreSesh, &memData)
	if err == nil {
		t.Fatalf("failed mem delete, not deleted: %v", err)
	}
}

func TestExpired(t *testing.T) {
	// disk expired add
	diskSesh, err := diskStore.New(testData)
	if err != nil {
		t.Fatalf("failed disk expired, not added: %v", err)
	}

	// mem expired add
	memSesh, err := memStore.New(testData)
	if err != nil {
		t.Fatalf("failed mem expired, not added: %v", err)
	}

	time.Sleep(time.Second * 3)

	// disk expired get
	var diskData testStruct
	err = diskStore.Get(diskSesh, &diskData)
	if err == nil {
		t.Fatalf("failed disk expired, not deleted: %v", err)
	}

	// mem expired get
	var memData testStruct
	err = memStore.Get(memSesh, &memData)
	if err == nil {
		t.Fatalf("failed mem expired, not deleted: %v", err)
	}

	diskStore.Close()
	memStore.Close()
}
