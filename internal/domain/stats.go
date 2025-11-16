package domain

type Stats struct {
    ReviewAssignments map[string]int
    PRStatuses        map[PRStatus]int
}
