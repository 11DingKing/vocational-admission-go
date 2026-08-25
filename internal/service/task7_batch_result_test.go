package service
import("context";"testing";"github.com/11DingKing/vocational-admission-go/internal/domain")
func TestProcessBatchReportsEveryDecision(t *testing.T){s,_,pid,gid:=svc(t);a,e:=s.Submit(context.Background(),domain.Application{PlanID:pid,MajorGroupID:gid,StudentNo:"20267001",Score:600,Rank:1,IdempotencyKey:"b1"},1,"r");if e!=nil{t.Fatal(e)};out:=s.ProcessBatch(context.Background(),[]int64{999999,a.ID},domain.ApplicationReviewing,2,"batch");if len(out)!=2{t.Fatalf("batch stopped after failure: %d",len(out))}}
