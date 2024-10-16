package role

type ROLE string

const ROLE_ADMIN ROLE = "ADMIN"
const ROLE_USER ROLE = "USER"

func Contains(base []ROLE, compare string) bool {
	for _, role := range base {
		if role == ROLE(compare) {
			return true
		}
	}
	return false
}
