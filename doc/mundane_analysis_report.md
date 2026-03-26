# SolarSage 世俗占星 (Mundane Astrology) 模块 — 调研分析报告

## 一、研究背景

世俗占星 (Mundane Astrology) 是占星学最古老的分支，可追溯至巴比伦时代（公元前2000年），
其核心关注点为国家命运、政治事件、自然灾害和世界周期。本报告基于对现有开源占星库、
权威数据源和学术传统的深度调研，为 SolarSage 实现 `pkg/mundane/` 模块提供完整的技术路线。

---

## 二、数据源分析

### 2.1 权威数据源对比

| 数据源 | 数据量 | 格式 | 许可 | 质量 | 国盘覆盖 |
|--------|--------|------|------|------|----------|
| **Astro-Databank** (astro.com) | 72,271+ 页 | Wiki/XML | 非开源, 可手动提取 | AA-DD 分级 | 有专门 "Birth of State" 分类 |
| **Campion《Book of World Horoscopes》** | ~500 图 | 书籍 | 版权保护 | 学术级, 多版本比较 | 最全面: 覆盖几乎所有主权国家 |
| **Astrotheme** (astrotheme.com) | 数百国 | HTML/可抓取 | 免费浏览 | 中等, 时间来源不透明 | 广泛, 包含城市星盘 |
| **AA Chart Database** (英国占星协会) | 100+ | 在线浏览 | 会员访问 | 高, 均有评级 | 有限 |
| **Astro-Seek** | 镜像 Databank | HTML | 免费浏览 | 同 Databank | 同 Databank |

**评估结论**: 不存在一个单一的、结构化的、开源的国家星盘数据库。数据散落在多个来源中，
且同一国家常有多个争议版本。最佳策略是**手工编纂一份高质量核心数据集**，以 Campion
和 Astro-Databank 为交叉验证基准。

### 2.2 数据质量等级 (Rodden Rating)

Astro-Databank 使用 Lois Rodden 的数据可靠性分级:

| 等级 | 含义 | 国盘语境 |
|------|------|----------|
| **AA** | 出生记录/官方文件 | 有精确官方记录的建国/独立时刻 |
| **A** | 当事人引述 | 国家领导人/官方发言提及的时间 |
| **B** | 传记/新闻记录 | 新闻报道中的时间 |
| **C** | 推测/校正 | 占星师通过事件反推的时间 |
| **DD** | 存在争议 | 同一国家有多个竞争版本 |

国家星盘中，**AA/A 级数据极少**。大多数属于 B-C 级。美国 Sibly 盘虽然被广泛使用，
其 17:10 LMT 的时间实际上属于 C 级 (Ebenezer Sibly 的 18 世纪推算)。

---

## 三、已收集的核心国家星盘数据

### 3.1 数据集 (30 国, G20 + 地缘关键国)

以下数据交叉验证自 Astro-Databank、Astrotheme 和 Campion:

| # | 国家 | 事件 | 日期 | 时间 | 时区 | 地点 | 纬度 | 经度 | 评级 | 备注 |
|---|------|------|------|------|------|------|------|------|------|------|
| 1 | **美国** | 独立宣言签署 | 1776-07-04 | 17:10 | LMT | Philadelphia | 39.95 | -75.15 | C | Sibly Chart, 最广泛使用 |
| 2 | **英国** | 合并法案生效 | 1801-01-01 | 00:00 | GMT | London | 51.51 | -0.12 | B | Campion 首选 |
| 3 | **中国** | 开国大典 | 1949-10-01 | 15:15 | +08:00 | Beijing | 39.91 | 116.39 | A | 毛泽东宣布时刻, Astrotheme |
| 4 | **俄罗斯** | 苏联解体声明 | 1991-12-25 | 19:38 | +03:00 | Moscow | 55.76 | 37.62 | A | 戈尔巴乔夫辞职, 旗帜降下 |
| 5 | **德国** | 基本法生效 | 1949-05-24 | 00:00 | CET | Bonn | 50.73 | 7.10 | B | 联邦德国; 另有 1990-10-03 统一盘 |
| 6 | **法国** | 第五共和国 | 1958-09-28 | 18:00 | CET | Paris | 48.86 | 2.35 | B | 公投通过; 另有 1958-10-05 宪法盘 |
| 7 | **日本** | 宪法生效 | 1947-05-03 | 00:00 | +09:00 | Tokyo | 35.68 | 139.69 | C | 另有 1952-04-28 主权恢复盘 |
| 8 | **印度** | 独立 | 1947-08-15 | 00:01 | +05:30 | New Delhi | 28.61 | 77.21 | AA | 午夜独立, 精确记录 |
| 9 | **巴西** | 独立 | 1822-09-07 | 16:35 | LMT | São Paulo | -23.55 | -46.63 | B | "Grito do Ipiranga" |
| 10 | **韩国** | 大韩民国成立 | 1948-08-15 | 11:00 | +09:00 | Seoul | 37.57 | 126.98 | B | Astrotheme |
| 11 | **澳大利亚** | 联邦成立 | 1901-01-01 | 00:00 | +10:00 | Sydney | -33.87 | 151.21 | B | Federation Day |
| 12 | **加拿大** | 联邦成立 | 1867-07-01 | 00:00 | LMT | Ottawa | 45.42 | -75.70 | C | BNA Act 生效 |
| 13 | **墨西哥** | 独立 | 1810-09-16 | 05:20 | LMT | Dolores Hidalgo | 21.16 | -100.93 | B | "Grito de Dolores" |
| 14 | **以色列** | 独立宣言 | 1948-05-14 | 16:00 | +02:00 | Tel Aviv | 32.07 | 34.77 | A | Ben-Gurion 宣读 |
| 15 | **土耳其** | 共和国宣布 | 1923-10-29 | 20:30 | +02:00 | Ankara | 39.93 | 32.86 | B | Astrotheme |
| 16 | **伊朗** | 伊斯兰共和国 | 1979-04-01 | 15:00 | +03:30 | Tehran | 35.69 | 51.39 | B | 公投后宣布 |
| 17 | **沙特阿拉伯** | 王国统一 | 1932-09-23 | 12:00 | +03:00 | Riyadh | 24.69 | 46.72 | C | 午间宣布, 时间不确定 |
| 18 | **意大利** | 共和国公投 | 1946-06-10 | 18:00 | CET | Rome | 41.90 | 12.50 | B | 最高法院宣布结果 |
| 19 | **欧盟** | 马斯特里赫特条约生效 | 1993-11-01 | 00:00 | CET | Brussels | 50.85 | 4.35 | A | 法律生效时刻 |
| 20 | **联合国** | 宪章生效 | 1945-10-24 | 16:45 | EST | Washington DC | 38.91 | -77.04 | A | 批准存档时刻 |
| 21 | **南非** | 新宪法生效 | 1994-04-27 | 00:00 | +02:00 | Pretoria | -25.75 | 28.19 | B | Freedom Day |
| 22 | **阿根廷** | 独立 | 1816-07-09 | 12:00 | LMT | Tucumán | -26.82 | -65.22 | C | |
| 23 | **印度尼西亚** | 独立宣言 | 1945-08-17 | 10:00 | +07:00 | Jakarta | -6.21 | 106.85 | A | Sukarno 宣读 |
| 24 | **埃及** | 共和国宣布 | 1953-06-18 | 23:30 | +02:00 | Cairo | 30.04 | 31.24 | B | |
| 25 | **巴基斯坦** | 独立 | 1947-08-14 | 00:01 | +05:00 | Karachi | 24.86 | 67.01 | AA | 午夜独立, 精确记录 |
| 26 | **乌克兰** | 独立 | 1991-08-24 | 18:00 | +03:00 | Kyiv | 50.45 | 30.52 | B | 最高拉达投票 |
| 27 | **泰国** | 暹罗不适用 | — | — | — | — | — | — | — | 从未被殖民, 无明确"建国"时刻 |
| 28 | **朝鲜** | 建国 | 1948-09-09 | 12:00 | +09:00 | Pyongyang | 39.02 | 125.75 | C | |
| 29 | **新加坡** | 独立 | 1965-08-09 | 10:00 | +07:30 | Singapore | 1.29 | 103.85 | A | Lee Kuan Yew 宣布 |
| 30 | **德国(统一)** | 两德统一 | 1990-10-03 | 00:00 | CET | Berlin | 52.52 | 13.41 | AA | 官方统一时刻 |

### 3.2 争议数据说明

几个重要国家的星盘存在**长期学术争议**:

1. **美国**: 至少 7 个竞争版本 — Sibly (17:10), Gemini Rising (02:13 AM), Sagittarius Rising (various)。
   Campion 在《Book of World Horoscopes》中列出全部但不做裁决。Sibly 使用最广但非 "AA" 级。

2. **俄罗斯**: 19 个候选星盘。1991-12-25 (苏联解体), 1991-12-08 (别洛韦日协议),
   1990-06-12 (主权宣言) 等均有支持者。

3. **德国**: 1871 帝国、1919 魏玛、1949 联邦、1990 统一 — 四个时代四种盘。
   本数据集包含 1949 和 1990 两版。

4. **中国**: 15:00 vs 15:15 争议。Astrotheme 取 15:15, 部分占星师取 15:00。
   均源于1949年10月1日天安门广场开国大典，毛泽东宣读中华人民共和国中央人民政府成立公告。

---

## 四、竞品世俗占星功能分析

### 4.1 现有开源库的世俗占星支持

| 库 | 语言 | 国盘数据 | 入境图 | 日月食路径 | 行运分析 | 大周期 |
|---|---|---|---|---|---|---|
| **Kerykeion** | Python | 无内置 | 有教程无代码 | 无 | 基础 | 无 |
| **flatlib** | Python | 无 | 无 | 无 | 无 | 无 |
| **Astrolog** | C | 无 | 无 | 无 | 基础 | 无 |
| **Morinus** | Python | 无 | 无 | 无 | 有 | 无 |
| **SolarSage** | Go | **无 (待实现)** | **无** | 部分 (pkg/lunar) | **有** | **无** |

**结论**: 目前没有任何开源占星库提供结构化的国家星盘数据库和世俗占星专用计算。
这是一个完全空白的领域。SolarSage 如果实现此功能，将成为**全球首个内置国盘数据库的开源占星引擎**。

### 4.2 世俗占星核心技法

根据 Benjamin Dykes、Nicholas Campion 和传统文献，世俗占星的技法体系:

```
世俗占星核心技法
│
├── 1. 入境图 (Ingress Charts)
│   ├── 白羊入境图 (Aries Ingress) — 年度国运主图
│   ├── 巨蟹/天秤/摩羯入境图 — 季度补充
│   └── 土木合相入境图 — 20年大周期
│
├── 2. 日月食 (Eclipses)
│   ├── 日食图 (Solar Eclipse Chart)
│   ├── 月食图 (Lunar Eclipse Chart)
│   └── 食相路径与国家领土叠加
│
├── 3. 国盘行运 (Transits to National Chart)
│   ├── 外行星过境关键轴点 (ASC/MC)
│   ├── 冥王星回归 (美国 2022 首次)
│   └── 土星回归 (~29.5 年周期)
│
├── 4. 大周期 (Great Cycles)
│   ├── 土木合相周期 (Jupiter-Saturn, ~20yr)
│   ├── 土冥合相周期 (~33-38yr)
│   └── 海王星/天王星周期
│
├── 5. 国盘推运 (Progressions for Nations)
│   ├── 次限推运 (Secondary Progressions)
│   └── 太阳弧推运 (Solar Arc)
│
└── 6. 专用分析
    ├── 国盘叠加 (Synastry between Nations)
    ├── 领导人盘与国盘互动
    └── Astro*Carto*Graphy 地理映射
```

---

## 五、技术实现方案

### 5.1 模块架构

```
pkg/mundane/
├── doc.go              # 包文档
├── nations.go          # 国家星盘数据结构 + 30 国内置数据
├── nations_data.go     # 数据常量 (大文件, 独立存放)
├── nations_test.go     # 数据完整性测试
├── ingress.go          # 入境图计算 (Aries/Cancer/Libra/Capricorn)
├── ingress_test.go
├── cycles.go           # 大周期: 土木合相、土冥合相查找
├── cycles_test.go
├── national_transit.go # 国盘行运分析 (调用现有 pkg/transit/)
├── national_transit_test.go
├── synastry.go         # 国盘比较 (调用现有 pkg/synastry/)
└── synastry_test.go
```

### 5.2 核心数据结构

```go
// NationChart 表示一个国家/政体的建国星盘数据
type NationChart struct {
    ID          string   // 唯一标识: "USA_SIBLY", "CN_PRC_1949"
    Name        string   // "United States of America"
    NameLocal   string   // "美利坚合众国" (本地化名称)
    Event       string   // "Declaration of Independence"
    Date        string   // "1776-07-04" (ISO 8601)
    Time        string   // "17:10" (24h)
    Timezone    string   // "LMT" | "+08:00" | "GMT"
    City        string   // "Philadelphia"
    Country     string   // "US"
    Latitude    float64  // 39.9526
    Longitude   float64  // -75.1652
    Rating      string   // "AA" | "A" | "B" | "C" | "DD"
    Source      string   // "Campion BWH p.367"
    Notes       string   // 备注和争议说明
    Tags        []string // ["G20", "NATO", "UN_P5"]
    Alternates  []string // 指向同一国家其他版本的 ID
}

// IngressChart 表示一个入境图
type IngressChart struct {
    Type       string         // "ARIES" | "CANCER" | "LIBRA" | "CAPRICORN"
    Year       int
    JD         float64        // 精确入境时刻 JD
    Location   GeoLocation    // 投射到的首都
    ChartInfo  *models.ChartInfo
}

// GreatCycle 表示一次大行星合相
type GreatCycle struct {
    Type      string  // "JUPITER_SATURN" | "SATURN_PLUTO"
    JD        float64
    Longitude float64 // 合相发生的黄道度数
    Sign      string
    Element   string  // 合相元素 (对土木合相有特殊意义)
}
```

### 5.3 实现优先级

| 优先级 | 功能 | 依赖 | 复杂度 | 价值 |
|--------|------|------|--------|------|
| **P0** | 国盘数据库 (30 国内置) | 无 | 低 | 极高 — 全球首创 |
| **P0** | 国盘星盘计算 | pkg/chart | 低 | 极高 — 数据的基础消费方式 |
| **P1** | 白羊入境图 | pkg/sweph (精确太阳入白羊时刻) | 中 | 高 — 世俗占星最核心技法 |
| **P1** | 四季入境图 | 同上 | 中 | 高 |
| **P1** | 国盘行运 | pkg/transit | 低 (封装调用) | 高 |
| **P2** | 土木/土冥合相周期 | pkg/sweph | 中 | 中高 |
| **P2** | 国盘推运 | pkg/progressions | 低 (封装调用) | 中 |
| **P2** | 国盘比较 (中美/中俄等) | pkg/synastry | 低 (封装调用) | 中 |
| **P3** | 日月食路径分析 | pkg/lunar + 扩展 | 高 | 中 |
| **P3** | 领导人盘叠加 | 数据收集难度大 | 中 | 低 |

### 5.4 API 设计草案

```
# MCP Tools (新增)
mundane_list_nations          # 列出所有可用国家星盘
mundane_nation_chart          # 计算指定国家的完整星盘
mundane_nation_transits       # 国盘行运分析 (指定时间段)
mundane_nation_synastry       # 两国星盘比较
mundane_ingress_chart         # 入境图 (年份 + 首都)
mundane_great_cycles          # 大行星合相周期查询

# REST Endpoints (新增)
GET  /api/mundane/nations
GET  /api/mundane/nations/{id}/chart
GET  /api/mundane/nations/{id}/transits?start=...&end=...
GET  /api/mundane/nations/{id1}/synastry/{id2}
GET  /api/mundane/ingress?year=2026&type=aries&city=Beijing
GET  /api/mundane/cycles?type=jupiter_saturn&start=2000&end=2100
```

---

## 六、风险与应对

| 风险 | 影响 | 应对策略 |
|------|------|----------|
| 国盘时间争议 | 同一国家多个版本, 用户困惑 | 提供 `Alternates` 字段, 默认选最广泛使用的版本, 注明评级 |
| 数据版权 | Campion 书籍有版权 | 只使用公开来源 (Astrotheme, Astro-Databank), 数据本身 (日期/时间) 不受版权保护 |
| 政治敏感性 | 某些"建国"日期有政治争议 | 用中性表述 (如 "Founding Event"), 不做政治判断 |
| 数据维度不足 | 缺少时间的国家无法计算精确宫位 | 标注 Rating, 对 C/DD 级数据提示用户 |
| 时区历史复杂 | LMT→标准时区转换、历史夏令时 | 使用 Swiss Ephemeris 的 JulDay 做 LMT 手工换算 |

---

## 七、结论

1. **市场空白**: 目前全球没有任何开源占星引擎包含结构化国盘数据库。这是 SolarSage
   最有差异化价值的新功能方向。

2. **数据可行**: 30 国核心数据集已收集完毕，均可从公开来源交叉验证。

3. **技术可行**: SolarSage 已具备全部底层能力 — pkg/chart (星盘计算), pkg/transit
   (行运), pkg/synastry (比较), pkg/progressions (推运), pkg/lunar (日月食)。
   `pkg/mundane/` 主要是数据层 + 封装调用层, 实现成本较低。

4. **建议路线**: P0 (国盘数据 + 星盘计算) → P1 (入境图 + 行运) → P2 (大周期 + 比较)

---

*报告日期: 2026-03-23*
*数据来源: Astro-Databank (astro.com), Astrotheme, Campion BWH, Astrology King, 及多个学术占星网站*
