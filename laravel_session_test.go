package golaravelsession

import (
	"testing"
)

func TestGetSessionID(t *testing.T) {
	cookie := "eyJpdiI6IjVrTVVDSmlyb1FtN2NrbmlOTllOUkE9PSIsInZhbHVlIjoiUU9pbCtMTjhQQnQyamJ6ZE5qenVWanhuZktUcjBkOUVsWU5ibkhlWHJyc25DNnZYQlRrOWlFd01ObVJwam1yVUtNcGRUanV1aEJIWHBsYXNiZytNenc9PSIsIm1hYyI6ImMzYzVmMGE1NWY5ZjEzMzRjOTVkN2FlZGY2YzZhNDExOTVhZjUzMjYzZmE3OTE1ODIwYWYzNmY5ODQzYjIwOGEifQ=="
	key := "base64:qsDvCdhT+JPXEBD3ys/XraOXVNpshsyElzJmtgnBqEI="
	expectedSessionID := "RYodG5AekDidQCVLvs4fQIRAPSwarZV26U4shNVX"

	sessionID, err := GetSessionID(cookie, key)
	if err != nil {
		t.Errorf("fail - %s", err)
	}

	if sessionID == expectedSessionID {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}

func TestGetSessionIDWithoutSerialization(t *testing.T) {
	cookie := "eyJpdiI6IlwvVjl1MTBpbGxpMUwxdEZsemRZV0Z3PT0iLCJ2YWx1ZSI6IjZRTlV0cjNTK1NlQ2NkTTAxSUF3VExKSXAra3VRS3RRV0JCaHBvc1lNNkNFVzliY1k1WUlndVJCT09RNENvdnciLCJtYWMiOiIzZTMyYmVkODcwZmVjNjBhZjY1MjkxYWQyZGRiNWMxMTg4ODJlZmNkOTJmZmUxMDcwNWYwYjMzZTY0NzM4ZjUxIn0="
	key := "base64:qsDvCdhT+JPXEBD3ys/XraOXVNpshsyElzJmtgnBqEI="
	expectedSessionID := "uJUzLDZDWUq48O1RWVUbgWDUeQ1vJz5YQdgbQaDO"

	sessionID, err := GetSessionID(cookie, key)
	if err != nil {
		t.Errorf("fail - %s", err)
	}

	if sessionID == expectedSessionID {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}

func TestParseSessionData(t *testing.T) {
	sessionData := `a:5:{s:6:"_token";s:40:"eE5gVNqGSn6wneCJAzhtTMulPFwvOfDyZRSoVStA";s:3:"url";a:0:{}s:9:"_previous";a:1:{s:3:"url";s:29:"https://jshb-admin.dev/update";}s:6:"_flash";a:2:{s:3:"old";a:0:{}s:3:"new";a:0:{}}s:52:"login_admin_59ba36addc2b2f9401580f014c7f58ea4e30989d";i:1;}`
	session, err := ParseSessionData(sessionData)

	if err != nil {
		t.Errorf("fail - %s", err)
	}

	if session["_token"].(string) == "eE5gVNqGSn6wneCJAzhtTMulPFwvOfDyZRSoVStA" {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}
