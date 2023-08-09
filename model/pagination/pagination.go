// Package pagination defines helper models for pagination related functionalities
package pagination

const (
	// PageDefaultNumber int value 1
	PageDefaultNumber int = 1
	// PageDefaultSize int value 10
	PageDefaultSize int = 10
	// PageDefaultSortBy default sortBy string value
	PageDefaultSortBy string = "created_at"
	// PageDefaultSortDirectionDesc default sort direction descending order status
	PageDefaultSortDirectionDesc bool = true
	// PageSortDirectionAscending string value asc
	PageSortDirectionAscending string = "asc"
	// PageSortDirectionDescending string value desc
	PageSortDirectionDescending string = "desc"
	// SortByEnterDate sort by entered_at on instalments table
	SortByEnterDate string = "entered_at"
	// SortByDaysLate sort by due_date on instalments table
	SortByDaysLate string = "due_date"
	// SortByCreatedAt sort by created_at on instalments table
	SortByCreatedAt string = "created_at"
	// SortByUpdatedAt sort by updated_at on instalments table
	SortByUpdatedAt string = "updated_at"
	// SortByName sort by name on customers table
	SortByName string = "name"
	// SortByTags sort by tags on customers table
	SortByTags string = "tags"
	// SortByAmountToPayIncludingFees sort by amount_to_pay_including_fees on instalments table
	SortByAmountToPayIncludingFees string = "amount_to_pay_including_fees"
	// SortByOriginalAmountToPay sort by original_principal_to_pay on instalments table
	SortByOriginalAmountToPay string = "original_principal_to_pay"
)

// Page object for pagination purpose. Not persisted
type Page struct {
	Number            *int
	Size              *int
	SortBy            *string
	SortDirectionDesc *bool
}

// PageInfo holds pagination response info
type PageInfo struct {
	Page            int
	Size            int
	HasNextPage     bool
	HasPreviousPage bool
	TotalCount      int64
}

// NewPage creates a new pagination Page object
func NewPage(n int, s int, sBy string, sDirectionD bool) Page {
	return Page{
		Number:            &n,
		Size:              &s,
		SortBy:            &sBy,
		SortDirectionDesc: &sDirectionD,
	}
}

// NewPageWithDefaultSorting creates a new pagination Page object with system default values
func NewPageWithDefaultSorting(n int, s int) Page {
	return Page{
		Number: &n,
		Size:   &s,
	}
}
