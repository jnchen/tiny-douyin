package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"io"
	"strings"
)

func UUID() string {
	id := uuid.New()
	return id.String()
}

func UUIDNoLine() string {
	str := UUID()
	return strings.ReplaceAll(str, "-", "")
}

func Md5(content string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, content)
	if err != nil {
		return "", err
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum[:]), nil
}

func SHA256(content string) (string, error) {
	h := sha256.New()
	_, err := io.WriteString(h, content)
	if err != nil {
		return "", err
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum[:]), nil
}
