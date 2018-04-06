package labeltile

// Tenant manages labels
type Tenant struct {
	ID          TenantID
	DefaultLang LangID
	Languages   []LangID
	Categories  []*Category
}
