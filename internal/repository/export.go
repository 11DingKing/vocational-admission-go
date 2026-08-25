package repository

import (
	"context"
	"database/sql"
	"encoding/csv"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"io"
	"strconv"
)

func ExportApplications(ctx context.Context, db *sql.DB, w io.Writer, planID int64) error {
	cw := csv.NewWriter(w)
	if e := cw.Write([]string{"id", "student_no", "score", "rank", "status"}); e != nil {
		return e
	}
	rows, e := db.QueryContext(ctx, "SELECT id,student_no,score,rank,status FROM applications WHERE plan_id=? ORDER BY rank", planID)
	if e != nil {
		return e
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var student, status string
		var score, rank int
		if e = rows.Scan(&id, &student, &score, &rank, &status); e != nil {
			return e
		}
		if e = cw.Write([]string{strconv.FormatInt(id, 10), student, strconv.Itoa(score), strconv.Itoa(rank), status}); e != nil {
			return e
		}
	}
	cw.Flush()
	return cw.Error()
}
func ApplicationStatusLabel(s domain.ApplicationStatus) string { return string(s) }
