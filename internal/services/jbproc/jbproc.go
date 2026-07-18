// Package jbproc detects, controls, and relaunches running JetBrains IDEs.
package jbproc

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/insigmo/jetreset/internal/config"
	"github.com/insigmo/jetreset/internal/logx"
)

// ErrTimeout is returned when a process did not exit within the wait window.
var ErrTimeout = errors.New("jbproc: process did not exit in time")

// ErrNoLauncher is returned when the IDE launcher could not be located.
var ErrNoLauncher = errors.New("jbproc: launcher not found")

// Proc describes a running JetBrains IDE.
type Proc struct {
	Product string
	PID     int
}

const ideaMain = "com.intellij.idea.Main"

var (
	reSelector = regexp.MustCompile(`-Didea\.paths\.selector=([A-Za-z]+)`)
	rePrefix   = regexp.MustCompile(`-Didea\.platform\.prefix=([A-Za-z]+)`)
	reInstall  = regexp.MustCompile(`[/\\]([A-Za-z]+)-\d+\.\d+\.\d+[/\\]`)
)

// minCSVFields is the minimum number of fields in a tasklist CSV row.
const minCSVFields = 2

// Canonicalize maps a selector/prefix token (e.g. "GoLand", "IdeaC") to its
// canonical product name; it returns "" for unknown tokens.
func Canonicalize(token string) string {
	if product, ok := config.TokenToProduct[token]; ok {
		return product
	}
	before, isCommunity := strings.CutSuffix(token, "C")
	if !isCommunity {
		return ""
	}
	if product, ok := config.TokenToProduct[before]; ok {
		return product
	}
	return ""
}

// ProductFromTokens identifies the product from a joined command line.
func ProductFromTokens(joined string) string {
	for _, re := range []*regexp.Regexp{reSelector, rePrefix, reInstall} {
		if m := re.FindStringSubmatch(joined); m != nil {
			if p := Canonicalize(m[1]); p != "" {
				return p
			}
		}
	}
	return ""
}

// DedupByProduct collapses a list to one entry per product, sorted by PID.
func DedupByProduct(ps []Proc) []Proc {
	seen := make(map[string]bool)
	var out []Proc
	for _, p := range ps {
		if p.Product == "" || seen[p.Product] {
			continue
		}
		seen[p.Product] = true
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PID < out[j].PID })
	return out
}

// ProcsFromPS parses "ps -o pid=,command=" output into running JetBrains IDEs.
func ProcsFromPS(out string) []Proc {
	var ps []Proc
	for line := range strings.SplitSeq(out, "\n") {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		before, after, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		pid, err := strconv.Atoi(before)
		if err != nil {
			continue
		}
		if !strings.Contains(after, ideaMain) {
			continue
		}
		product := ProductFromTokens(after)
		if product == "" {
			continue
		}
		logx.Debugf("detected %s pid=%d", product, pid)
		ps = append(ps, Proc{product, pid})
	}
	return DedupByProduct(ps)
}

// ProcsFromTasklist parses "tasklist /FO CSV /NH" output into running IDEs.
func ProcsFromTasklist(out string) []Proc {
	var ps []Proc
	for line := range strings.SplitSeq(out, "\n") {
		f := strings.Split(strings.TrimSpace(line), "\",\"")
		if len(f) < minCSVFields {
			continue
		}
		image := strings.Trim(f[0], "\"")
		pid, err := strconv.Atoi(strings.Trim(f[1], "\""))
		if err != nil {
			continue
		}
		for product, exe := range config.ProductExe {
			if strings.EqualFold(image, exe) {
				logx.Debugf("detected %s pid=%d", product, pid)
				ps = append(ps, Proc{product, pid})
			}
		}
	}
	return DedupByProduct(ps)
}
