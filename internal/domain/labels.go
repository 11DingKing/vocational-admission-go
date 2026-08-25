package domain

func PlanStatusLabel(s PlanStatus) string {
	switch s {
	case PlanDraft:
		return "草稿"
	case PlanLocked:
		return "已锁定"
	case PlanPublished:
		return "已发布"
	case PlanClosed:
		return "已关闭"
	default:
		return "未知"
	}
}
func BatchLabel(b Batch) string {
	switch b {
	case BatchEarly:
		return "提前批"
	case BatchRegular:
		return "普通批"
	case BatchSupplement:
		return "征集批"
	default:
		return "未知批次"
	}
}
func RoleCanDecide(r Role) bool { return r == RoleAdmin || r == RoleReviewer }
func RoleCanPlan(r Role) bool   { return r == RoleAdmin || r == RoleOfficer }
func RoleCanView(r Role) bool {
	return r == RoleAdmin || r == RoleOfficer || r == RoleReviewer || r == RoleViewer
}
func NormalizeProvince(v string) string {
	switch v {
	case "北京市":
		return "北京"
	case "上海市":
		return "上海"
	case "江苏省":
		return "江苏"
	case "浙江省":
		return "浙江"
	default:
		return v
	}
}
func NormalizeBatch(v string) Batch {
	switch v {
	case "early", "提前批":
		return BatchEarly
	case "supplement", "征集批":
		return BatchSupplement
	default:
		return BatchRegular
	}
}
