package worker
import("context";"errors";"testing";"time")
func TestRetryReturnsCancellation(t *testing.T){ctx,cancel:=context.WithCancel(context.Background());called:=false;e:=Retry(ctx,RetryPolicy{MaxAttempts:2,BaseDelay:time.Millisecond},func(context.Context)error{called=true;cancel();return errors.New("transient")});if !errors.Is(e,context.Canceled){t.Fatalf("retry returned wrong error: %v",e)};if !called{t.Fatal("retry did not make initial attempt")}}
