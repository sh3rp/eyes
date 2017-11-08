package util

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/matishsiao/goInfo"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/msg"
)

var VERSION_MAJOR = 0
var VERSION_MINOR = 1
var VERSION_PATCH = 0

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

func GenNodeInfo(id ID) msg.NodeInfo {
	info := goInfo.GetInfo()

	return msg.NodeInfo{
		Id:           string(id),
		Os:           info.GoOS,
		Kernel:       info.Core,
		Platform:     info.Platform,
		Ip:           "",
		Hostname:     info.Hostname,
		CoreCount:    int32(info.CPUs),
		MajorVersion: int32(VERSION_MAJOR),
		MinorVersion: int32(VERSION_MINOR),
		PatchVersion: int32(VERSION_PATCH),
	}
}
