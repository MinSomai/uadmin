package core

type AdminRequestPaginator struct {
	PerPage int
	Offset  int
}

type PaginationType string

var LimitPaginationType PaginationType = "limit"
var CursorPaginationType PaginationType = "cursor"

type IPaginationInterface interface {
	Paginate(afo IAdminFilterObjects)
}

type Paginator struct {
	PerPage                    int
	AllowEmptyFirstPage        bool
	ShowLastPageOnPreviousPage bool
	Count                      int
	NumPages                   int
	Offset                     int
	Template                   string
	PaginationType             PaginationType
}

func (p *Paginator) Paginate(afo IAdminFilterObjects) {

}

