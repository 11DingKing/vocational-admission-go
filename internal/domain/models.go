package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidState = errors.New("invalid state transition")
	ErrCapacity     = errors.New("capacity exceeded")
	ErrDuplicate    = errors.New("duplicate request")
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOfficer  Role = "officer"
	RoleReviewer Role = "reviewer"
	RoleViewer   Role = "viewer"
)

type PlanStatus string

const (
	PlanDraft     PlanStatus = "draft"
	PlanLocked    PlanStatus = "locked"
	PlanPublished PlanStatus = "published"
	PlanClosed    PlanStatus = "closed"
)

type ApplicationStatus string

const (
	ApplicationSubmitted ApplicationStatus = "submitted"
	ApplicationReviewing ApplicationStatus = "reviewing"
	ApplicationAdmitted  ApplicationStatus = "admitted"
	ApplicationRejected  ApplicationStatus = "rejected"
	ApplicationWithdrawn ApplicationStatus = "withdrawn"
)

type Batch string

const (
	BatchEarly      Batch = "early"
	BatchRegular    Batch = "regular"
	BatchSupplement Batch = "supplement"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         Role
	Active       bool
	CreatedAt    time.Time
}
type Session struct {
	ID        string
	UserID    int64
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
type ProvinceRule struct {
	ID            int64
	Province      string
	Batch         Batch
	MinScore      int
	RankLimit     int
	AllowTransfer bool
	Version       int
	UpdatedAt     time.Time
}
type AdmissionPlan struct {
	ID            int64
	Year          int
	Province      string
	Name          string
	Status        PlanStatus
	TotalCapacity int
	UsedCapacity  int
	Version       int
	LockedAt      *time.Time
	CreatedAt     time.Time
}
type MajorGroup struct {
	ID           int64
	PlanID       int64
	Code         string
	Name         string
	Capacity     int
	UsedCapacity int
	Version      int
}
type Application struct {
	ID               int64
	PlanID           int64
	MajorGroupID     int64
	StudentNo        string
	Score            int
	Rank             int
	Preferences      []string
	Status           ApplicationStatus
	TransferAccepted bool
	IdempotencyKey   string
	Version          int
	SubmittedAt      time.Time
	UpdatedAt        time.Time
}
type Decision struct {
	ID            int64
	ApplicationID int64
	ActorID       int64
	FromStatus    string
	ToStatus      string
	Reason        string
	RequestID     string
	CreatedAt     time.Time
}
type AuditEvent struct {
	ID        int64
	ActorID   int64
	Entity    string
	EntityID  int64
	Action    string
	Outcome   string
	RequestID string
	Details   string
	CreatedAt time.Time
}
type Job struct {
	ID          int64
	Kind        string
	EntityID    int64
	Attempts    int
	AvailableAt time.Time
	CompletedAt *time.Time
	LastError   string
	CreatedAt   time.Time
}

func (p PlanStatus) CanTransition(to PlanStatus) bool {
	switch p {
	case PlanDraft:
		return to == PlanLocked
	case PlanLocked:
		return to == PlanPublished || to == PlanDraft
	case PlanPublished:
		return to == PlanClosed
	case PlanClosed:
		return false
	}
	return false
}
func (a ApplicationStatus) CanTransition(to ApplicationStatus) bool {
	switch a {
	case ApplicationSubmitted:
		return to == ApplicationReviewing || to == ApplicationWithdrawn
	case ApplicationReviewing:
		return to == ApplicationAdmitted || to == ApplicationRejected || to == ApplicationWithdrawn
	case ApplicationAdmitted, ApplicationRejected, ApplicationWithdrawn:
		return false
	}
	return false
}
func (p AdmissionPlan) Remaining() int { return p.TotalCapacity - p.UsedCapacity }
func (g MajorGroup) Remaining() int    { return g.Capacity - g.UsedCapacity }
func ValidateScore(r ProvinceRule, score int, rank int) error {
	if score < r.MinScore {
		return fmt.Errorf("%w: score below threshold", ErrForbidden)
	}
	if r.RankLimit > 0 && rank > r.RankLimit {
		return fmt.Errorf("%w: rank exceeds limit", ErrForbidden)
	}
	return nil
}
