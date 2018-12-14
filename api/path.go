package main

import (
	"log"
	"strconv"
	"strings"
)

const PathSeparator = "/"
const ParameterCount = 3

type Path struct {
	Path          string
	UserName      string
	SeasonName    string
	MinStoryCount string
}

func NewPath(p string) *Path {
	var user_name, season_name, lower_bound string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	log.Println(s)
	if len(s) > ParameterCount {
		user_name = s[len(s)-3]
		season_name = s[len(s)-2]
		lower_bound = s[len(s)-1]
		p = strings.Join(s[:len(s)-ParameterCount], PathSeparator)
	}
	return &Path{Path: p, UserName: user_name, SeasonName: season_name, MinStoryCount: lower_bound}
}

func (p *Path) HasParameter() bool {
	if _, err := strconv.Atoi(p.MinStoryCount); err != nil {
		return false
	}
	return len(p.UserName) > 0 && len(p.SeasonName) > 0
}
