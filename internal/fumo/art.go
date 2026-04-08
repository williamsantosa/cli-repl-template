package fumo

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/bmp"

	"github.com/fumo-cli/fumo-command-line-interface/internal/config"
)

// Frame holds a single pre-rendered frame and its display duration.
type Frame struct {
	Rendered string
	Delay    time.Duration
}

var imageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".bmp":  true,
}

var fumoPixels = []string{
	"......HHHHHH............",
	"....HHHHHHHHH...........",
	"...HHHHHHHHHHHH.........",
	"..HHHHSSSSSSHHH.........",
	"..HHSSSSSSSSSHH.........",
	"..HSSSKSSKSSSSH.........",
	"..HSSSSSSSSSSSH.........",
	"..HSSSSMMSSSSHH.........",
	"..HHSSSSSSSSHHH.........",
	"...RRRRRRRRRR...........",
	"..RRWWWWWWWWWRR.........",
	".RRWWWWWWWWWWWRR........",
	".RWWWWWWWWWWWWWR........",
	".RWWWSSSWWWSSSWR........",
	".RWWWWWWWWWWWWWR........",
	"..RWWWWWWWWWWWR.........",
	"...RRWWWWWWWRR..........",
	"....SSSSSSSSS...........",
}

var palette = map[rune]lipgloss.Color{
	'W': lipgloss.Color("255"),
	'P': lipgloss.Color("205"),
	'R': lipgloss.Color("196"),
	'S': lipgloss.Color("223"),
	'H': lipgloss.Color("94"),
	'B': lipgloss.Color("130"),
	'K': lipgloss.Color("16"),
	'M': lipgloss.Color("210"),
}

func renderBuiltinArt() string {
	var sb strings.Builder
	for i, row := range fumoPixels {
		for _, ch := range row {
			if color, ok := palette[ch]; ok {
				style := lipgloss.NewStyle().Foreground(color).Background(color)
				sb.WriteString(style.Render("██"))
			} else {
				sb.WriteString("  ")
			}
		}
		if i < len(fumoPixels)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return imageExtensions[ext]
}

func isGIF(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".gif"
}

func isOpaque(r, g, b, a uint32) bool {
	return a >= 0x8000
}

func renderHalfBlocks(img image.Image, width int) string {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	cellRatio := config.C.Art.CellRatio
	if cellRatio <= 0 {
		cellRatio = 0.45
	}

	targetW := width
	// Each half-block character is cellRatio wide : 1 tall on screen.
	// Multiply height by cellRatio*2 to compensate (half-blocks pack 2 pixels per row).
	targetH := int(float64(srcH) / float64(srcW) * float64(targetW) * cellRatio * 2)
	if targetH%2 != 0 {
		targetH++
	}

	scaled := resize.Resize(uint(targetW), uint(targetH), img, resize.Lanczos3)

	var sb strings.Builder
	for y := 0; y < targetH; y += 2 {
		for x := 0; x < targetW; x++ {
			topR, topG, topB, topA := scaled.At(x, y).RGBA()
			botR, botG, botB, botA := scaled.At(x, y+1).RGBA()

			topVis := isOpaque(topR, topG, topB, topA)
			botVis := isOpaque(botR, botG, botB, botA)

			switch {
			case topVis && botVis:
				bg := fmt.Sprintf("#%02x%02x%02x", topR>>8, topG>>8, topB>>8)
				fg := fmt.Sprintf("#%02x%02x%02x", botR>>8, botG>>8, botB>>8)
				style := lipgloss.NewStyle().
					Foreground(lipgloss.Color(fg)).
					Background(lipgloss.Color(bg))
				sb.WriteString(style.Render("▄"))
			case topVis:
				fg := fmt.Sprintf("#%02x%02x%02x", topR>>8, topG>>8, topB>>8)
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(fg))
				sb.WriteString(style.Render("▀"))
			case botVis:
				fg := fmt.Sprintf("#%02x%02x%02x", botR>>8, botG>>8, botB>>8)
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(fg))
				sb.WriteString(style.Render("▄"))
			default:
				sb.WriteRune(' ')
			}
		}
		if y+2 < targetH {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

// compositeGIFFrames decodes all frames of an animated GIF,
// properly handling disposal methods, and returns each composited
// full-canvas frame as an image.Image along with its delay.
func compositeGIFFrames(path string) ([]image.Image, []time.Duration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, nil, err
	}

	bounds := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, image.NewUniform(color.Transparent), image.Point{}, draw.Src)

	images := make([]image.Image, len(g.Image))
	delays := make([]time.Duration, len(g.Image))

	for i, frame := range g.Image {
		// Apply disposal from the previous frame before drawing this one
		if i > 0 {
			var disposal byte
			if i-1 < len(g.Disposal) {
				disposal = g.Disposal[i-1]
			}
			switch disposal {
			case gif.DisposalBackground:
				prevBounds := g.Image[i-1].Bounds()
				draw.Draw(canvas, prevBounds, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				prevBounds := g.Image[i-1].Bounds()
				draw.Draw(canvas, prevBounds, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
			}
		}

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

		// Snapshot the canvas
		snapshot := image.NewRGBA(bounds)
		draw.Draw(snapshot, bounds, canvas, image.Point{}, draw.Src)
		images[i] = snapshot

		delay := time.Duration(g.Delay[i]) * 10 * time.Millisecond
		if delay == 0 {
			delay = 100 * time.Millisecond
		}
		delays[i] = delay
	}

	return images, delays, nil
}

func openStaticImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadCustomArt(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\n\r"), nil
}

func applyBorder(art string) string {
	cfg := config.C.Art
	if !cfg.Border {
		return art
	}
	borderColor := lipgloss.Color(cfg.BorderColor)
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)
	return style.Render(art)
}

// RenderArt returns a single static frame (first frame for GIFs).
func RenderArt() string {
	frames := RenderFrames()
	if len(frames) > 0 {
		return frames[0].Rendered
	}
	return applyBorder(renderBuiltinArt())
}

// RenderFrames returns all frames for the configured art source.
// For static images, built-in art, and text files this returns a single frame.
// For animated GIFs, this returns all frames with their proper delays.
func RenderFrames() []Frame {
	cfg := config.C.Art

	if cfg.Source == "" || cfg.Source == "built-in" {
		return []Frame{{Rendered: applyBorder(renderBuiltinArt()), Delay: 0}}
	}

	if isImageFile(cfg.Source) {
		if isGIF(cfg.Source) {
			images, delays, err := compositeGIFFrames(cfg.Source)
			if err != nil {
				return []Frame{{Rendered: applyBorder(renderBuiltinArt()), Delay: 0}}
			}
			frames := make([]Frame, len(images))
			for i, img := range images {
				rendered := applyBorder(renderHalfBlocks(img, cfg.Width))
				frames[i] = Frame{Rendered: rendered, Delay: delays[i]}
			}
			return frames
		}

		img, err := openStaticImage(cfg.Source)
		if err != nil {
			return []Frame{{Rendered: applyBorder(renderBuiltinArt()), Delay: 0}}
		}
		return []Frame{{Rendered: applyBorder(renderHalfBlocks(img, cfg.Width)), Delay: 0}}
	}

	custom, err := loadCustomArt(cfg.Source)
	if err != nil {
		return []Frame{{Rendered: applyBorder(renderBuiltinArt()), Delay: 0}}
	}
	return []Frame{{Rendered: applyBorder(custom), Delay: 0}}
}
