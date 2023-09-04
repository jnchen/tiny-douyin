package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

// intToExcelColName 将 int 转换为 Excel 列名。
// 例如：0 -> A, 1 -> B, 25 -> Z, 26 -> AA, 27 -> AB, ...
func intToExcelColName(num int) string {
	if num < 26 {
		return string(rune(num + 'A'))
	} else {
		return intToExcelColName(num/26-1) + intToExcelColName(num%26)
	}
}

func GenerateTestVideos(
	videoFilesOutputPath,
	audioFilesPath string,
	numUsers int,
	poolSize int,
) (numGenerated, numToBeGenerated uint32) {
	if err := os.MkdirAll(videoFilesOutputPath, 0750); err != nil {
		log.Printf("创建目录 %s 失败：%v\n", videoFilesOutputPath, err)
		return 0, 0
	}
	files, err := os.ReadDir(audioFilesPath)
	if err != nil {
		log.Println("读取音频文件夹失败：", err)
		return 0, 0
	}
	nAudio := len(files)

	errLogger := func(err error) {
		log.Println(err)
	}
	generateUserVideos := func(user string, numVideos uint32) uint32 {
		userVideoFilesOutputPath := filepath.Join(videoFilesOutputPath, user)
		if err := os.Mkdir(userVideoFilesOutputPath, 0750); err != nil && !os.IsExist(err) {
			log.Printf("创建用户 %s 目录失败：%v\n", user, err)
			return 0
		}

		log.Printf("开始生成 %s 用户的 %d 个测试视频\n", user, numVideos)

		var generated uint32 = 0
		generator := NewVideoGenerator()
		defer func() {
			_ = generator.Close()
		}()
		for j := uint32(0); j < numVideos; j++ {
			text := fmt.Sprintf("%s%d", user, j+1)
			duration := rand.Intn(nAudio)
			videoFileOutputPath := filepath.Join(userVideoFilesOutputPath, fmt.Sprintf("%d.mp4", j+1))
			audioFilePath := filepath.Join(audioFilesPath, fmt.Sprintf("%d.mp3", duration))
			err := generator.
				SetVideoOutputPath(videoFileOutputPath).
				SetAudioFilePath(audioFilePath).
				SetText(text).
				SetDuration(duration).
				Generate().
				Catch(errLogger)
			if err == nil {
				generated++
			}
		}
		log.Printf("成功生成 %s 用户的 %d/%d 个测试视频\n", user, generated, numVideos)

		return generated
	}

	var wg sync.WaitGroup
	var numGeneratedAtomic atomic.Uint32
	pool := make(chan struct{}, poolSize)

	log.Printf("开始生成 %d 个用户的测试视频\n", numUsers)

	numToBeGenerated = 0
	for i := 0; i < numUsers; i++ {
		user := intToExcelColName(i)
		numVideos := uint32(rand.Intn(16) + 1)
		numToBeGenerated += numVideos

		wg.Add(1)
		pool <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-pool }()
			numGeneratedAtomic.Add(generateUserVideos(user, numVideos))
		}()
	}
	close(pool)
	wg.Wait()

	numGenerated = numGeneratedAtomic.Load()

	log.Printf(
		"成功生成 %d 个用户的 %d/%d 个测试视频\n",
		numUsers,
		numGenerated,
		numToBeGenerated,
	)

	return
}

func main() {
	var (
		numUsers       int
		outputPath     string
		audioFilesPath string
		poolSize       int
	)
	flag.IntVar(&numUsers, "numUsers", 1, "要生成的用户数")
	flag.StringVar(&outputPath, "outputPath", "../../public/test_videos", "输出目录")
	flag.StringVar(&audioFilesPath, "audioFilesPath", "countdown_audio", "音频文件目录")
	flag.IntVar(&poolSize, "poolSize", runtime.NumCPU(), "并发生成视频的协程数")
	flag.Parse()

	GenerateTestVideos(outputPath, audioFilesPath, numUsers, poolSize)
}
