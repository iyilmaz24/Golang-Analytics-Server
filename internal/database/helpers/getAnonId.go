package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/iyilmaz24/Go-Analytics-Server/internal/config"
)

func GetAnonymousID(ip, region string) string { 

	truncatedIP := getTruncatedIP(ip) // truncate the IP address for privacy
	if truncatedIP == "" { // if the IP address is invalid, return a default value
		return "invalid-ip"
	}

	stringHash := getStringHash(truncatedIP) // create a irreversible hash of the truncated IP address
	anonID := generateUUID(stringHash, region) // generate a UUID based on the hash and the namespace for the region

	return anonID
}

func getTruncatedIP(id string) string {

	delimitter := "."
	if strings.Contains(id, ":") { // check if the IP address is IPv6
		delimitter = ":" 
	}

	expandedIp := net.ParseIP(id) // expand the IP address to a full 16 byte representation

	if expandedIp == nil { // if the IP address is invalid, net.ParseIP returns nil
		return ""
	}

	parts := strings.Split(expandedIp.String(), delimitter)
	for len(parts) < 3 { // if the IP address has less than 4 octets, pad with 0s
		parts = append(parts, "0") 
	}

	newIp := parts[0]+parts[1]+parts[2]	// take the first 3 octets of the IP address, drop rest for privacy

	return newIp
}

func getStringHash(str string) string {
	appConfig := config.LoadConfig()
	salt := appConfig.Salt

	data := str + salt
	hash := sha256.Sum256([]byte(data)) // create a SHA-256 hash of the string and salt

	return hex.EncodeToString(hash[:]) // return the hash as a hex string
}

func generateUUID(stringHash, region string) string {
	appConfig := config.LoadConfig()
	var namespace uuid.UUID 

	if region == "FL" { // get the UUID namespace based on the region
		namespace = appConfig.FL_NS
	} else {
		namespace = appConfig.NM_NS
	}
	
	id := uuid.NewSHA1(namespace, []byte(stringHash)) // generate a UUID based on the hash and the namespace

	return id.String()
}

