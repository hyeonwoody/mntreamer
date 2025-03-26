package service

import (
	"context"
	"fmt"
	"mntreamer/media/cmd/api/domain/business"
	"mntreamer/media/cmd/api/infrastructure/repository"
	model "mntreamer/media/cmd/model"
	mntreamerModel "mntreamer/shared/model"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type ShellScriptService struct {
	bizStrat *business.BusinessStrategy
	repo     repository.IRepository
}

func NewShellScriptService(bizStrat *business.BusinessStrategy, repo repository.IRepository) *ShellScriptService {
	return &ShellScriptService{bizStrat: bizStrat, repo: repo}
}

func (s *ShellScriptService) Download(media *mntreamerModel.Media, channelName string, platformId uint16) error {
	now := time.Now()
	channelNameWithNoSpace := strings.ReplaceAll(channelName, " ", "")
	path := s.getFilePath(now, platformId, channelNameWithNoSpace)
	s.createFolder(path)
	filename := s.getBaseFilename(now, channelNameWithNoSpace)
	filename = s.getTitle(media.Title, filename)
	filePath := s.getNumbering(path, filename)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-i",
		media.VideoUrl,
		"-c", "copy",
		"-f", "segment",
		"-segment_time", "60",
		"-segment_format", "mpegts",
		"-segment_list", filePath+".m3u8",
		"-segment_list_type", "m3u8",
		"-strftime", "1",
		filePath+".%H%M%S.ts")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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

func (s *ShellScriptService) getNumbering(path string, filename string) string {
	cnt := 1
	filePath := filepath.Join(path, fmt.Sprintf("%s.%d.m3u8", filename, cnt))
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		cnt++
		filePath = filepath.Join(path, fmt.Sprintf("%s.%d.m3u8", filename, cnt))
	}
	return filepath.Join(path, fmt.Sprintf("%s.%d", filename, cnt))
}

func (s *ShellScriptService) createFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellScriptService) GetFilePath(mediaRecord *model.MediaRecord, channelName string) string {
	basePath := s.bizStrat.GetDownloadPath(mediaRecord.PlatformId)
	year := mediaRecord.Date.Format("2006")
	month := mediaRecord.Date.Format("01")
	day := mediaRecord.Date.Format("02")
	filePath := filepath.Join(basePath, channelName, year, month, day)
	return filePath
}

func (s *ShellScriptService) getFilePath(now time.Time, platformId uint16, channelName string) string {
	basePath := s.bizStrat.GetDownloadPath(platformId)
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

func (s *ShellScriptService) Save(platformId uint16, streamerId uint32) {
	s.repo.Save(model.NewMediaRecord(platformId, streamerId))
}

func (s *ShellScriptService) GetFiles(filePath string) ([]model.FileInfo, error) {
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	var fileInfos []model.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(filePath, file.Name())
		fileInfos = append(fileInfos, model.FileInfo{
			Name:        file.Name(),
			IsDirectory: file.IsDir(),
			Path:        fullPath,
			Size:        info.Size(),
			UpdatedAt:   info.ModTime().UTC().Format(http.TimeFormat),
		})
	}
	return fileInfos, nil
}

func (s *ShellScriptService) GetM3u8(filePath string) ([]model.FileInfo, error) {
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	var fileInfos []model.FileInfo
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".m3u8") {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(filePath, file.Name())
		fileInfos = append(fileInfos, model.FileInfo{
			Name:        file.Name(),
			IsDirectory: file.IsDir(),
			Path:        fullPath,
			Size:        info.Size(),
			UpdatedAt:   info.ModTime().UTC().Format(http.TimeFormat),
		})
	}
	return fileInfos, nil
}

func (s *ShellScriptService) GetMediaToRefine() ([]model.MediaRecord, error) {
	return s.repo.FindByStatus(mntreamerModel.IDLE)
}
