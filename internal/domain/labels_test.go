package domain

import "testing"

func TestLabelsAndRoles(t *testing.T) {
	plans := []PlanStatus{PlanDraft, PlanLocked, PlanPublished, PlanClosed}
	for _, p := range plans {
		if PlanStatusLabel(p) == "未知" {
			t.Errorf("label %s", p)
		}
	}
	batches := []Batch{BatchEarly, BatchRegular, BatchSupplement}
	for _, b := range batches {
		if BatchLabel(b) == "未知批次" {
			t.Errorf("batch %s", b)
		}
	}
	roles := []Role{RoleAdmin, RoleOfficer, RoleReviewer, RoleViewer}
	for _, r := range roles {
		if !RoleCanView(r) {
			t.Errorf("view %s", r)
		}
	}
	if !RoleCanPlan(RoleOfficer) || RoleCanPlan(RoleViewer) {
		t.Fatal("plan role")
	}
	if !RoleCanDecide(RoleReviewer) || RoleCanDecide(RoleOfficer) {
		t.Fatal("decision role")
	}
}
func TestProvinceNormalization(t *testing.T) {
	cases := map[string]string{"北京市": "北京", "上海市": "上海", "江苏省": "江苏", "浙江省": "浙江", "广东": "广东", "四川": "四川", "湖北": "湖北", "山东": "山东", "福建": "福建", "安徽": "安徽", "河南": "河南", "河北": "河北", "湖南": "湖南", "江西": "江西", "陕西": "陕西", "云南": "云南", "辽宁": "辽宁", "吉林": "吉林", "黑龙江": "黑龙江", "广西": "广西", "贵州": "贵州", "甘肃": "甘肃", "海南": "海南", "青海": "青海", "宁夏": "宁夏", "新疆": "新疆", "西藏": "西藏", "内蒙古": "内蒙古", "天津": "天津", "重庆": "重庆"}
	for in, want := range cases {
		if got := NormalizeProvince(in); got != want {
			t.Errorf("%s=%s", in, got)
		}
	}
}
func TestBatchNormalization(t *testing.T) {
	for _, v := range []string{"early", "提前批"} {
		if NormalizeBatch(v) != BatchEarly {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"supplement", "征集批"} {
		if NormalizeBatch(v) != BatchSupplement {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"regular", "普通批", "unknown"} {
		if NormalizeBatch(v) != BatchRegular {
			t.Fatal(v)
		}
	}
}
