package constant

import "fmt"

const (
	userInfoKey  = "user:info:%d"
	userTokenKey = "user:token:%s"
)

func GenUserInfoKey(id uint) string {
	return fmt.Sprintf(userInfoKey, id)
}

func GenUserTokenKey(id string) string {
	return fmt.Sprintf(userTokenKey, id)
}
