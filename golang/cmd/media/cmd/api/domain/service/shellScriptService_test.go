package service_test

import (
	"io"
	"mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/model"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// Decode implements business.IBusiness.
func (m *MockRepository) Decode(reader io.Reader) (interface{}, error) {
	panic("unimplemented")
}

// Encode implements business.IBusiness.
func (m *MockRepository) Encode(interface{}) model.IBuffer {
	panic("unimplemented")
}

// FindByStatus implements repository.IRepository.
func (m *MockRepository) FindByStatus(status int8) ([]model.MediaRecord, error) {
	panic("unimplemented")
}

// Save implements repository.IRepository.
func (m *MockRepository) Save(mediaRecord *model.MediaRecord) (*model.MediaRecord, error) {
	panic("unimplemented")
}

// Terminate implements repository.IRepository.
func (m *MockRepository) Terminate(platformId uint16, streamerId uint32, date time.Time, sequence uint16) (*model.MediaRecord, error) {
	panic("unimplemented")
}

type MockM3u8ParserBiz struct {
	mock.Mock
}

type ShellScriptService struct {
	RootPath string
}

func (s *ShellScriptService) GetRootPath() string {
	return s.RootPath
}

func TestGetFilePath(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockM3u8Parser := new(MockRepository)
	service := service.NewShellScriptService(mockRepo, mockM3u8Parser, "/zzz/mntreamer")
	testDate := time.Date(2025, 04, 06, 0, 0, 0, 0, time.UTC)

	mediaRecord := &model.MediaRecord{
		PlatformId: 1,
		StreamerId: 1001,
		Date:       testDate,
		Sequence:   5,
	}
	tests := []struct {
		name         string
		platformName string
		channelName  string
		expected     string
	}{
		{
			name:         "basic path",
			platformName: "youtube",
			channelName:  "tech-channel",
			expected:     filepath.Join("/zzz/mntreamer", "youtube", "tech-channel", "2025", "04", "06"),
		},
		{
			name:         "special characters",
			platformName: "platform#2",
			channelName:  "channel@test",
			expected:     filepath.Join("/zzz/mntreamer", "platform#2", "channel@test", "2025", "04", "06"),
		},
		{
			name:         "spaces in names",
			platformName: "live platform",
			channelName:  "daily stream",
			expected:     filepath.Join("/zzz/mntreamer", "live platform", "daily stream", "2025", "04", "06"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetFilePath(mediaRecord, tt.platformName, tt.channelName)
			assert.Equal(t, tt.expected, result, "mismatched file path")
		})
	}
}
