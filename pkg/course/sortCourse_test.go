package course

import (
	"testing"
)

func TestParseSksj(t *testing.T) {
	// 多时间段
	slots := ParseSksj("星期四第3-5节{7-8周};星期四第3-5节{1-6周,9-17周}")
	if len(slots) != 2 {
		t.Fatalf("期望2个时间段，得到%d", len(slots))
	}
	if slots[0].Day != 4 || slots[0].P1 != 3 || slots[0].P2 != 5 {
		t.Errorf("时间段解析错误: %+v", slots[0])
	}
	if !slots[0].Weeks[7] || !slots[0].Weeks[8] || slots[0].Weeks[9] {
		t.Errorf("周次解析错误: %+v", slots[0].Weeks)
	}
	if !slots[1].Weeks[1] || slots[1].Weeks[7] || !slots[1].Weeks[17] {
		t.Errorf("多段周次解析错误: %+v", slots[1].Weeks)
	}

	// 双周
	slots = ParseSksj("星期三第6-7节{2-8周(双)}")
	if len(slots) != 1 {
		t.Fatalf("期望1个时间段，得到%d", len(slots))
	}
	for w, want := range map[int]bool{2: true, 3: false, 4: true, 7: false, 8: true} {
		if slots[0].Weeks[w] != want {
			t.Errorf("双周解析错误: 第%d周应为%v", w, want)
		}
	}

	// 单节次
	slots = ParseSksj("星期一第5节{3周}")
	if len(slots) != 1 || slots[0].P1 != 5 || slots[0].P2 != 5 || !slots[0].Weeks[3] {
		t.Errorf("单节次解析错误: %+v", slots)
	}
}

func TestCourseCodeOfJxbmc(t *testing.T) {
	if code := CourseCodeOfJxbmc("(2026-2027-1)-A0600910-02"); code != "A0600910" {
		t.Errorf("课程号提取错误: %s", code)
	}
}

func TestSortCoursePairs(t *testing.T) {
	idx := CourseIndex{
		"(2026-2027-1)-DROP001-01": {Kcmc: "退课A", Sksj: "星期二第3-5节{1-17周}"},
		"(2026-2027-1)-CONF001-01": {Kcmc: "冲突选课", Sksj: "星期二第4节{1-17周}"},
		"(2026-2027-1)-FREE001-01": {Kcmc: "无冲突选课", Sksj: "星期五第1-2节{1-17周}"},
	}
	pairs := [][]string{
		{"(2026-2027-1)-CONF001-01", "1"}, // 与退课A时间冲突
		{"(2026-2027-1)-DROP001-01", "0"},
		{"(2026-2027-1)-FREE001-01", "1"},
	}
	sorted := SortCoursePairs(pairs, idx)
	order := map[string]int{}
	for i, p := range sorted {
		order[p[0]] = i
	}
	// 无冲突的选课应最先执行
	if order["(2026-2027-1)-FREE001-01"] != 0 {
		t.Errorf("无冲突选课应排在最前: %v", sorted)
	}
	// 退课必须先于与其冲突的选课
	if order["(2026-2027-1)-DROP001-01"] > order["(2026-2027-1)-CONF001-01"] {
		t.Errorf("退课应先于冲突的选课: %v", sorted)
	}
}

func TestSortCoursePairsSameCode(t *testing.T) {
	// 同课程号不同教学班：必须先退后选，即使时间不冲突
	idx := CourseIndex{
		"(2026-2027-1)-A6500043-42": {Kcmc: "形势与政策3", Sksj: "星期三第6-7节{14-17周}"},
		"(2026-2027-1)-A6500043-02": {Kcmc: "形势与政策3", Sksj: "星期三第8-9节{10-13周}"},
	}
	pairs := [][]string{
		{"(2026-2027-1)-A6500043-02", "1"},
		{"(2026-2027-1)-A6500043-42", "0"},
	}
	sorted := SortCoursePairs(pairs, idx)
	if sorted[0][0] != "(2026-2027-1)-A6500043-42" {
		t.Errorf("同课程号应先退课再选课: %v", sorted)
	}
}

func TestSortCoursePairsStable(t *testing.T) {
	// 无冲突时保持原顺序
	pairs := [][]string{
		{"(2026-2027-1)-AAA0001-01", "1"},
		{"(2026-2027-1)-BBB0001-01", "1"},
		{"(2026-2027-1)-CCC0001-01", "1"},
	}
	sorted := SortCoursePairs(pairs, CourseIndex{})
	for i := range pairs {
		if sorted[i][0] != pairs[i][0] {
			t.Errorf("无冲突时应保持原顺序: %v", sorted)
		}
	}
}

func TestSortCoursePairsCategoryPriority(t *testing.T) {
	// 体育分项与通识选修课通常竞争最激烈，应优先抢
	idx := CourseIndex{
		"(2026-2027-1)-AAA0001-01": {Kcmc: "主修课", Sksj: "星期一第1-2节{1-17周}", Kklxmc: "主修课程"},
		"(2026-2027-1)-BBB0001-01": {Kcmc: "通识课", Sksj: "星期二第1-2节{1-17周}", Kklxmc: "通识选修课"},
		"(2026-2027-1)-CCC0001-01": {Kcmc: "体育课", Sksj: "星期三第1-2节{1-17周}", Kklxmc: "体育分项"},
	}
	pairs := [][]string{
		{"(2026-2027-1)-AAA0001-01", "1"},
		{"(2026-2027-1)-BBB0001-01", "1"},
		{"(2026-2027-1)-CCC0001-01", "1"},
	}
	sorted := SortCoursePairs(pairs, idx)
	want := []string{"(2026-2027-1)-CCC0001-01", "(2026-2027-1)-BBB0001-01", "(2026-2027-1)-AAA0001-01"}
	for i, w := range want {
		if sorted[i][0] != w {
			t.Errorf("开课类型优先级错误: %v", sorted)
			break
		}
	}
}
