package domain

// UserID is alias for identifying user
type UserID string

// CategoryID is alias for identifying category
type CategoryID string

// LabelID is alias for identifying label
type LabelID string

// LangID is alias for identifying lang
type LangID string

// TenantID is alias for identifying tenant
type TenantID string

// LabelStatus is status for Label
type LabelStatus int

// String provides exchange UserID to string
func (i UserID) String() string {
	return string(i)
}

// String provides exchange CategoryID to string
func (i CategoryID) String() string {
	return string(i)
}

// String provides exchange LabelID to string
func (i LabelID) String() string {
	return string(i)
}

// String provides exchange LangID to string
func (i LangID) String() string {
	return string(i)
}

// String provides exchange TenantID to string
func (i TenantID) String() string {
	return string(i)
}
