package main

import (
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"gocv.io/x/gocv"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func renameFile(oldPath, newPath string) (err error) {
	if err = os.Rename(oldPath, newPath); err != nil {
		log.Printf("重命名 %s 为 %s 失败：", oldPath, newPath)
		log.Println(err)
		return
	}
	return
}

func loadFont(fontFile string) (face font.Face, ctx *freetype.Context, err error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return
	}
	ttfOptions := &truetype.Options{
		Size: 64,
		DPI:  72,
	}

	ctx = freetype.NewContext()
	ctx.SetFont(f)
	ctx.SetFontSize(ttfOptions.Size)
	ctx.SetDPI(ttfOptions.DPI)

	face = truetype.NewFace(f, ttfOptions)

	return
}

func setMatRGBFromImage(mat gocv.Mat, img image.Image) error {
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()

	// Ensure that the provided Mat has the correct size and type
	if mat.Cols() != dx ||
		mat.Rows() != dy ||
		mat.Type() != gocv.MatTypeCV8UC3 {
		return errors.New("Mat size or type mismatch")
	}

	data, err := mat.DataPtrUint8()
	if err != nil {
		return fmt.Errorf("获取图像矩阵数据指针失败：%w", err)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[y*dx*3+x*3+0] = byte(b >> 8)
			data[y*dx*3+x*3+1] = byte(g >> 8)
			data[y*dx*3+x*3+2] = byte(r >> 8)
		}
	}
	return nil
}

func getTextSize(fontFace font.Face, text string) (width, height int) {
	width = font.MeasureString(fontFace, text).Ceil()
	height = fontFace.Metrics().Height.Ceil()
	return
}

type VideoGenerator struct {
	videoOutputPath     string
	tempVideoOutputPath string
	audioFilePath       string
	text                string
	duration            int

	// optionalArgs
	fps       int
	size      image.Point
	bgColor   color.RGBA
	textColor color.RGBA
	fontFace  font.Face
	fontCtx   *freetype.Context

	err              error
	bgColorUniform   image.Uniform
	textColorUniform image.Uniform
	img              *image.RGBA
	imgMat           gocv.Mat
}

func NewVideoGenerator(
	args ...interface{},
) *VideoGenerator {
	if len(args)%2 != 0 {
		log.Panicln("省略参数个数必须为偶数！")
	}

	var ok bool
	var err error
	g := &VideoGenerator{
		videoOutputPath: "",
		audioFilePath:   "",
		text:            "",
		duration:        0,
		fps:             30,
		size:            image.Pt(480, 640),
		bgColor:         color.RGBA{R: 0, G: 43, B: 54, A: 255},
		textColor:       color.RGBA{R: 253, G: 246, B: 227, A: 255},
		fontFace:        nil,
		fontCtx:         nil,
		err:             nil,
	}
	fontFile := "LXGWWenKaiMono.ttf"
	for i := 0; i < len(args); i += 2 {
		switch args[i] {
		case "videoOutputPath":
			g.videoOutputPath, ok = args[i+1].(string)
			if !ok {
				log.Panicln("videoOutputPath 必须是 string 类型")
			}
		case "audioFilePath":
			g.audioFilePath, ok = args[i+1].(string)
			if !ok {
				log.Panicln("audioFilePath 必须是 string 类型")
			}
			_, err = os.Stat(g.audioFilePath)
			if os.IsNotExist(err) {
				log.Panicln("音频文件不存在！")
			} else if err != nil {
				log.Panicln(err)
			}
		case "text":
			g.text, ok = args[i+1].(string)
			if !ok {
				log.Panicln("text 必须是 string 类型")
			}
		case "duration":
			g.duration, ok = args[i+1].(int)
			if !ok {
				log.Panicln("duration 必须是 int 类型")
			}
		case "fps":
			g.fps, ok = args[i+1].(int)
			if !ok {
				log.Panicln("fps 必须是 int 类型")
			}
		case "size":
			g.size, ok = args[i+1].(image.Point)
			if !ok {
				log.Panicln("size 必须是 image.Point 类型")
			}
		case "bgColor":
			g.bgColor, ok = args[i+1].(color.RGBA)
			if !ok {
				log.Panicln("bgColor 必须是 color.RGBA 类型")
			}
		case "textColor":
			g.textColor, ok = args[i+1].(color.RGBA)
			if !ok {
				log.Panicln("textColor 必须是 color.RGBA 类型")
			}
		case "fontFile":
			fontFile, ok = args[i+1].(string)
			if !ok {
				log.Panicln("fontFile 必须是 string 类型")
			}
		}
	}

	g.bgColorUniform = image.Uniform{C: g.bgColor}
	g.textColorUniform = image.Uniform{C: g.textColor}
	g.img = image.NewRGBA(image.Rect(0, 0, g.size.X, g.size.Y))
	if g.imgMat, err = gocv.ImageToMatRGB(g.img); err != nil {
		log.Panicf("创建 imgMat 失败：%v\n", err)
	}
	_ = g.SetFont(fontFile).
		Catch(func(err error) {
			log.Panicln(err)
		})
	return g
}

func (g *VideoGenerator) Close() (err error) {
	if g.fontFace != nil || g.fontCtx != nil {
		err = g.fontFace.Close() // 对于 truetype，总返回 nil
		g.fontFace = nil
		g.fontCtx = nil
	}
	g.img = nil
	runtime.GC()
	err = g.imgMat.Close()
	return
}

func (g *VideoGenerator) SetVideoOutputPath(videoOutputPath string) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.videoOutputPath = videoOutputPath
	return g
}

func (g *VideoGenerator) SetAudioFilePath(audioFilePath string) *VideoGenerator {
	if g.err != nil {
		return g
	}

	if _, g.err = os.Stat(audioFilePath); g.err != nil {
		if os.IsNotExist(g.err) {
			g.err = fmt.Errorf("音频文件不存在！%w", g.err)
		}
		return g
	}
	g.audioFilePath = audioFilePath
	return g
}

func (g *VideoGenerator) SetText(text string) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.text = text
	return g
}

func (g *VideoGenerator) SetDuration(duration int) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.duration = duration
	return g
}

func (g *VideoGenerator) SetFPS(fps int) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.fps = fps
	return g
}

func (g *VideoGenerator) SetSize(size image.Point) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.size = size

	g.img = image.NewRGBA(image.Rect(0, 0, g.size.X, g.size.Y))
	if g.fontFace != nil && g.fontCtx != nil {
		g.fontCtx.SetClip(g.img.Bounds())
		g.fontCtx.SetDst(g.img)
	}
	_ = g.imgMat.Close()
	g.imgMat, g.err = gocv.ImageToMatRGB(g.img)
	return g
}

func (g *VideoGenerator) SetBGColor(bgColor color.RGBA) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.bgColor = bgColor
	g.bgColorUniform = image.Uniform{C: g.bgColor}
	return g
}

func (g *VideoGenerator) SetTextColor(textColor color.RGBA) *VideoGenerator {
	if g.err != nil {
		return g
	}

	g.textColor = textColor
	g.textColorUniform = image.Uniform{C: g.textColor}
	if g.fontCtx != nil {
		g.fontCtx.SetSrc(&g.textColorUniform)
	}
	return g
}

func (g *VideoGenerator) SetFont(fontFile string) *VideoGenerator {
	if g.err != nil {
		return g
	}

	if g.fontFace, g.fontCtx, g.err = loadFont(fontFile); g.err != nil {
		g.err = fmt.Errorf("字体文件加载失败：%w", g.err)
		return g
	}
	if g.img != nil {
		g.fontCtx.SetClip(g.img.Bounds())
		g.fontCtx.SetDst(g.img)
		g.fontCtx.SetSrc(&image.Uniform{C: g.textColor})
	}
	return g
}

func (g *VideoGenerator) Catch(handlers ...func(err error)) error {
	err := g.err
	if err != nil && len(handlers) > 0 {
		for _, handler := range handlers {
			handler(err)
		}
	}
	g.err = nil
	return err
}

func (g *VideoGenerator) drawBackground() *VideoGenerator {
	if g.err != nil {
		return g
	}

	draw.Draw(g.img, g.img.Bounds(), &g.bgColorUniform, image.Point{}, draw.Src)
	return g
}

func (g *VideoGenerator) drawText(text string, pt fixed.Point26_6) *VideoGenerator {
	if g.err != nil {
		return g
	}

	_, g.err = g.fontCtx.DrawString(text, pt)
	if g.err != nil {
		g.err = fmt.Errorf("绘制文字 %s 失败：%w", text, g.err)
		return g
	}
	return g
}

var mu sync.Mutex

func (g *VideoGenerator) Print() *VideoGenerator {
	fields := reflect.TypeOf(*g)
	values := reflect.ValueOf(*g)

	maxFieldWidth := 0
	for i := 0; i < fields.NumField(); i++ {
		fieldName := fields.Field(i).Name
		fieldValue := fmt.Sprintf("%v", values.Field(i))
		if len(fieldName) > maxFieldWidth {
			maxFieldWidth = len(fieldName)
		}
		if len(fieldValue) > maxFieldWidth {
			maxFieldWidth = len(fieldValue)
		}
	}
	format := "| %-" + fmt.Sprintf("%d", maxFieldWidth+2) + "s | %-" + fmt.Sprintf("%d", maxFieldWidth+2) + "s |\n"

	mu.Lock()
	defer mu.Unlock()

	fmt.Println()
	fmt.Println("+", strings.Repeat("-", maxFieldWidth+2), "+", strings.Repeat("-", maxFieldWidth+2), "+")
	fmt.Printf(format, "Field", "Value")
	fmt.Println("+", strings.Repeat("-", maxFieldWidth+2), "+", strings.Repeat("-", maxFieldWidth+2), "+")
	for i := 0; i < fields.NumField(); i++ {
		fieldName := fields.Field(i).Name
		fieldValue := fmt.Sprintf("%v", values.Field(i))
		fmt.Printf(format, fieldName, fieldValue)
	}
	fmt.Println("+", strings.Repeat("-", maxFieldWidth+2), "+", strings.Repeat("-", maxFieldWidth+2), "+")
	fmt.Println()

	return g
}

func (g *VideoGenerator) generateVideo() (err error) {
	g.tempVideoOutputPath = strings.Replace(g.videoOutputPath, ".mp4", ".temp.mp4", 1)
	var videoWriter *gocv.VideoWriter
	if videoWriter, err = gocv.VideoWriterFile(
		g.tempVideoOutputPath,
		"mp4v",
		float64(g.fps),
		g.size.X, g.size.Y,
		true,
	); err != nil {
		err = fmt.Errorf("VideoWriter创建失败：%w", err)
		return
	}
	defer func() {
		if err = videoWriter.Close(); err != nil {
			err = fmt.Errorf("VideoWriter关闭失败：%w", err)
		}
		videoWriter = nil
	}()

	countdown := g.duration
	textWidth, textHeight := getTextSize(g.fontFace, g.text)
	countdownWidth, countdownHeight := getTextSize(g.fontFace, strconv.Itoa(countdown))
	textPt := freetype.Pt(
		(g.size.X-textWidth)/2,
		(g.size.Y-textHeight)/4,
	)
	countdownPt := freetype.Pt(
		(g.size.X-countdownWidth)/2,
		(g.size.Y-countdownHeight)*3/4,
	)

	frames := (g.duration + 1) * g.fps
	for i := 0; i < frames; i++ {
		if i%g.fps == 0 {
			g.drawBackground().
				drawText(g.text, textPt).
				drawText(strconv.Itoa(countdown), countdownPt)

			if err = setMatRGBFromImage(g.imgMat, g.img); err != nil {
				err = fmt.Errorf("图片转换失败：%w", err)
				return
			}
		} else if i%g.fps == g.fps-1 { // 每秒减一
			countdown--
		}

		if err = videoWriter.Write(g.imgMat); err != nil {
			err = fmt.Errorf("视频写入失败：%w", err)
			return
		}
	}

	return
}

func (g *VideoGenerator) combineAudioAndVideo() (err error) {
	command := []string{
		"ffmpeg",
		"-y",
		"-i", g.tempVideoOutputPath,
		"-i", g.audioFilePath,
		"-c", "copy",
		"-strict", "experimental",
		g.videoOutputPath,
	}

	cmd := exec.Command(command[0], command[1:]...)
	var output []byte
	output, err = cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf(
			"合并音频和视频失败！%w\n程序返回值：%d\n程序输出：%s\n",
			err,
			cmd.ProcessState.ExitCode(),
			string(output),
		)
		return
	}
	return
}

func (g *VideoGenerator) Generate() *VideoGenerator {
	if g.err != nil {
		return g
	}
	defer func() {
		// 不论什么情况，都尝试删除临时文件
		if err := os.Remove(g.tempVideoOutputPath); err != nil && !os.IsNotExist(err) {
			log.Println(err)
		}
	}()

	if g.err = g.generateVideo(); g.err != nil {
		return g
	}

	if g.audioFilePath == "" {
		// 没有音频文件，直接将临时文件重命名为最终文件
		_ = renameFile(g.tempVideoOutputPath, g.videoOutputPath)
	} else if g.err = g.combineAudioAndVideo(); g.err != nil {
		// 或合并音频和视频出错，直接将临时文件重命名为最终文件
		_ = renameFile(g.tempVideoOutputPath, g.videoOutputPath)
	}

	return g
}
