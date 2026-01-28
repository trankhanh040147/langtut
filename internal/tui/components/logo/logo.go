// Package logo renders a Langtut wordmark in a stylized way.
package logo

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/slice"
	"github.com/trankhanh040147/langtut/internal/tui/styles"
)

// letterform represents a letterform. It can be stretched horizontally by
// a given amount via the boolean argument.
type letterform func(bool) string

const diag = `╱`

// Opts are the options for rendering the Langtut title art.
type Opts struct {
	FieldColor   color.Color // diagonal lines
	TitleColorA  color.Color // left gradient ramp point
	TitleColorB  color.Color // right gradient ramp point
	CharmColor   color.Color // Charm™ text color
	VersionColor color.Color // Version text color
	Width        int         // width of the rendered logo, used for truncation
}

// Render renders the Langtut logo. Set the argument to true to render the narrow
// version, intended for use in a sidebar.
//
// The compact argument determines whether it renders compact for the sidebar
// or wider for the main pane.
func Render(version string, compact bool, o Opts) string {
	const charm = " Charm™"

	fg := func(c color.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	// Title.
	const spacing = 1
	letterforms := []letterform{
		letterP,
		letterR,
		letterE,
		letterP,
		letterF,
	}
	stretchIndex := -1 // -1 means no stretching.
	if !compact {
		stretchIndex = cachedRandN(len(letterforms))
	}

	langtut := renderWord(spacing, stretchIndex, letterforms...)
	langtutWidth := lipgloss.Width(langtut)
	b := new(strings.Builder)
	for r := range strings.SplitSeq(langtut, "\n") {
		fmt.Fprintln(b, styles.ApplyForegroundGrad(r, o.TitleColorA, o.TitleColorB))
	}
	langtut = b.String()

	// Charm and version.
	metaRowGap := 1
	maxVersionWidth := langtutWidth - lipgloss.Width(charm) - metaRowGap
	version = ansi.Truncate(version, maxVersionWidth, "…") // truncate version if too long.
	gap := max(0, langtutWidth-lipgloss.Width(charm)-lipgloss.Width(version))
	metaRow := fg(o.CharmColor, charm) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)

	// Join the meta row and big Langtut title.
	langtut = strings.TrimSpace(metaRow + "\n" + langtut)

	// Narrow version.
	if compact {
		field := fg(o.FieldColor, strings.Repeat(diag, langtutWidth))
		return strings.Join([]string{field, field, langtut, field, ""}, "\n")
	}

	fieldHeight := lipgloss.Height(langtut)

	// Left field.
	const leftWidth = 6
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diag, leftWidth))
	leftField := new(strings.Builder)
	for range fieldHeight {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightWidth := max(15, o.Width-langtutWidth-leftWidth-2) // 2 for the gap.
	const stepDownAt = 0
	rightField := new(strings.Builder)
	for i := range fieldHeight {
		width := rightWidth
		if i >= stepDownAt {
			width = rightWidth - (i - stepDownAt)
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diag, width)), "\n")
	}

	// Return the wide version.
	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, langtut, hGap, rightField.String())
	if o.Width > 0 {
		// Truncate the logo to the specified width.
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		logo = strings.Join(lines, "\n")
	}
	return logo
}

// SmallRender renders a smaller version of the Langtut logo, suitable for
// smaller windows or sidebar usage.
func SmallRender(width int) string {
	t := styles.CurrentTheme()
	title := t.S().Base.Foreground(t.Secondary).Render("Charm™")
	title = fmt.Sprintf("%s %s", title, styles.ApplyBoldForegroundGrad("Langtut", t.Secondary, t.Primary))
	remainingWidth := width - lipgloss.Width(title) - 1 // 1 for the space after "Langtut"
	if remainingWidth > 0 {
		lines := strings.Repeat("╱", remainingWidth)
		title = fmt.Sprintf("%s %s", title, t.S().Base.Foreground(t.Primary).Render(lines))
	}
	return title
}

// renderWord renders letterforms to fork a word. stretchIndex is the index of
// the letter to stretch, or -1 if no letter should be stretched.
func renderWord(spacing int, stretchIndex int, letterforms ...letterform) string {
	if spacing < 0 {
		spacing = 0
	}

	renderedLetterforms := make([]string, len(letterforms))

	// pick one letter randomly to stretch
	for i, letter := range letterforms {
		renderedLetterforms[i] = letter(i == stretchIndex)
	}

	if spacing > 0 {
		// Add spaces between the letters and render.
		renderedLetterforms = slice.Intersperse(renderedLetterforms, strings.Repeat(" ", spacing))
	}
	return strings.TrimSpace(
		lipgloss.JoinHorizontal(lipgloss.Top, renderedLetterforms...),
	)
}

// letterP renders the letter P in a stylized way. It takes a boolean that
// determines whether to stretch the letter. If the stretch is false, it defaults to no stretching.
func letterP(stretch bool) string {
	// Here's what we're making:
	//
	// █▀▀▀▄
	// █   █
	// █▀▀▀
	// █
	// ▀

	left := heredoc.Doc(`
		█
		█
		█
		█
		▀
	`)
	top := heredoc.Doc(`
		▀
	`)
	middle := heredoc.Doc(`
	`)
	bottom := heredoc.Doc(`
		▀
	`)
	right := heredoc.Doc(`
		▄
		█
	`)
	return joinLetterform(
		left,
		stretchLetterformPart(top, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 6,
			maxStretch: 10,
		}),
		right,
		joinLetterform(
			left,
			stretchLetterformPart(middle, letterformProps{
				stretch:    stretch,
				width:      3,
				minStretch: 5,
				maxStretch: 8,
			}),
		),
		joinLetterform(
			left,
			stretchLetterformPart(bottom, letterformProps{
				stretch:    stretch,
				width:      3,
				minStretch: 5,
				maxStretch: 8,
			}),
		),
		left,
	)
}

// letterE renders the letter E in a stylized way. It takes a boolean that
// determines whether to stretch the letter. If the stretch is false, it defaults to no stretching.
func letterE(stretch bool) string {
	// Here's what we're making:
	//
	// █▀▀▀▀
	// █
	// █▀▀▀
	// █
	// ▀▀▀▀

	left := heredoc.Doc(`
		█
		█
		█
		█
		▀
	`)
	top := heredoc.Doc(`
		▀
	`)
	middle := heredoc.Doc(`
		▀
	`)
	bottom := heredoc.Doc(`
		▀
	`)
	return joinLetterform(
		left,
		stretchLetterformPart(top, letterformProps{
			stretch:    stretch,
			width:      4,
			minStretch: 6,
			maxStretch: 10,
		}),
		left,
		stretchLetterformPart(middle, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 5,
			maxStretch: 9,
		}),
		left,
		stretchLetterformPart(bottom, letterformProps{
			stretch:    stretch,
			width:      4,
			minStretch: 6,
			maxStretch: 10,
		}),
	)
}

// letterF renders the letter F in a stylized way. It takes a boolean that
// determines whether to stretch the letter. If the stretch is false, it defaults to no stretching.
func letterF(stretch bool) string {
	// Here's what we're making:
	//
	// █▀▀▀▀
	// █
	// █▀▀▀
	// █
	// ▀

	left := heredoc.Doc(`
		█
		█
		█
		█
		▀
	`)
	top := heredoc.Doc(`
		▀
	`)
	middle := heredoc.Doc(`
		▀
	`)
	return joinLetterform(
		left,
		stretchLetterformPart(top, letterformProps{
			stretch:    stretch,
			width:      4,
			minStretch: 6,
			maxStretch: 10,
		}),
		left,
		stretchLetterformPart(middle, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 5,
			maxStretch: 9,
		}),
		left,
		left,
	)
}

// letterR renders the letter R in a stylized way. It takes an integer that
// determines how many cells to stretch the letter. If the stretch is less than
// 1, it defaults to no stretching.
func letterR(stretch bool) string {
	// Here's what we're making:
	//
	// █▀▀▀▄
	// █▀▀▀▄
	// ▀   ▀

	left := heredoc.Doc(`
		█
		█
		▀
	`)
	center := heredoc.Doc(`
		▀
		▀
	`)
	right := heredoc.Doc(`
		▄
		▄
		▀
	`)
	return joinLetterform(
		left,
		stretchLetterformPart(center, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 7,
			maxStretch: 12,
		}),
		right,
	)
}

// letterSStylized renders the letter S in a stylized way, more so than
// [letterS]. It takes an integer that determines how many cells to stretch the
// letter. If the stretch is less than 1, it defaults to no stretching.
func letterSStylized(stretch bool) string {
	// Here's what we're making:
	//
	// ▄▀▀▀▀▀
	// ▀▀▀▀▀█
	// ▀▀▀▀▀

	left := heredoc.Doc(`
		▄
		▀
		▀
	`)
	center := heredoc.Doc(`
		▀
		▀
		▀
	`)
	right := heredoc.Doc(`
		▀
		█
	`)
	return joinLetterform(
		left,
		stretchLetterformPart(center, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 7,
			maxStretch: 12,
		}),
		right,
	)
}

// letterU renders the letter U in a stylized way. It takes an integer that
// determines how many cells to stretch the letter. If the stretch is less than
// 1, it defaults to no stretching.
func letterU(stretch bool) string {
	// Here's what we're making:
	//
	// █   █
	// █   █
	//	▀▀▀

	side := heredoc.Doc(`
		█
		█
	`)
	middle := heredoc.Doc(`


		▀
	`)
	return joinLetterform(
		side,
		stretchLetterformPart(middle, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 7,
			maxStretch: 12,
		}),
		side,
	)
}

func joinLetterform(letters ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, letters...)
}

// letterformProps defines letterform stretching properties.
// for readability.
type letterformProps struct {
	width      int
	minStretch int
	maxStretch int
	stretch    bool
}

// stretchLetterformPart is a helper function for letter stretching. If randomize
// is false the minimum number will be used.
func stretchLetterformPart(s string, p letterformProps) string {
	if p.maxStretch < p.minStretch {
		p.minStretch, p.maxStretch = p.maxStretch, p.minStretch
	}
	n := p.width
	if p.stretch {
		n = cachedRandN(p.maxStretch-p.minStretch) + p.minStretch //nolint:gosec
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = s
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
