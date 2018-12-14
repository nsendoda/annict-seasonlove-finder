package main

import "github.com/naninunenosi/annict-seasonlove-finder/recfilter"

func NarrowRecord(raw []recfilter.Record, n int) []recfilter.Record {
	title_map := map[string]int{}
	narrowed_count := 0
	for i := range raw {
		title_map[raw[i].Work.Title]++
		if title_map[raw[i].Work.Title] == n {
			narrowed_count += n
		} else if title_map[raw[i].Work.Title] > n {
			narrowed_count++
		}
	}
	narrowed := make([]recfilter.Record, narrowed_count)
	for i := range raw {
		if title_map[raw[i].Work.Title] >= n {
			narrowed_count--
			narrowed[narrowed_count] = raw[i]
		}
	}
	return narrowed
}
