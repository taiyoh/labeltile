package domain

// Tag manages label keywords
type Tag struct {
	ID     TagID
	Tenant TenantID
	Name   string
}
