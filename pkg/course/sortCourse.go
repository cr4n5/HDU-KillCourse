package course

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cr4n5/HDU-KillCourse/client"
	"github.com/cr4n5/HDU-KillCourse/config"
	"github.com/cr4n5/HDU-KillCourse/log"
	"github.com/iancoleman/orderedmap"
)

// TimeSlot 一段上课时间：星期几、第几节到第几节、哪些周
type TimeSlot struct {
	Day    int // 1=星期一 ... 7=星期日
	P1, P2 int // 节次范围
	Weeks  map[int]bool
}

var (
	sksjRe    = regexp.MustCompile(`星期(.)第(\d+)(?:-(\d+))?节\{([^}]*)\}`)
	weekSegRe = regexp.MustCompile(`^(\d+)(?:-(\d+))?周?(?:\((单|双)\))?$`)
	dayMap    = map[string]int{"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "日": 7, "天": 7}
)

// ParseSksj 解析上课时间字符串，如 星期二第3-5节{1-17周};星期四第6-7节{2-8周(双)}
func ParseSksj(sksj string) []TimeSlot {
	var slots []TimeSlot
	for _, part := range strings.Split(sksj, ";") {
		m := sksjRe.FindStringSubmatch(part)
		if m == nil {
			continue
		}
		day, ok := dayMap[m[1]]
		if !ok {
			continue
		}
		p1, _ := strconv.Atoi(m[2])
		p2 := p1
		if m[3] != "" {
			p2, _ = strconv.Atoi(m[3])
		}
		slots = append(slots, TimeSlot{Day: day, P1: p1, P2: p2, Weeks: parseWeeks(m[4])})
	}
	return slots
}

// parseWeeks 解析周次描述，如 "1-17周"、"1-17周(单)"、"1-3周,5-17周"
func parseWeeks(s string) map[int]bool {
	weeks := make(map[int]bool)
	for _, seg := range strings.Split(s, ",") {
		wm := weekSegRe.FindStringSubmatch(strings.TrimSpace(seg))
		if wm == nil {
			continue
		}
		w1, _ := strconv.Atoi(wm[1])
		w2 := w1
		if wm[2] != "" {
			w2, _ = strconv.Atoi(wm[2])
		}
		for w := w1; w <= w2; w++ {
			if wm[3] == "单" && w%2 == 0 {
				continue
			}
			if wm[3] == "双" && w%2 == 1 {
				continue
			}
			weeks[w] = true
		}
	}
	return weeks
}

// slotsOverlap 两段时间是否冲突：同一天、节次相交、周次相交
func slotsOverlap(a, b []TimeSlot) bool {
	for _, s1 := range a {
		for _, s2 := range b {
			if s1.Day != s2.Day || s1.P1 > s2.P2 || s2.P1 > s1.P2 {
				continue
			}
			for w := range s1.Weeks {
				if s2.Weeks[w] {
					return true
				}
			}
		}
	}
	return false
}

// sortNode 参与排序的一次选退课操作
type sortNode struct {
	Key     string // 教学班名称
	Value   string // "1"选课 "0"退课
	code    string
	slots   []TimeSlot
	catRank int // 课程类型优先级，越小越先抢
}

// categoryRank 按开课类型估计竞争激烈程度：体育分项与通识选修课通常最抢手
func categoryRank(kklxmc string) int {
	switch kklxmc {
	case "体育分项":
		return 0
	case "通识选修课":
		return 1
	default:
		return 2
	}
}

// coursesConflict 两门课是否冲突：同课程号（教务系统不允许同时持有同一课程的两个教学班）或上课时间重叠
func coursesConflict(a, b *sortNode) bool {
	if a.code == b.code {
		return true
	}
	return slotsOverlap(a.slots, b.slots)
}

// SortCoursePairs 对 [教学班名称, 值] 列表做拓扑排序：
// 退课必须先于与其冲突的选课执行；在满足约束的前提下尽量让选课靠前（先抢课），
// 选课之间按开课类型优先（体育分项 > 通识选修课 > 其他），
// 其余保持原有相对顺序稳定（作为用户手动排序的次序参考）。
func SortCoursePairs(pairs [][]string, idx CourseIndex) [][]string {
	n := len(pairs)
	nodes := make([]*sortNode, n)
	for i, p := range pairs {
		node := &sortNode{Key: p[0], Value: p[1], code: CourseCodeOfJxbmc(p[0]), catRank: 2}
		if info, ok := idx[p[0]]; ok {
			node.slots = ParseSksj(info.Sksj)
			node.catRank = categoryRank(info.Kklxmc)
		} else {
			log.Info("未在本地课程信息中找到课程(仅按课程号判断冲突): ", p[0])
		}
		nodes[i] = node
	}

	// 建图：退课 -> 与之冲突的选课
	adj := make([][]int, n)
	indeg := make([]int, n)
	for i, d := range nodes {
		if d.Value != "0" {
			continue
		}
		for j, e := range nodes {
			if e.Value == "1" && coursesConflict(d, e) {
				adj[i] = append(adj[i], j)
				indeg[j]++
			}
		}
	}

	// Kahn拓扑排序：可执行的操作中优先选课(按开课类型优先级，其次原顺序)，再退课(保持原顺序)
	visited := make([]bool, n)
	result := make([][]string, 0, n)
	for len(result) < n {
		pick := -1
		for _, wantValue := range []string{"1", "0"} {
			for i := 0; i < n; i++ {
				if !visited[i] && indeg[i] == 0 && nodes[i].Value == wantValue {
					if pick < 0 || (wantValue == "1" && nodes[i].catRank < nodes[pick].catRank) {
						pick = i
					}
				}
			}
			if pick >= 0 {
				break
			}
		}
		if pick < 0 {
			// 理论上不可能出现环（边只从退课指向选课），保险起见取第一个未访问的
			for i := 0; i < n; i++ {
				if !visited[i] {
					pick = i
					break
				}
			}
		}
		visited[pick] = true
		result = append(result, []string{nodes[pick].Key, nodes[pick].Value})
		for _, j := range adj[pick] {
			indeg[j]--
		}
	}
	return result
}

// AutoSortCourses 自动排序配置文件中的课程并保存
func AutoSortCourses(cfg *config.Config, courses *client.GetCourseResp) error {
	// 构建课程索引：course.json + Excel补充
	idx := indexFromResp(courses)
	if excelIdx, err := indexFromExcel(cfg); err == nil {
		for k, v := range excelIdx {
			idx[k] = v
		}
	}

	// OrderedMap -> pairs
	var pairs [][]string
	for _, k := range cfg.Course.Keys() {
		v, _ := cfg.Course.Get(k)
		pairs = append(pairs, []string{k, v.(string)})
	}

	sorted := SortCoursePairs(pairs, idx)

	// 是否有变化
	changed := false
	for i := range sorted {
		if sorted[i][0] != pairs[i][0] {
			changed = true
			break
		}
	}
	if !changed {
		log.Info("课程自动排序: 顺序已满足约束，无需调整")
		return nil
	}

	// pairs -> OrderedMap
	omap := orderedmap.New()
	for _, p := range sorted {
		omap.Set(p[0], p[1])
	}
	cfg.Course = omap

	log.Info("课程自动排序完成，新的执行顺序:")
	for _, p := range sorted {
		op := "选课"
		if p[1] == "0" {
			op = "退课"
		}
		name := ""
		if info, ok := idx[p[0]]; ok {
			name = info.Kcmc
		}
		log.Info("  ", op, "  ", p[0], "  ", name)
	}

	return config.SaveConfig(cfg)
}
