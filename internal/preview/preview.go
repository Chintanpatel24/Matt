package preview

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/Chintanpatel24/Matt/internal/filetree"
)

// PreviewResult contains rendered preview text and metadata.
type PreviewResult struct {
	Content   string
	LineCount int
	IsBinary  bool
	Info      string
}

// GeneratePreview produces syntax-highlighted code, image preview, or formatted text.
func GeneratePreview(entry filetree.FileEntry, maxWidth, maxHeight int) PreviewResult {
	if entry.IsDir {
		return previewDirectory(entry)
	}

	info, err := os.Stat(entry.Path)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  [Error reading file metadata: %v]", err),
			Info:    "Error",
		}
	}

	if info.Size() == 0 {
		return PreviewResult{
			Content: "\n  [Empty file]",
			Info:    "0 Bytes",
		}
	}

	// Check for image files
	ext := strings.ToLower(filepath.Ext(entry.Name))
	if isImageExtension(ext) {
		return previewImage(entry, maxWidth, maxHeight)
	}

	// Read initial chunk to test for binary content
	f, err := os.Open(entry.Path)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  [Error opening file: %v]", err),
			Info:    "Permission Denied / Error",
		}
	}
	defer f.Close()

	headBuf := make([]byte, 2048)
	n, _ := f.Read(headBuf)
	headBuf = headBuf[:n]

	if isBinaryContent(headBuf) {
		return previewBinary(entry, headBuf)
	}

	// Reset file pointer to read text content
	_, _ = f.Seek(0, 0)
	limitReader := io.LimitReader(f, 256*1024)
	contentBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  [Error reading content: %v]", err),
			Info:    "Read Error",
		}
	}

	// Rich preview for Markdown
	if ext == ".md" || ext == ".markdown" {
		formatted := FormatMarkdown(string(contentBytes))
		lines := strings.Split(formatted, "\n")
		return PreviewResult{
			Content:   formatted,
			LineCount: len(lines),
			IsBinary:  false,
			Info:      "Rich Markdown",
		}
	}

	// Rich preview for CSV
	if ext == ".csv" {
		formatted := FormatCSV(string(contentBytes), maxWidth)
		lines := strings.Split(formatted, "\n")
		return PreviewResult{
			Content:   formatted,
			LineCount: len(lines),
			IsBinary:  false,
			Info:      "Rich CSV Table",
		}
	}

	return highlightText(entry, string(contentBytes), maxWidth, maxHeight)
}

func isImageExtension(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp":
		return true
	default:
		return false
	}
}

func previewImage(entry filetree.FileEntry, maxWidth, maxHeight int) PreviewResult {
	f, err := os.Open(entry.Path)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  [Error opening image: %v]", err),
			Info:    "Image Error",
		}
	}
	defer f.Close()

	imgConfig, format, err := image.DecodeConfig(f)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  🖼️ Image File: %s\n  Format: %s (Unrecognized image encoding)\n  Size: %s", entry.Name, filepath.Ext(entry.Name), entry.FormatSize()),
			Info:    "Image",
		}
	}

	// Reset and decode full image for ASCII block rendering
	_, _ = f.Seek(0, 0)
	img, _, err := image.Decode(f)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  🖼️ Image: %s\n  Dimensions: %dx%d px (%s)\n  Size: %s", entry.Name, imgConfig.Width, imgConfig.Height, strings.ToUpper(format), entry.FormatSize()),
			Info:    fmt.Sprintf("%dx%d %s", imgConfig.Width, imgConfig.Height, strings.ToUpper(format)),
		}
	}

	// Generate ASCII pixel preview
	asciiArt := renderASCIIImage(img, max(20, maxWidth-6), max(8, maxHeight-6))

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n  🖼️ Image Preview: %s\n", entry.Name))
	sb.WriteString(fmt.Sprintf("  -----------------------------------\n"))
	sb.WriteString(fmt.Sprintf("  Resolution  : %d x %d px\n", imgConfig.Width, imgConfig.Height))
	sb.WriteString(fmt.Sprintf("  Format      : %s\n", strings.ToUpper(format)))
	sb.WriteString(fmt.Sprintf("  File Size   : %s\n\n", entry.FormatSize()))
	sb.WriteString(asciiArt)

	return PreviewResult{
		Content:   sb.String(),
		LineCount: maxHeight,
		IsBinary:  false,
		Info:      fmt.Sprintf("%dx%d %s", imgConfig.Width, imgConfig.Height, strings.ToUpper(format)),
	}
}

func renderASCIIImage(img image.Image, targetWidth, targetHeight int) string {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	if width <= 0 || height <= 0 {
		return "  [Invalid image dimensions]"
	}

	scaleX := float64(width) / float64(targetWidth)
	scaleY := float64(height) / float64(targetHeight)

	// ASCII ramp for shading
	ramp := " .:-=+*#%@"
	rampLen := len(ramp)

	var sb strings.Builder

	for y := 0; y < targetHeight; y++ {
		sb.WriteString("  ")
		for x := 0; x < targetWidth; x++ {
			srcX := int(float64(x) * scaleX)
			srcY := int(float64(y) * scaleY)

			if srcX >= width {
				srcX = width - 1
			}
			if srcY >= height {
				srcY = height - 1
			}

			r, g, b, _ := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY).RGBA()
			// Convert to 8-bit brightness
			gray := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0 / 256.0
			idx := int((gray / 255.0) * float64(rampLen-1))
			if idx < 0 {
				idx = 0
			}
			if idx >= rampLen {
				idx = rampLen - 1
			}

			sb.WriteByte(ramp[idx])
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func previewDirectory(entry filetree.FileEntry) PreviewResult {
	entries, err := os.ReadDir(entry.Path)
	if err != nil {
		return PreviewResult{
			Content: fmt.Sprintf("\n  [Cannot read directory: %v]", err),
			Info:    "Directory (Unreadable)",
		}
	}

	dirCount, fileCount := 0, 0
	var items []string

	for i, e := range entries {
		if i < 20 {
			icon := "📄"
			if e.IsDir() {
				icon = "📁"
				dirCount++
			} else {
				fileCount++
			}
			items = append(items, fmt.Sprintf("  %s %s", icon, e.Name()))
		} else {
			if e.IsDir() {
				dirCount++
			} else {
				fileCount++
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n  📁 Directory: %s\n", entry.Name))
	sb.WriteString(fmt.Sprintf("  -----------------------------------\n"))
	sb.WriteString(fmt.Sprintf("  Total Items : %d (%d dirs, %d files)\n", len(entries), dirCount, fileCount))
	sb.WriteString(fmt.Sprintf("  Modified    : %s\n", entry.FormatModTime()))
	sb.WriteString(fmt.Sprintf("  Permissions : %s\n\n", entry.Permissions))
	sb.WriteString(fmt.Sprintf("  Contents Overview:\n"))
	sb.WriteString(strings.Join(items, "\n"))

	if len(entries) > 20 {
		sb.WriteString(fmt.Sprintf("\n  ... and %d more items", len(entries)-20))
	}

	return PreviewResult{
		Content:   sb.String(),
		LineCount: len(entries) + 8,
		IsBinary:  false,
		Info:      fmt.Sprintf("%d items", len(entries)),
	}
}

func previewBinary(entry filetree.FileEntry, head []byte) PreviewResult {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n  📦 Binary File: %s\n", entry.Name))
	sb.WriteString(fmt.Sprintf("  -----------------------------------\n"))
	sb.WriteString(fmt.Sprintf("  Size        : %s\n", entry.FormatSize()))
	sb.WriteString(fmt.Sprintf("  Modified    : %s\n", entry.FormatModTime()))
	sb.WriteString(fmt.Sprintf("  Permissions : %s\n\n", entry.Permissions))
	sb.WriteString(fmt.Sprintf("  Hex Peek:\n"))

	for i := 0; i < len(head) && i < 128; i += 16 {
		end := i + 16
		if end > len(head) {
			end = len(head)
		}
		sb.WriteString(fmt.Sprintf("  %04x: % -48x\n", i, head[i:end]))
	}

	return PreviewResult{
		Content:   sb.String(),
		LineCount: 15,
		IsBinary:  true,
		Info:      fmt.Sprintf("Binary (%s)", entry.FormatSize()),
	}
}

func highlightText(entry filetree.FileEntry, text string, maxWidth, maxHeight int) PreviewResult {
	lexer := lexers.Match(entry.Name)
	if lexer == nil {
		lexer = lexers.Analyse(text)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("tty16m")
	if formatter == nil {
		formatter = formatters.TTY16m
	}

	iterator, err := lexer.Tokenise(nil, text)
	if err != nil {
		return formatPlainLines(text, entry)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return formatPlainLines(text, entry)
	}

	lines := strings.Split(buf.String(), "\n")
	var formattedLines []string
	for i, l := range lines {
		formattedLines = append(formattedLines, fmt.Sprintf("\033[90m%4d │\033[0m %s", i+1, l))
	}

	return PreviewResult{
		Content:   strings.Join(formattedLines, "\n"),
		LineCount: len(lines),
		IsBinary:  false,
		Info:      fmt.Sprintf("%s | %d lines", strings.ToUpper(strings.TrimPrefix(filepath.Ext(entry.Name), ".")), len(lines)),
	}
}

func formatPlainLines(text string, entry filetree.FileEntry) PreviewResult {
	lines := strings.Split(text, "\n")
	var formatted []string
	for i, l := range lines {
		formatted = append(formatted, fmt.Sprintf("%4d │ %s", i+1, l))
	}
	return PreviewResult{
		Content:   strings.Join(formatted, "\n"),
		LineCount: len(lines),
		IsBinary:  false,
		Info:      fmt.Sprintf("Plain Text | %d lines", len(lines)),
	}
}

func isBinaryContent(buf []byte) bool {
	if len(buf) == 0 {
		return false
	}
	nullCount := 0
	for _, b := range buf {
		if b == 0 {
			nullCount++
		}
	}
	if nullCount > 0 {
		return true
	}
	return !utf8.Valid(buf)
}
