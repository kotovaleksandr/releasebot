package main

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestReadAllUsers(t *testing.T) {
	DataFileName = tempFileName()
	connectedUsers := GetUsers()
	if len(connectedUsers) != 0 {
		t.Error("Users not found")
	}
	_, err := os.Open(DataFileName)
	if err != nil {
		t.Errorf("Data file %s is not created: %v", DataFileName, err)
	}
}

func TestAddUserToDataFile(t *testing.T) {
	DataFileName = tempFileName()
	AddUser(123)
	connectedUsers := GetUsers()
	if len(connectedUsers) == 0 {
		t.Error("User not added")
	}
	if connectedUsers[0] != 123 {
		t.Errorf("User id corrupted: %d", connectedUsers[0])
	}
}

func tempFileName() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), hex.EncodeToString(randBytes))
}
