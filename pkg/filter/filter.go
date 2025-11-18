package filter

import (
	"logGen/model"
	"sort"
	"time"
)

func FilterEntries(
	Segments []model.Segment,
	levels, components, hosts, reqIDs []string, startTime time.Time, endTime time.Time,
) []model.LogEntry {
	var result []model.LogEntry

	for _, segment := range Segments {

		if !startTime.IsZero() && segment.EndTime.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && segment.StartTime.After(endTime) {
			continue
		}

		matchedIndex := make(map[int]bool)
		totalFilters := 0

		// ---- LEVELS ----
		if len(levels) > 0 {
			totalFilters++
			for _, level := range levels {
				for _, idx := range segment.Index.ByLevel[level] {

					// BOUNDS CHECK
					if idx < 0 || idx >= len(segment.LogEntries) {
						continue
					}

					matchedIndex[idx] = true
				}
			}
		}

		// ---- COMPONENTS ----
		if len(components) > 0 {
			totalFilters++
			componentFilter := make(map[int]bool)
			for _, component := range components {
				for _, idx := range segment.Index.ByComponent[component] {

					if idx < 0 || idx >= len(segment.LogEntries) {
						continue
					}

					if len(matchedIndex) == 0 || matchedIndex[idx] {
						componentFilter[idx] = true
					}
				}
			}
			matchedIndex = componentFilter
		}

		// ---- HOSTS ----
		if len(hosts) > 0 {
			totalFilters++
			hostFilter := make(map[int]bool)
			for _, host := range hosts {
				for _, idx := range segment.Index.ByHost[host] {

					if idx < 0 || idx >= len(segment.LogEntries) {
						continue
					}

					if len(matchedIndex) == 0 || matchedIndex[idx] {
						hostFilter[idx] = true
					}
				}
			}
			matchedIndex = hostFilter
		}

		// ---- REQUEST ID ----
		if len(reqIDs) > 0 {
			totalFilters++
			reqFilter := make(map[int]bool)
			for _, req := range reqIDs {
				for _, idx := range segment.Index.ByReqId[req] {

					if idx < 0 || idx >= len(segment.LogEntries) {
						continue
					}

					if len(matchedIndex) == 0 || matchedIndex[idx] {
						reqFilter[idx] = true
					}
				}
			}
			matchedIndex = reqFilter
		}

		// ---- NO FILTERS (return everything in time range) ----
		if totalFilters == 0 {
			for _, entry := range segment.LogEntries {

				if !startTime.IsZero() && entry.Time.Before(startTime) {
					continue
				}
				if !endTime.IsZero() && entry.Time.After(endTime) {
					continue
				}

				result = append(result, entry)
			}
			continue
		}

		// ---- SORT INDEXES FOR CONSISTENT ORDER ----
		var idxs []int
		for idx := range matchedIndex {
			idxs = append(idxs, idx)
		}
		sort.Ints(idxs)

		for _, idx := range idxs {
			entry := segment.LogEntries[idx]

			if !startTime.IsZero() && entry.Time.Before(startTime) {
				continue
			}
			if !endTime.IsZero() && entry.Time.After(endTime) {
				continue
			}

			result = append(result, entry)
		}
	}

	return result
}
