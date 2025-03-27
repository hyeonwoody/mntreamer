package service

import (
	"bufio"
	"context"
	"fmt"
	parserBusiness "mntreamer/media/cmd/api/domain/business/parser"
	platform "mntreamer/media/cmd/api/domain/business/platform"
	"mntreamer/media/cmd/api/infrastructure/repository"
	"mntreamer/media/cmd/model"
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
	bizStrat      *platform.BusinessStrategy
	m3u8ParserBiz parserBusiness.IBusiness
	repo          repository.IRepository
	rootPath      string
}

func NewShellScriptService(bizStrat *platform.BusinessStrategy, repo repository.IRepository, m3u8ParserBiz parserBusiness.IBusiness, basePath string) *ShellScriptService {
	return &ShellScriptService{bizStrat: bizStrat, m3u8ParserBiz: m3u8ParserBiz, repo: repo, rootPath: basePath}
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
		if !strings.HasSuffix(file.Name(), ".m3u8") {
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

func (s *ShellScriptService) Stream(fullPath string) (string, error) {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", err
	}
	return fullPath, nil
}

func (s *ShellScriptService) StreamMediaPlaylist(fullPath string) (string, error) {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", err
	}
	return fullPath, nil
}

func (s *ShellScriptService) StreamSegment(fullPath string) (string, error) {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", err
	}
	return fullPath, nil
}

func (s *ShellScriptService) Decode(path string) (interface{}, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("🛑failed to open file %s: %w", path, err)
	}
	defer file.Close()

	playList, err := s.m3u8ParserBiz.Decode(bufio.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("🛑failed to decode file %s: %w", path, err)
	}
	return playList, nil
}

func (s *ShellScriptService) Excise(path string, begin float64, end float64) error {
	var mpl *model.MediaPlaylist
	playlist, _ := s.Decode(path)
	mpl, ok := playlist.(*model.MediaPlaylist)
	if !ok {
		return fmt.Errorf("🛑decoded playlist is unknown")
	}
	var duration float64
	segmentsToRemove := []uint{}
	for i := uint(0); duration < end; i++ {
		segment, err := mpl.GetSegment(i)
		if err != nil {
			return fmt.Errorf("🛑failed to get segment %d: %w", i, err)
		}
		if begin <= duration && duration < end {
			segmentsToRemove = append(segmentsToRemove, i)
		}
		duration += segment.Duration
	}
	if len(segmentsToRemove) == 0 {
		return fmt.Errorf("🛑nothing to remove")
	}

	filePath := filepath.Dir(path)
	removeIdx := segmentsToRemove[0]
	for range segmentsToRemove {
		segment, _ := mpl.GetSegment(removeIdx)
		segmentPath := filepath.Join(filePath, segment.Uri)
		if err := os.Remove(segmentPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("🛑failed to delete segment %s: %w", segment.Uri, err)
		}
		mpl.PullSegment(removeIdx)
	}
	mpl.SetDiscontinuityWithIndex(removeIdx, true)

	buf := s.m3u8ParserBiz.Encode(mpl)
	return s.WriteBufferToFile(path, buf)
}

func (s *ShellScriptService) WriteBufferToFile(path string, buf model.IBuffer) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("🛑failed to open file %s for writing: %w", path, err)
	}
	defer file.Close()
	if _, err := file.Write(buf.GetData()); err != nil {
		return fmt.Errorf("🛑failed to write updated playlist: %w", err)
	}
	return nil
}
