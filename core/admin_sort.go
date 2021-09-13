package core

type ISortBy interface {
	Sort(afo IAdminFilterObjects, direction int)
	GetDirection() int
}

type SortBy struct {
	Direction int // -1 descending order, 1 ascending order
	Field     *Field
}

func (sb *SortBy) Sort(afo IAdminFilterObjects, direction int) {
	sortBy := sb.Field.DBName
	if direction == -1 {
		sortBy += " desc"
	}
	afo.SetPaginatedQuerySet(afo.GetPaginatedQuerySet().Order(sortBy))
}

func (sb *SortBy) GetDirection() int {
	return sb.Direction
}
