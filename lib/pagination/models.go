package pagination

type TableInfo struct {
	PrimaryField         string
	PrimaryFieldToString func(any) string
}

type DatabaseInfo struct {
	Tables map[string]TableInfo
}

type CursorInfo struct {
	Direction string
	PrimaryId string
	SortId    string
	SortField string
}

type SortInfo struct {
	Field      string
	Order      string
	Limit      int
	Table      string
	CursorInfo *CursorInfo
}

type DBResult struct {
	PID    string
	SortID string
}
