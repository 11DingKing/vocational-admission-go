package service
import("context";"testing";"github.com/11DingKing/vocational-admission-go/internal/domain")
func TestLockAuditFailureKeepsDraft(t *testing.T){s,db,pid,_:=svc(t);if _,e:=db.SQL.Exec("CREATE TRIGGER fail_lock_audit BEFORE INSERT ON audit_events BEGIN SELECT RAISE(ABORT,'audit down'); END");e!=nil{t.Fatal(e)};if e:=s.LockPlan(context.Background(),pid,1,"lock");e==nil{t.Fatal("expected audit error")};p,e:=s.Plans.ByID(context.Background(),pid);if e!=nil{t.Fatal(e)};if p.Status!=domain.PlanDraft{t.Fatalf("plan committed despite audit failure: %s",p.Status)}}
