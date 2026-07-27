package fuzzy

import (
	"strings"
)

// Match returns true if pattern characters appear sequentially in target (case-insensitive).
func Match(pattern, target string) bool {
	if pattern == "" {
		return true
	}

	pattern = strings.ToLower(pattern)
	target = strings.ToLower(target)

	pIdx := 0
	pLen := len(pattern)

	for i := 0; i < len(target); i++ {
		if target[i] == pattern[pIdx] {
			pIdx++
			if pIdx == pLen {
				return true
			}
		}
	}
	return false
}

// Score computes match score for ranking search results (higher is better).
func Score(pattern, target string) int {
	if pattern == "" {
		return 0
	}

	pattern = strings.ToLower(pattern)
	target = strings.ToLower(target)

	score := 0
	if strings.Contains(target, pattern) {
		score += 100 - len(target)
	}

	if strings.HasPrefix(target, pattern) {
		score += 50
	}

	return score
}
