package service

import (
	"context"
	"fmt"
	mntreamerModel "mntreamer/shared/model"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type ShellScriptService struct {
}

func NewShellScriptService() *ShellScriptService {
	return &ShellScriptService{}
}

func (s *ShellScriptService) Download(media *mntreamerModel.Media, streamer *mntreamerModel.Streamer) error {
	now := time.Now()
	channelNameWithNoSpace := strings.ReplaceAll(streamer.ChannelName, " ", "")
	path := s.getFilePath(now, channelNameWithNoSpace)
	s.createFolder(path)
	filename := s.getBaseFilename(now, channelNameWithNoSpace)
	filename = s.getTitle(media.Title, filename)
	filePath := s.getNumberingAndExtension(path, filename)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", media.VideoUrl, "-c", "copy", filePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Running command: %v\n", cmd.Args)
	cmd.Start()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-ctx.Done():
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGINT)
		}
		<-done
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (s *ShellScriptService) getTitle(title string, filename string) string {
	titleWithNoSpace := strings.ReplaceAll(title, " ", "")
	filename += "." + titleWithNoSpace
	return filename
}

func (s *ShellScriptService) getNumberingAndExtension(path string, filename string) string {
	cnt := 1
	filePath := filepath.Join(path, fmt.Sprintf("%s.%d.mp4", filename, cnt))
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		cnt++
		filePath = filepath.Join(path, fmt.Sprintf("%s.%d.mp4", filename, cnt))
	}
	return filePath
}

func (s *ShellScriptService) createFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellScriptService) getFilePath(now time.Time, channelName string) string {
	basePath := "/zzz/mntreamer/chzzk/"
	year := fmt.Sprintf("%d", now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())
	filePath := filepath.Join(basePath, channelName, year, month, day)
	return filePath

}

func (s *ShellScriptService) getBaseFilename(now time.Time, channelName string) string {
	year := now.Year() % 100 // Extract the last two digits of the year
	date := fmt.Sprintf("%02d%02d%02d", year, now.Month(), now.Day())
	return fmt.Sprintf("%s.%s", channelName, date)
}
