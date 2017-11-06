package util

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
)

func GetLocalIP() string {
	var ip string
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if !strings.HasPrefix(addr.String(), "127.0.0.1") && !strings.Contains(addr.String(), ":") {
			ipAddr := addr.String()
			elements := strings.Split(ipAddr, "/")
			ip = elements[0]
			log.Info().Msgf("GetLocalIP: using %s as local addr", ip)
			break
		}
	}
	return ip
}

func GenerateHash(username string, req *http.Request) string {
	tokenSeed := username + ";" + req.RemoteAddr

	hasher := sha256.New()
	hasher.Write([]byte(tokenSeed))
	token := hasher.Sum(nil)

	return base64.URLEncoding.EncodeToString(token)
}

func Now() int64 {
	return time.Now().UnixNano() / 1000000
}

type ID string

func NewId() ID {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return ID(id.String())
}
