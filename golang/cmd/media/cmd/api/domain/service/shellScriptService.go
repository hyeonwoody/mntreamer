package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	parserBusiness "mntreamer/media/cmd/api/domain/business/parser"
	"mntreamer/media/cmd/api/infrastructure/repository"
	"mntreamer/media/cmd/model"
	mntreamerModel "mntreamer/shared/model"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ShellScriptService struct {
	m3u8ParserBiz parserBusiness.IBusiness
	repo          repository.IRepository
	rootPath      string
}

func NewShellScriptService(repo repository.IRepository, m3u8ParserBiz parserBusiness.IBusiness, basePath string) *ShellScriptService {
	return &ShellScriptService{m3u8ParserBiz: m3u8ParserBiz, repo: repo, rootPath: basePath}
}

func (s *ShellScriptService) GetRootPath() string {
	return s.rootPath
}

func (s *ShellScriptService) Download(media *mntreamerModel.Media, channelName string, platformId uint16) error {
	now := time.Now()
	channelNameWithNoSpace := strings.ReplaceAll(channelName, " ", "")
	path := s.getFilePath(now, media.OutputUrl, channelName)
	s.createFolder(path)
	filename := s.getBaseFilename(now, channelNameWithNoSpace)
	filename = s.getTitle(media.Title, filename)

	sequence := s.getSequence(path, filename)
	filePath := filepath.Join(path, fmt.Sprintf("%s.%d", filename, sequence))

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
		filepath.Join(path, fmt.Sprintf("%d", sequence))+".%H%M%S.ts")
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
	reg := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	titleCleaned := reg.ReplaceAllString(title, "")
	titleWithUnderscores := strings.ReplaceAll(titleCleaned, " ", "_")
	titleTrimmed := strings.Trim(titleWithUnderscores, "_.")
	if titleTrimmed == "" {
		titleTrimmed = "untitled"
	}
	return filename + "." + titleTrimmed
}

func (s *ShellScriptService) getSequence(path string, filename string) uint {
	cnt := 0
	filePath := filepath.Join(path, fmt.Sprintf("%s.%d.m3u8", filename, cnt))
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		cnt++
		filePath = filepath.Join(path, fmt.Sprintf("%s.%d.m3u8", filename, cnt))
	}
	return uint(cnt)
}

func (s *ShellScriptService) createFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShellScriptService) GetFilePath(mediaRecord *model.MediaRecord, platformName, channelName string) string {
	basePath := fmt.Sprintf("%s/%s", s.GetRootPath(), platformName)
	year := mediaRecord.Date.Format("2006")
	month := mediaRecord.Date.Format("01")
	day := mediaRecord.Date.Format("02")
	filePath := filepath.Join(basePath, channelName, year, month, day)
	return filePath
}

func (s *ShellScriptService) getRelativePath(filePath string) string {
	relativePath, err := filepath.Rel(s.rootPath, filePath)
	if err != nil {
		return ""
	}
	return relativePath
}

func (s *ShellScriptService) GetPlatformNameByFilePath(filePath string) (string, error) {
	relPath := s.getRelativePath(filePath)
	if len(relPath) == 0 {
		return "", fmt.Errorf("invalid file path")
	}
	return strings.Split(relPath, "/")[0], nil
}

func (s *ShellScriptService) GetChannelNameByFilePath(filePath string) (string, error) {
	relPath := s.getRelativePath(filePath)
	if len(relPath) == 0 {
		return "", fmt.Errorf("invalid file path")
	}
	return strings.Split(relPath, "/")[1], nil
}

func (s *ShellScriptService) GetDateByFilePath(fullPath string) (time.Time, error) {
	date := strings.Split(filepath.Base(fullPath), ".")[1]
	return time.Parse("060102", date)
}

func (s *ShellScriptService) GetSequenceByFilePath(fullPath string) (uint16, error) {
	filenameComponent := strings.Split(filepath.Base(fullPath), ".")
	sequenceStr := filenameComponent[len(filenameComponent)-2]
	sequence, err := strconv.ParseUint(sequenceStr, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sequence number: %w", err)
	}
	return uint16(sequence), nil
}

func (s *ShellScriptService) getFilePath(now time.Time, url string, channelName string) string {
	basePath := fmt.Sprintf("%s/%s", s.GetRootPath(), url)
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

func (s *ShellScriptService) Save(platformId uint16, streamerId uint32, status int8) {
	s.repo.Save(model.NewMediaRecord(platformId, streamerId, status))
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

func (s *ShellScriptService) GetMediaFiles(fullPath string) ([]model.FileInfo, error) {
	filePath := filepath.Dir(fullPath)
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	filename := filepath.Base(fullPath)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	var fileInfos []model.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), filename) {
			fullPath := filepath.Join(filePath, file.Name())
			fileInfos = append(fileInfos, model.FileInfo{
				Name:        file.Name(),
				IsDirectory: file.IsDir(),
				Path:        fullPath,
				Size:        info.Size(),
				UpdatedAt:   info.ModTime().UTC().Format(http.TimeFormat),
			})
		}
	}
	return fileInfos, nil
}

func (s *ShellScriptService) GetM3u8(filePath string, sequence uint16) ([]model.FileInfo, error) {
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	//TODO
	suffix := fmt.Sprintf("%d.m3u8", sequence)
	var fileInfos []model.FileInfo
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), suffix) {
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
		break
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

func (s *ShellScriptService) StreamMp4(fullPath string) (*os.File, error) {
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return file, nil
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
	if strings.HasSuffix(path, ".mp4") {
		return s.exciseMp4(path, begin, end)
	}
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

func (s *ShellScriptService) exciseMp4(fullPath string, begin float64, end float64) error {
	filePath := filepath.Dir(fullPath)
	filename := filepath.Base(fullPath)
	tmpPath := filepath.Join(filePath, "FFMPEG_"+filename)
	os.Rename(fullPath, tmpPath)
	parts := []string{fmt.Sprintf("%s.part1.mp4", fullPath), fmt.Sprintf("%s.part2.mp4", fullPath)}
	var wg sync.WaitGroup
	var err error
	wg.Add(2)
	go func() {
		defer wg.Done()
		if trimErr := s.Trim(tmpPath, 0, begin, parts[0]); trimErr != nil {
			err = trimErr
		}
	}()
	go func() {
		defer wg.Done()
		if trimErr := s.Trim(tmpPath, end, -1, parts[1]); trimErr != nil {
			err = trimErr
		}
	}()
	wg.Wait()
	if err != nil {
		os.Rename(tmpPath, fullPath)
		return err
	}

	err = s.Combine(parts, fullPath)
	for _, file := range parts {
		os.Remove(file)
	}

	if err != nil {
		os.Rename(tmpPath, fullPath)
		return err
	}

	os.Remove(tmpPath)
	return nil
}

func (s *ShellScriptService) Combine(files []string, output string) error {
	filePath := filepath.Dir(output)
	fileList := filePath + "/list.txt"
	f, err := os.Create(fileList)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, file := range files {
		_, err := f.WriteString(fmt.Sprintf("file '%s'\n", file))
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(
		"ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", fileList,
		"-c", "copy", output)
	fmt.Printf("Running command: %v\n", cmd.Args)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = cmd.Run()
	os.Remove(fileList)
	return err
}

func (s *ShellScriptService) Trim(fullPath string, begin, end float64, output string) error {
	beginStr := fmt.Sprintf("%.3f", begin)
	durationStr := fmt.Sprintf("%.3f", end-begin)
	args := []string{
		"-y",
		"-i", fullPath,
		"-c", "copy",
		"-ss", beginStr}
	if end != -1 {
		args = append(args, "-t", durationStr)
	}
	args = append(args, output)
	cmd := exec.Command("ffmpeg", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	fmt.Printf("Running command: %v\n", cmd.Args)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("FFmpeg error: %v\nError details: %s\n", err, stderr.String())
	}
	return err
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
	file.Close()
	return nil
}

func (s *ShellScriptService) Delete(platformId uint16, streamerId uint32, fullPath string) (*model.MediaRecord, error) {
	if strings.HasSuffix(fullPath, ".mp4") {
		s.deleteMp4(fullPath)
		return nil, nil
	}
	s.deleteMedia(fullPath)
	s.deleteMediaPlaylist(fullPath)
	date, sequence, err := s.getDateAndSequeceFromFullPath(fullPath)
	if err != nil {
		return nil, err
	}
	terminated, err := s.UpdateStatus(model.NewInstance(platformId, streamerId, date, sequence, mntreamerModel.TERMINATED), mntreamerModel.TERMINATED)
	if err != nil {
		return nil, err
	}
	return terminated, nil
}

func (s *ShellScriptService) Confirm(platformId uint16, streamerId uint32, fullPath string) (*model.MediaRecord, error) {
	date, sequence, err := s.getDateAndSequeceFromFullPath(fullPath)
	if err != nil {
		return nil, err
	}
	confirmed, err := s.UpdateStatus(model.NewInstance(platformId, streamerId, date, sequence, mntreamerModel.DONE), mntreamerModel.DONE)
	if err != nil {
		return nil, err
	}
	return confirmed, nil
}

func (s *ShellScriptService) getDateAndSequeceFromFullPath(fullPath string) (time.Time, uint16, error) {
	date, err := s.GetDateByFilePath(fullPath)
	if err != nil {
		return time.Time{}, 0, err
	}
	sequence, err := s.GetSequenceByFilePath(fullPath)
	if err != nil {
		return time.Time{}, 0, err
	}
	return date, sequence, nil
}

func (s *ShellScriptService) UpdateStatus(mediaRecord *model.MediaRecord, status int8) (*model.MediaRecord, error) {
	mediaRecord.Status = status
	updated, err := s.repo.Save(mediaRecord)
	return updated, err
}

func (s *ShellScriptService) deleteMp4(fullPath string) error {
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("🛑failed to delete media playlist %s: %w", fullPath, err)
	}
	return nil
}

func (s *ShellScriptService) deleteMediaPlaylist(fullPath string) error {
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("🛑failed to delete media playlist %s: %w", fullPath, err)
	}
	return nil
}

func (s *ShellScriptService) deleteMedia(fullPath string) error {
	playlist, _ := s.Decode(fullPath)
	mpl, ok := playlist.(*model.MediaPlaylist)
	if !ok {
		return fmt.Errorf("🛑decoded playlist is unknown")
	}
	filePath := filepath.Dir(fullPath)
	for i := range mpl.Count() {
		seg, _ := mpl.GetSegment(i)
		if err := os.Remove(filepath.Join(filePath, seg.Uri)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("🛑failed to delete media %s: %w", filepath.Join(filePath, seg.Uri), err)
		}
	}
	return nil
}
