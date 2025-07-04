/*
Copyright (C) 2016 Felipe da Cunha Gonçalves
All Rights Reserved.

MIT LICENSE

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// copied from github.com/xrash/smetrics
// fix transpositions calculation

package util

import "math"

func JaroWinkler(a, b string, boostThreshold float64, prefixSize int) float64 {
	j := jaro(a, b)

	if j <= boostThreshold {
		return j
	}

	prefixSize = min(len(a), len(b), prefixSize)

	var prefixMatch float64
	for i := 0; i < prefixSize; i++ {
		if a[i] == b[i] {
			prefixMatch++
		} else {
			break
		}
	}

	return j + 0.1*prefixMatch*(1.0-j)
}

// The Jaro distance. The result is 1 for equal strings, and 0 for completely different strings.
func jaro(a, b string) float64 {
	// If both strings are zero-length, they are completely equal,
	// therefore return 1.
	if len(a) == 0 && len(b) == 0 {
		return 1
	}

	// If one string is zero-length, strings are completely different,
	// therefore return 0.
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	// Define the necessary variables for the algorithm.
	la := float64(len(a))
	lb := float64(len(b))
	matchRange := int(math.Max(0, math.Floor(math.Max(la, lb)/2.0)-1))
	matchesA := make([]bool, len(a))
	matchesB := make([]bool, len(b))
	var matches float64 = 0

	// Step 1: Matches
	// Loop through each character of the first string,
	// looking for a matching character in the second string.
	for i := 0; i < len(a); i++ {
		start := int(math.Max(0, float64(i-matchRange)))
		end := int(math.Min(lb-1, float64(i+matchRange)))

		for j := start; j <= end; j++ {
			if matchesB[j] {
				continue
			}

			if a[i] == b[j] {
				matchesA[i] = true
				matchesB[j] = true
				matches++
				break
			}
		}
	}

	// If there are no matches, strings are completely different,
	// therefore return 0.
	if matches == 0 {
		return 0
	}

	// Step 2: Transpositions
	// Loop through the matches' arrays, looking for
	// unaligned matches. Count the number of unaligned matches.
	unaligned := 0
	j := 0
	for i := 0; i < len(a); i++ {
		if !matchesA[i] {
			continue
		}

		for !matchesB[j] {
			j++
		}

		if a[i] != b[j] {
			unaligned++
		}

		j++
	}

	// The number of unaligned matches divided by two, is the number of _transpositions_.
	transpositions := math.Floor(float64(unaligned) / 2)

	// Jaro distance is the average between these three numbers:
	// 1. matches / length of string A
	// 2. matches / length of string B
	// 3. (matches - transpositions/matches)
	// So, all that divided by three is the final result.
	return ((matches / la) + (matches / lb) + ((matches - transpositions) / matches)) / 3.0
}
