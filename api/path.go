package main

import (
	"log"
	"strings"
)

const PathSeparator = "/"

type Path struct {
	Path       string
	UserName   string
	SeasonName string
}

func NewPath(p string) *Path {
	var user_name, season_name string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	log.Println(s)
	if len(s) > 2 {
		user_name = s[len(s)-2]
		season_name = s[len(s)-1]
		p = strings.Join(s[:len(s)-2], PathSeparator)
	}
	return &Path{Path: p, UserName: user_name, SeasonName: season_name}
}

func (p *Path) HasParameter() bool {
	return len(p.UserName) > 0 && len(p.SeasonName) > 0
}
