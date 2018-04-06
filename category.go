package labeltile

// Category manages label category
type Category struct {
	ID          CategoryID
	Tenant      TenantID
	Name        string
	Description string
}
