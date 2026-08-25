package pagination

type Query struct {
	Limit  int
	Offset int
}

func Parse(limit, offset int) Query {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return Query{limit, offset}
}
func (q Query) Next(total int) bool { return q.Offset+q.Limit < total }
