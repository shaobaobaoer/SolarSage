package transit

import (
	"sort"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

// CalcTransitEvents computes all transit events in the given time range.
// Refactored version: delegates to task-based architecture.
func CalcTransitEvents(input TransitCalcInput) ([]models.TransitEvent, error) {
	// Step 1: Build calculation context (pre-compute fixed data)
	ctx, err := buildCalcContext(input)
	if err != nil {
		return nil, err
	}

	// Step 2: Build task list (declarative enumeration)
	tasks := buildTasks(ctx)

	// Step 3: Execute all tasks
	allEvents := runAll(tasks, ctx)

	// Step 4: Add Void of Course Moon (post-processing)
	if ctx.Input.EventFilter.VoidOfCourse {
		vocEvents := findVoidOfCourse(allEvents, ctx.StartJD, ctx.EndJD)
		allEvents = append(allEvents, vocEvents...)
	}

	// Step 5: Sort events by JD
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].JD < allEvents[j].JD
	})

	return allEvents, nil
}
