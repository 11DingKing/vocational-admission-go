package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"strings"
	"time"
)

type Notification struct {
	ApplicationID int64
	Channel       string
	Recipient     string
	Message       string
	CreatedAt     time.Time
}

func BuildNotification(a domain.Application, channel string) Notification {
	return Notification{ApplicationID: a.ID, Channel: channel, Recipient: a.StudentNo, Message: fmt.Sprintf("申请 %d 当前状态：%s", a.ID, a.Status), CreatedAt: time.Now().UTC()}
}
func ValidateChannel(ch string) bool {
	switch strings.ToLower(ch) {
	case "sms", "email", "portal":
		return true
	default:
		return false
	}
}
func (s AdmissionService) QueueNotification(ctx context.Context, a domain.Application, channel string) error {
	if !ValidateChannel(channel) {
		return domain.ErrConflict
	}
	return s.DB.Tx(ctx, func(tx *sql.Tx) error { return s.Jobs.Enqueue(ctx, tx, "notify_"+channel, a.ID) })
}
func RenderStatus(status domain.ApplicationStatus) string {
	switch status {
	case domain.ApplicationSubmitted:
		return "待审核"
	case domain.ApplicationReviewing:
		return "审核中"
	case domain.ApplicationAdmitted:
		return "已录取"
	case domain.ApplicationRejected:
		return "未录取"
	case domain.ApplicationWithdrawn:
		return "已撤回"
	default:
		return "未知"
	}
}
func IsTerminal(status domain.ApplicationStatus) bool {
	return status == domain.ApplicationAdmitted || status == domain.ApplicationRejected || status == domain.ApplicationWithdrawn
}
