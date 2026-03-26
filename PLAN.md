# Transit 模块重构规划

## 一、现状问题诊断

### 核心病灶：`CalcTransitEvents` 承担了三种职责

```
现状：一个函数做三件事
  1. 配置解析（哪些组合启用？）
  2. 任务枚举（对哪些行星、哪些参考点？）
  3. 计算执行（调用 RQ1/RQ2 扫描器）
```

**症状表现：**
- 三层 `for` 嵌套，逻辑重复（TR/SP/SA 三段几乎平行）
- 站点被重复计算（TR-NA 用一次，TR-TR 又重算）
- 新增一种组合类型（如 SP-SA）需要侵入主函数

---

## 二、你的代码已经隐含的体系

你的代码结构非常健康，核心抽象已经存在，只是**没有被显式命名**：

| 隐含概念 | 当前代码中的体现 | 缺少什么 |
|---------|---------------|---------|
| **移动天体** | `makeXxxCalcFn` 工厂函数族 | 没有统一的结构体承载 |
| **RQ1/RQ2 分类** | 函数名已区分 | 没有对应的任务类型 |
| **固定参考点** | `natalRefs` 局部变量 | 没有独立类型 |
| **站点缓存** | 每段逻辑各自计算 | 没有共享缓存 |
| **组合启用判断** | `if IncludeTrNa` 散落各处 | 没有集中的映射表 |

**结论：重构不是推翻，是把已有的隐式结构显式化。**

---

## 三、重构目标

```
将 CalcTransitEvents 的职责拆分为：

  配置 ──→ 任务列表 ──→ 执行 ──→ 合并排序
  (声明)    (枚举)      (计算)    (后处理)
```

四个步骤，每个步骤独立、可测试。

---

## 四、三个新概念（最小化新增类型）

### 概念一：`MovingBody`（移动天体）

**作用：** 统一描述"一个在时间范围内运动的计算天体"，替代当前散落的 `calcFn + planet + chartType + orbConfig` 四元组。

**包含字段：**
- 天体 ID、盘面类型、计算函数、轨道配置
- `CanRetrograde bool`（控制是否做站点检测，替代当前 `if planet == Sun/Moon` 的散落判断）

---

### 概念二：`NatalRef`（出生参考点）

**作用：** 统一描述"一个固定的出生盘参考位置"，替代当前的匿名 `refPoint` 结构体。

**包含字段：**
- ID（用于事件的 Target 字段）、经度、盘面类型

---

### 概念三：`CalcContext`（计算上下文）

**作用：** 持有所有预计算的只读数据，避免重复计算。

**包含字段：**
- `NatalHouses`、`NatalRefs`（原本是局部变量）
- `StationCache map[PlanetID][]StationInfo`（新增：站点缓存）

---

## 五、文件职责划分

```
transit/
├── calc.go       # 主入口（精简为调度器，~30行）
├── context.go    # CalcContext 定义 + 预计算逻辑
├── body.go       # MovingBody + 三个工厂函数
├── tasks.go      # buildTasks：声明式枚举所有任务
│
│   ── 以下文件基本不动，只做整理 ──
├── scanner.go    # RQ1/RQ2 扫描（现有逻辑）
├── station.go    # 站点检测（现有逻辑）
├── ingress.go    # 过座/过宫（现有逻辑）
├── voc.go        # 虚空月亮（现有逻辑）
└── bisect.go     # 数学工具（现有逻辑）
```

**新增代码量估计：** ~150行（全部是结构定义和胶水代码）  
**删除/替换代码量：** CalcTransitEvents 函数体约 180行 → 约 25行

---

## 六、重构后的主函数形态

重构后 `CalcTransitEvents` 应该呈现这个形态：

```
func CalcTransitEvents(input):
  1. ctx = buildCalcContext(input)    // 预计算固定数据
  2. tasks = buildTasks(ctx)          // 声明式生成任务列表
  3. events = runAll(tasks, ctx)      // 统一执行
  4. events += findVoidOfCourse(...)  // 后处理
  5. sort(events)
  6. return events
```

---

## 七、`buildTasks` 的逻辑结构

这是重构的核心，替代原来的三段平行嵌套循环：

```
buildTasks:
  
  收集三类天体列表：
    transitBodies    = buildTransitBodies(input)
    progressBodies   = buildProgressionBodies(input)
    solarArcBodies   = buildSolarArcBodies(input)

  对 transitBodies 中每个 trBody：
    生成：StationTask、SignIngressTask、HouseIngressTask
    若启用 TR-NA：对每个 NatalRef 生成 AspectRQ1Task
    若启用 TR-TR：对每对 transit 天体生成 AspectRQ2Task（去重）
    若启用 TR-SP：对每个 progressBody 生成 AspectRQ2Task
    若启用 TR-SA：对每个 solarArcBody 生成 AspectRQ2Task

  对 progressBodies 中每个 spBody：
    生成：StationTask、SignIngressTask、HouseIngressTask
    若启用 SP-NA：对每个 NatalRef 生成 AspectRQ1Task
    若启用 SP-SP：对每对 progress 天体生成 AspectRQ2Task（去重）

  对 solarArcBodies 中每个 saBody：
    生成：SignIngressTask、HouseIngressTask（无站点）
    若启用 SA-NA：对每个 NatalRef 生成 AspectRQ1Task
```

---

## 八、关键设计决策

### 决策一：Task 接口 vs 直接调用

**选择：使用 `Task` 接口**

```
interface Task {
    Run(ctx) → []Event
}
```

**理由：**
- `buildTasks` 返回 `[]Task`，主函数无需关心具体类型
- 未来支持并发只需改 `runAll` 一处
- 各任务类型可独立单元测试

**任务类型清单（共4种，覆盖所有事件）：**

| Task 类型 | 对应事件 | 对应原函数 |
|----------|---------|----------|
| `StationTask` | 站点顺逆 | `makeStationEvent` |
| `SignIngressTask` | 过座 | `findSignIngressEvents` |
| `HouseIngressTask` | 过宫 | `findHouseIngressEvents` |
| `AspectRQ1Task` | TR-NA / SP-NA / SA-NA 相位 | `findAspectEventsRQ1` |
| `AspectRQ2Task` | TR-TR / TR-SP / TR-SA / SP-SP 相位 | `findAspectEventsRQ2` |

---

### 决策二：站点缓存放在 CalcContext

**问题：** 现在 TR-NA 会调用一次 `findStations(tPlanet)`，TR-TR 里又会重算同一颗行星的站点。

**方案：** `CalcContext` 提供 `GetStations(planet, calcFn)` 方法，内部做 `map` 缓存。

---

### 决策三：特殊点与行星统一为 MovingBody

**现状：** 特殊点（ASC/MC）在原代码中有单独的 `if SpecialPoints != nil` 分支。

**方案：** `buildTransitBodies` / `buildProgressionBodies` / `buildSolarArcBodies` 各自在内部处理特殊点，对外统一返回 `[]MovingBody`，主任务枚举循环无需区分。

---

### 决策四：`shouldPairRQ2` 去重规则

RQ2 组合 (A,B) 和 (B,A) 只需计算一次。

**规则：** `string(planet1) < string(planet2)` 时才生成任务（与原代码 `>=` 跳过逻辑一致，但集中到一处）。

---

## 九、不需要改动的部分

以下函数**逻辑正确，无需修改**，重构只是改变调用方式：

- `findAspectEventsRQ1` / `findAspectEventsRQ2`
- `findStations` / `buildMonoIntervals`
- `findSignIngressEvents` / `findHouseIngressEvents`
- `findVoidOfCourse`
- 所有 `bisect*` 函数
- 所有 `make*CalcFn` 工厂函数
- `adaptiveStep` / `angleDiffToAspect` 等数学工具

---

## 十、重构步骤（建议顺序）

```
Step 1  新建 context.go
        定义 CalcContext、NatalRef
        将 buildCalcContext、buildNatalRefs 从主函数中提取

Step 2  新建 body.go
        定义 MovingBody
        将三个 buildXxxBodies 工厂函数实现
        实现 canRetrograde 辅助函数

Step 3  新建 tasks.go
        定义 Task 接口和四种 Task 结构体
        实现 buildTasks
        各 Task.Run 内部直接调用现有扫描函数

Step 4  重写 calc.go 主函数
        替换为五行调度逻辑
        删除原有 180 行嵌套循环

Step 5  验证
        对比重构前后同一输入的事件列表
        确保事件数量、时间、类型完全一致
```

---

## 十一、预期收益

| 指标 | 重构前 | 重构后 |
|-----|-------|-------|
| `CalcTransitEvents` 行数 | ~200 行 | ~25 行 |
| 新增组合类型所需改动 | 侵入主函数 | 只加一个 Task + 在 buildTasks 里加几行 |
| 站点重复计算 | 是 | 否（缓存） |
| 可并发执行 | 否 | 改 `runAll` 一处即可 |
| 单元测试粒度 | 只能测整体 | 每个 Task 独立可测 |


# `CalcTransitEvents` 输入重构设计

## 一、现状问题

当前 `TransitCalcInput` 是一个**扁平的大结构体**：

```
TransitCalcInput
├── NatalLat/Lon/JD/Planets      ← 出生数据
├── TransitLat/Lon               ← 过运地点
├── StartJD/EndJD                ← 时间范围
├── TransitPlanets               ← 过运行星
├── ProgressionsConfig           ← 推运配置（含 Enabled 开关）
├── SolarArcConfig               ← 太阳弧配置（含 Enabled 开关）
├── SpecialPoints                ← 特殊点（四类混在一个结构体）
├── EventConfig                  ← 事件过滤
├── OrbConfigTransit             ← 三套轨道配置并列
├── OrbConfigProgressions
├── OrbConfigSolarArc
└── HouseSystem
```

**问题：**
- 三套 `OrbConfig` 平铺，与对应盘面分离
- `SpecialPoints` 把四类特殊点混在一起，与盘面配置分离
- `ProgressionsConfig` 和 `SolarArcConfig` 各自带 `Enabled`，但 `TransitPlanets` 没有对等结构

---

## 二、核心设计原则

**每种盘面的配置应该自包含**，即"这个盘面用哪些行星、哪些特殊点、什么轨道"应该聚合在一起。

---

## 三、重构后的输入结构

### 顶层结构

```
CalcTransitEvents 的输入 = 三块正交的关注点

TransitCalcInput
├── NatalChart      NatalChartConfig     ← 出生盘（固定，参考基准）
├── TimeRange       TimeRangeConfig      ← 计算时间范围
├── Charts          ChartSetConfig       ← 各盘面配置（TR/SP/SA）
├── EventFilter     EventFilterConfig    ← 启用哪些事件类型
└── HouseSystem     models.HouseSystem   ← 宫位系统（全局共享）
```

---

### 各子结构

**`NatalChartConfig`** — 出生盘配置（固定参考点来源）

```
NatalChartConfig
├── Lat, Lon    float64              ← 出生地
├── JD          float64              ← 出生时刻
├── Planets     []PlanetID           ← 参与计算的出生行星
└── Points      []SpecialPointID     ← 参与计算的出生特殊点
```

---

**`TimeRangeConfig`** — 时间范围

```
TimeRangeConfig
├── StartJD   float64
└── EndJD     float64
```

---

**`ChartSetConfig`** — 各盘面配置集合

```
ChartSetConfig
├── Transit      *TransitChartConfig       ← nil 表示不启用
├── Progressions *ProgressionsChartConfig  ← nil 表示不启用
└── SolarArc     *SolarArcChartConfig      ← nil 表示不启用
```

每个盘面配置**自包含**其所需的全部信息：

```
TransitChartConfig
├── Lat, Lon   float64           ← 过运地点（影响 ASC/MC 特殊点）
├── Planets    []PlanetID        ← 参与过运的行星
├── Points     []SpecialPointID  ← 参与过运的特殊点
└── Orbs       OrbConfig         ← 该盘面使用的轨道配置

ProgressionsChartConfig
├── Planets    []PlanetID
├── Points     []SpecialPointID
└── Orbs       OrbConfig

SolarArcChartConfig
├── Planets    []PlanetID
├── Points     []SpecialPointID
└── Orbs       OrbConfig
```

---

**`EventFilterConfig`** — 事件过滤（启用哪些事件类型）

```
EventFilterConfig
├── Station         bool   ← 站点事件
├── SignIngress      bool   ← 过座
├── HouseIngress     bool   ← 过宫
├── VoidOfCourse     bool   ← 虚空月亮
│
├── TrNa            bool   ← Transit → Natal 相位
├── TrTr            bool   ← Transit → Transit 相位
├── TrSp            bool   ← Transit → Progressions 相位
├── TrSa            bool   ← Transit → SolarArc 相位
├── SpNa            bool   ← Progressions → Natal 相位
├── SpSp            bool   ← Progressions → Progressions 相位
└── SaNa            bool   ← SolarArc → Natal 相位
```

---

## 四、与原结构的对照

| 原字段 | 重构后位置 |
|-------|----------|
| `NatalLat/Lon/JD` | `NatalChart.Lat/Lon/JD` |
| `NatalPlanets` | `NatalChart.Planets` |
| `SpecialPoints.NatalPoints` | `NatalChart.Points` |
| `TransitLat/Lon` | `Charts.Transit.Lat/Lon` |
| `TransitPlanets` | `Charts.Transit.Planets` |
| `SpecialPoints.TransitPoints` | `Charts.Transit.Points` |
| `OrbConfigTransit` | `Charts.Transit.Orbs` |
| `ProgressionsConfig.Planets` | `Charts.Progressions.Planets` |
| `SpecialPoints.ProgressionsPoints` | `Charts.Progressions.Points` |
| `OrbConfigProgressions` | `Charts.Progressions.Orbs` |
| `ProgressionsConfig.Enabled` | `Charts.Progressions == nil` |
| `SolarArcConfig.Enabled` | `Charts.SolarArc == nil` |
| `EventConfig.*` | `EventFilter.*` |
| `HouseSystem` | 顶层 `HouseSystem` |

---

## 五、关键改进点说明

### 改进一：用 `nil` 替代 `Enabled bool`

```
原来：ProgressionsConfig.Enabled == false → 跳过
现在：Charts.Progressions == nil          → 跳过
```

语义更清晰，也避免了"配置了行星列表但忘记设 Enabled"的错误。

---

### 改进二：轨道配置随盘面聚合

```
原来：OrbConfigTransit、OrbConfigProgressions、OrbConfigSolarArc 三个平铺字段

现在：
  Charts.Transit.Orbs       ← Transit 专属
  Charts.Progressions.Orbs  ← Progressions 专属
  Charts.SolarArc.Orbs      ← SolarArc 专属
```

新增一种盘面时，只需新增一个 `XxxChartConfig` 结构体，不需要在顶层再加字段。

---

### 改进三：特殊点随盘面聚合

```
原来：SpecialPoints 一个结构体里有四个列表（Natal/Transit/Progressions/SolarArc）

现在：
  NatalChart.Points          ← 出生特殊点
  Charts.Transit.Points      ← 过运特殊点
  Charts.Progressions.Points ← 推运特殊点
  Charts.SolarArc.Points     ← 太阳弧特殊点
```

---

## 六、`buildTasks` 中的读取方式变化

重构后，`buildTasks` 读取配置的方式从：

```
// 原来：到处判断 nil 和 Enabled
if input.ProgressionsConfig != nil && input.ProgressionsConfig.Enabled {
    for _, pid := range input.ProgressionsConfig.Planets { ... }
}
if input.SpecialPoints != nil {
    for _, sp := range input.SpecialPoints.ProgressionsPoints { ... }
}
```

变为：

```
// 重构后：nil 判断统一在工厂函数入口
progressBodies = buildProgressionBodies(input.Charts.Progressions, input.NatalChart.JD)
// buildProgressionBodies 内部：if cfg == nil { return nil }
```

所有"是否启用"的判断收敛到三个工厂函数的入口处。
