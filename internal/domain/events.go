package domain

type EventType string

const (
	EventPlanLocked           EventType = "plan.locked"
	EventApplicationSubmitted EventType = "application.submitted"
	EventDecisionRecorded     EventType = "decision.recorded"
	EventNotificationQueued   EventType = "notification.queued"
)

type Event struct {
	Type      EventType
	Entity    string
	EntityID  int64
	RequestID string
	ActorID   int64
	Payload   map[string]any
}

func (e Event) Valid() bool { return e.Type != "" && e.Entity != "" && e.EntityID > 0 }
func NewEvent(t EventType, entity string, id, actor int64, req string, p map[string]any) Event {
	return Event{Type: t, Entity: entity, EntityID: id, ActorID: actor, RequestID: req, Payload: p}
}
