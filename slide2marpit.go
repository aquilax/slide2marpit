package slide2marpit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	KW_IMAGE         = ".image"
	KW_IFRAME        = ".iframe"
	KW_BACKGROUND    = ".background"
	KW_LINK          = ".link"
	KW_CAPTION       = ".caption"
	KW_COMMENT       = "// "
	KW_SPEAKER_NOTES = ": "
	KW_SLIDE_HEADING = "## "
)

var imageRegex = regexp.MustCompile(`\` + KW_IMAGE + `\s+(\S+)(\s+([\d_]+)%?)?(\s+([\d_]+)%?)?`)
var linkRegex = regexp.MustCompile(`\` + KW_LINK + `\s+(\S+)\s+(.+)`)

func Convert(input io.Reader, output io.Writer) error {
	// Write Marpit front matter
	fmt.Fprintln(output, "---")
	fmt.Fprintln(output, "marpitScanPlugin: true")
	fmt.Fprintln(output, "theme: default")
	fmt.Fprintln(output, "---")

	scanner := bufio.NewScanner(input)
	var inCodeBlock bool

	for scanner.Scan() {
		line := scanner.Text()

		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			fmt.Fprintln(output, line)
			continue
		}

		if inCodeBlock {
			fmt.Fprintln(output, line)
			continue
		}

		// Handle slides
		if strings.HasPrefix(line, KW_SLIDE_HEADING) {
			fmt.Fprintln(output, "---")
			fmt.Fprintln(output)
			fmt.Fprintln(output, line)
			continue
		}

		// Handle presenter notes
		if strings.HasPrefix(line, KW_SPEAKER_NOTES) {
			fmt.Fprintln(output, "<!--", strings.TrimPrefix(line, KW_SPEAKER_NOTES), "-->")
			continue
		}

		// Handle .image syntax
		if strings.HasPrefix(line, KW_IMAGE) {
			fmt.Fprintln(output, handleImage(line))
			continue
		}

		// Handle .background syntax
		if strings.HasPrefix(line, KW_BACKGROUND) {
			fmt.Fprintln(output, handleBackground(line))
			continue
		}

		// Handle .iframe syntax
		if strings.HasPrefix(line, KW_IFRAME) {
			fmt.Fprintln(output, handleIFrame(line))
			continue
		}

		// Handle .link syntax
		if strings.HasPrefix(line, KW_LINK) {
			fmt.Fprintln(output, handleLink(line))
			continue
		}

		// Handle .caption syntax
		if strings.HasPrefix(line, KW_CAPTION) {
			fmt.Fprintln(output, handleCaption(line))
			continue
		}

		// Handle comments syntax
		if strings.HasPrefix(line, KW_COMMENT) {
			continue
		}

		// Handle normal Markdown content
		fmt.Fprintln(output, line)
	}

	if err := scanner.Err(); err != nil {
		return errors.Join(errors.New("error reading input file"), err)
	}
	return nil
}

func handleLink(line string) string {
	if matches := linkRegex.FindStringSubmatch(line); matches != nil {
		url := matches[1]
		label := matches[2]
		return fmt.Sprintf("[%s](%s)\n", label, url)
	}
	return ""
}

func handleCaption(line string) string {
	return strings.TrimPrefix(line, KW_CAPTION)
}

func handleIFrame(line string) string {
	src := strings.TrimPrefix(line, KW_IFRAME)
	return fmt.Sprintf(`<iframe src="%s"></iframe>`, src)
}

func handleBackground(line string) string {
	image := strings.TrimPrefix(line, KW_BACKGROUND)
	return fmt.Sprintf("![bg](%s)", image)
}

func handleImage(line string) string {
	matches := imageRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		imagePath := matches[1]
		width := ""
		height := ""
		if len(matches) > 3 && matches[3] != "" && matches[3] != "_" {
			width = matches[3]
			if !strings.HasSuffix(width, "%") {
				width += "px"
			}
		}
		if len(matches) > 5 && matches[5] != "" && matches[5] != "_" {
			height = matches[5]
			if !strings.HasSuffix(height, "%") {
				height += "px"
			}
		}

		style := ""
		if width != "" {
			style += fmt.Sprintf("width:%s ", width)
		}
		if height != "" {
			style += fmt.Sprintf("height:%s ", height)
		}

		return fmt.Sprintf("![%s](%s)", strings.TrimSpace(style), imagePath)
	}
	return line // If it doesn't match our expected format, print as-is
}
