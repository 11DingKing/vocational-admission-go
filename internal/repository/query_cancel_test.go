package repository
import("context";"testing";"github.com/11DingKing/vocational-admission-go/internal/storage")
func TestPlanSummaryHonorsCanceledContext(t *testing.T){db,e:=storage.Open(context.Background(),"file:summary?mode=memory&cache=shared");if e!=nil{t.Fatal(e)};defer db.Close();if e=storage.Migrate(context.Background(),db.SQL);e!=nil{t.Fatal(e)};ctx,cancel:=context.WithCancel(context.Background());cancel();if _,e=PlanSummary(ctx,db.SQL,1);e==nil{t.Fatal("canceled summary succeeded")}}
