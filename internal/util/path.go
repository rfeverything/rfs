package util

import (
	"path"
	"strings"
)

type Path string

func (p *Path) IsDir() bool {
	return strings.HasSuffix(string(*p), "/")
}

func (p *Path) IsFile() bool {
	return !p.IsDir()
}

func (p *Path) IsAbs() bool {
	return strings.HasPrefix(string(*p), "/")
}

func (p *Path) IsRel() bool {
	return !p.IsAbs()
}

func (p *Path) GetFileName() string {
	return path.Base(string(*p))
}

func (p *Path) GetDirectory() string {
	return path.Dir(string(*p))
}

func (p Path) Split() []string {
	if p == "" || p == "/" {
		return []string{}
	}
	return strings.Split(string(p)[1:], "/")
}
