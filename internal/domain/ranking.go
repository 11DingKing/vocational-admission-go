package domain

import "sort"

type RankedApplication struct {
	Application Application
	Priority    int
}

func RankApplications(items []Application) []Application {
	out := append([]Application(nil), items...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].Rank != out[j].Rank {
			return out[i].Rank < out[j].Rank
		}
		return out[i].SubmittedAt.Before(out[j].SubmittedAt)
	})
	return out
}
func Eligible(app Application, plan AdmissionPlan, group MajorGroup) bool {
	return app.PlanID == plan.ID && app.MajorGroupID == group.ID && plan.Status == PlanPublished && plan.Remaining() > 0 && group.Remaining() > 0
}
func Allocate(items []Application, capacity int) (admitted, rejected []Application) {
	ranked := RankApplications(items)
	for i, a := range ranked {
		if i < capacity {
			admitted = append(admitted, a)
		} else {
			rejected = append(rejected, a)
		}
	}
	return
}
