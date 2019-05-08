package golaravelsession

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"

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
		// laravel > 5.5.40 or 5.6.29 session keys will not be serialized
		// https://github.com/laravel/framework/pull/25121
		return string(removePadding(encryptedText)), nil
	}

	return sessionID.(string), nil
}

func removePadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

//ParseSessionData Parse session data into readable golang struct
func ParseSessionData(data string) (php_serialize.PhpArray, error) {
	decoder := php_serialize.NewUnSerializer(data)
	sessionDataDecoded, err := decoder.Decode()
	if err != nil {
		return nil, err
	}

	// Some cache-based session stores such as redis result in double
	// serialization of session data. If the session data is a string
	// then this has almost certainly happened.
	if sessionDataDecodedString, ok := sessionDataDecoded.(string); ok {
		return ParseSessionData(sessionDataDecodedString)
	}

	return sessionDataDecoded.(php_serialize.PhpArray), nil
}
