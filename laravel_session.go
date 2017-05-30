package golaravelsession

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"strings"

	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

//GetSessionID Decypt laravel's session cookie and return it's SessionID
func GetSessionID(cookie, key string) (string, error) {

	decodeBytes, err := base64.StdEncoding.DecodeString(cookie)
	if err != nil {
		return "", errors.New("cookie value must in base64 format")
	}
	var payload struct {
		IV    string
		Value string
		Mac   string
	}
	err = json.Unmarshal(decodeBytes, &payload)
	if err != nil {
		return "", errors.New("cookie value must be valid")
	}
	encryptedText, err := base64.StdEncoding.DecodeString(payload.Value)
	if err != nil {
		return "", errors.New("encrypted text must be valid base64 format")
	}
	iv, err := base64.StdEncoding.DecodeString(payload.IV)
	if err != nil {
		return "", errors.New("iv in payload must be valid base64 format")
	}

	var keyBytes []byte
	if strings.HasPrefix(key, "base64:") {
		keyBytes, err = base64.StdEncoding.DecodeString(string(key[7:]))
		if err != nil {
			return "", errors.New("seems like you provide a key in base64 format, but it's not valid")
		}
	} else {
		keyBytes = []byte(key)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptedText, encryptedText)

	unserializer := php_serialize.NewUnSerializer(string(encryptedText))
	var sessionID php_serialize.PhpValue
	if sessionID, err = unserializer.Decode(); err != nil {
		return "", fmt.Errorf("unserialize data error, raw value is %s", string(encryptedText))
	}

	return sessionID.(string), nil
}

//ParseSessionData Parse session data into readable golang struct
func ParseSessionData(data string) (php_serialize.PhpArray, error) {
	decoder := php_serialize.NewUnSerializer(data)
	sessionDataDecoded, err := decoder.Decode()
	if err != nil {
		return nil, err
	}
	return sessionDataDecoded.(php_serialize.PhpArray), nil
}
