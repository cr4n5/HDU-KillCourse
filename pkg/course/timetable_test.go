package course

import (
	"testing"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/iancoleman/orderedmap"
)

func TestParseWeekday(t *testing.T) {
	cases := []struct {
		xqj, xqjmc string
		want       int
	}{
		{"2", "", 2},
		{"", "星期二", 2},
		{"", "星期日", 7},
		{"", "星期天", 7},
		{"7", "", 7},
		{"", "无效", 0},
	}
	for _, c := range cases {
		if got := parseWeekday(c.xqj, c.xqjmc); got != c.want {
			t.Errorf("parseWeekday(%q,%q)=%d, want %d", c.xqj, c.xqjmc, got, c.want)
		}
	}
}

func TestParsePeriods(t *testing.T) {
	cases := []struct {
		in     string
		p1, p2 int
	}{
		{"3-5", 3, 5},
		{"6", 6, 6},
		{"10-12", 10, 12},
		{"", 0, 0},
	}
	for _, c := range cases {
		if p1, p2 := parsePeriods(c.in); p1 != c.p1 || p2 != c.p2 {
			t.Errorf("parsePeriods(%q)=%d,%d want %d,%d", c.in, p1, p2, c.p1, c.p2)
		}
	}
}

func TestParseWeeksParity(t *testing.T) {
	// 单周
	w := parseWeeks("1-17周(单)")
	for _, odd := range []int{1, 3, 17} {
		if !w[odd] {
			t.Errorf("单周应包含第%d周", odd)
		}
	}
	for _, even := range []int{2, 4, 16} {
		if w[even] {
			t.Errorf("单周不应包含第%d周", even)
		}
	}
	// 多段
	w = parseWeeks("1-3周,5-8周")
	for _, in := range []int{1, 2, 3, 5, 8} {
		if !w[in] {
			t.Errorf("应包含第%d周", in)
		}
	}
	if w[4] {
		t.Errorf("不应包含第4周")
	}
}

func TestScheduleItemsToEvents(t *testing.T) {
	cfg := &config.Config{Course: orderedmap.New()}
	cfg.Time = config.Time{XueNian: "2026", XueQi: "1"}
	cfg.SemesterStartDate = "2026-09-07" // 星期一

	items := []client.ScheduleItem{
		{Kcmc: "高等数学", Xqj: "1", Jcs: "1-2", Zcd: "1-17周", Cdmc: "1教101", Xm: "张三", Jxbmc: "(2026-2027-1)-X0000001-01"},
		{Kcmc: "物理", Xqjmc: "星期三", Jcs: "3-4", Zcd: "1-17周(单)", Cdmc: "2教202", Xm: "李四"},
	}
	events, err := scheduleItemsToEvents(items, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// 高数17次(全周) + 物理9次(单周1,3,...,17) = 26
	if len(events) != 17+9 {
		t.Fatalf("事件数=%d, 期望26", len(events))
	}
	// 第1周星期一高数应为 2026-09-07 08:05
	first := events[0]
	if first.start.Format("2006-01-02 15:04") != "2026-09-07 08:05" {
		t.Errorf("首个事件时间错误: %s", first.start.Format("2006-01-02 15:04"))
	}
	if first.summary != "高等数学" || first.location != "1教101" {
		t.Errorf("首个事件内容错误: %+v", first)
	}
}
