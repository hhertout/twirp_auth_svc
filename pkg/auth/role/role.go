package role

type ROLE string

// Role constants.
// ROLE_ADMIN is the admin role.
const ROLE_ADMIN ROLE = "ADMIN"

// ROLE_USER is the user role.
const ROLE_USER ROLE = "USER"

// Contains checks if a role is present in a list of roles.
func Contains(base []ROLE, compare string) bool {
	for _, role := range base {
		if role == ROLE(compare) {
			return true
		}
	}
	return false
}

// AddRole adds a role to a list of roles if it is not already present.
func AddRole(base []ROLE, newRole ROLE) []ROLE {
	if !Contains(base, string(newRole)) {
		return append(base, newRole)
	}
	return base
}

// RemoveRole removes a role from a list of roles.
func RemoveRole(base []ROLE, removeRole ROLE) []ROLE {
	for i, role := range base {
		if role == removeRole {
			return append(base[:i], base[i+1:]...)
		}
	}
	return base
}

// ToString converts a list of roles to a string slice.
func ToString() []string {
	return []string{string(ROLE_ADMIN), string(ROLE_USER)}
}

// FromString converts a string slice to a list of roles.
func FromString(roles []string) []ROLE {
	var result []ROLE
	for _, role := range roles {
		result = append(result, ROLE(role))
	}
	return result
}
