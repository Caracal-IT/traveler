package colours

import (
	"fmt"
	"strings"
)

type Color int

const (
	Red Color = 1 << iota
	Green
	Blue

	// Yellow composites and convenience
	Yellow  = Red | Green
	Cyan    = Green | Blue
	Magenta = Red | Blue
	White   = Red | Green | Blue
)

// Has reports whether a flag is present in c.
func (c Color) Has(flag Color) bool { return c&flag != 0 }

// Add sets flag(s) on c.
func (c *Color) Add(flag Color) { *c |= flag }

// Remove clears flag(s) from c.
func (c *Color) Remove(flag Color) { *c &^= flag }

// String returns a human-readable representation like "Red|Blue".
func (c Color) String() string {
	if c == 0 {
		return "None"
	}
	parts := make([]string, 0, 3)
	if c&Red != 0 {
		parts = append(parts, "Red")
	}
	if c&Green != 0 {
		parts = append(parts, "Green")
	}
	if c&Blue != 0 {
		parts = append(parts, "Blue")
	}
	if len(parts) == 0 {
		return fmt.Sprintf("Color(%d)", uint8(c))
	}
	return strings.Join(parts, "|")
}
