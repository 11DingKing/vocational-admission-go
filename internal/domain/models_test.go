package domain

import "testing"

func TestPlanTransitions(t *testing.T) {
	cases := []struct {
		from, to PlanStatus
		ok       bool
	}{
		{PlanDraft, PlanLocked, true}, {PlanDraft, PlanPublished, false}, {PlanDraft, PlanClosed, false},
		{PlanLocked, PlanPublished, true}, {PlanLocked, PlanDraft, true}, {PlanLocked, PlanClosed, false},
		{PlanPublished, PlanClosed, true}, {PlanPublished, PlanDraft, false}, {PlanPublished, PlanLocked, false},
		{PlanClosed, PlanDraft, false}, {PlanClosed, PlanLocked, false}, {PlanClosed, PlanPublished, false},
	}
	for _, tc := range cases {
		if got := tc.from.CanTransition(tc.to); got != tc.ok {
			t.Errorf("%s to %s = %v", tc.from, tc.to, got)
		}
	}
}
func TestApplicationTransitions(t *testing.T) {
	cases := []struct {
		from, to ApplicationStatus
		ok       bool
	}{
		{ApplicationSubmitted, ApplicationReviewing, true}, {ApplicationSubmitted, ApplicationWithdrawn, true}, {ApplicationSubmitted, ApplicationAdmitted, false}, {ApplicationSubmitted, ApplicationRejected, false},
		{ApplicationReviewing, ApplicationAdmitted, true}, {ApplicationReviewing, ApplicationRejected, true}, {ApplicationReviewing, ApplicationWithdrawn, true}, {ApplicationReviewing, ApplicationSubmitted, false},
		{ApplicationAdmitted, ApplicationReviewing, false}, {ApplicationAdmitted, ApplicationRejected, false}, {ApplicationRejected, ApplicationAdmitted, false}, {ApplicationWithdrawn, ApplicationReviewing, false},
	}
	for _, tc := range cases {
		if got := tc.from.CanTransition(tc.to); got != tc.ok {
			t.Errorf("%s to %s = %v", tc.from, tc.to, got)
		}
	}
}
func TestValidationBoundaries(t *testing.T) {
	for _, v := range []struct {
		province string
		valid    bool
	}{{"北京", true}, {"浙江", true}, {"", false}, {"a", false}, {"\n", false}, {"广东省", true}} {
		if got := ValidProvince(v.province); got != v.valid {
			t.Errorf("province %q", v.province)
		}
	}
	for _, v := range []struct {
		student string
		valid   bool
	}{{"20260001", true}, {"A-9001", true}, {"abc-123", true}, {"", false}, {"12", false}, {"2026_01", false}} {
		if got := ValidStudentNo(v.student); got != v.valid {
			t.Errorf("student %q", v.student)
		}
	}
}
func TestRankingAndAllocation(t *testing.T) {
	items := []Application{{ID: 1, Score: 600, Rank: 2}, {ID: 2, Score: 650, Rank: 8}, {ID: 3, Score: 650, Rank: 3}, {ID: 4, Score: 500, Rank: 1}}
	got := RankApplications(items)
	if got[0].ID != 3 || got[1].ID != 2 {
		t.Fatalf("rank order %#v", got)
	}
	a, r := Allocate(items, 2)
	if len(a) != 2 || len(r) != 2 {
		t.Fatalf("allocation %d %d", len(a), len(r))
	}
}
func TestScoreRule(t *testing.T) {
	rule := ProvinceRule{MinScore: 500, RankLimit: 100}
	if ValidateScore(rule, 499, 1) == nil {
		t.Fatal("low score accepted")
	}
	if ValidateScore(rule, 500, 101) == nil {
		t.Fatal("rank accepted")
	}
	if ValidateScore(rule, 500, 100) != nil {
		t.Fatal("valid score rejected")
	}
}
func TestRemaining(t *testing.T) {
	p := AdmissionPlan{TotalCapacity: 10, UsedCapacity: 4}
	if p.Remaining() != 6 {
		t.Fatal(p.Remaining())
	}
	g := MajorGroup{Capacity: 3, UsedCapacity: 3}
	if g.Remaining() != 0 {
		t.Fatal(g.Remaining())
	}
}
