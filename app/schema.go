package app

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"jhmeeting.com/adminserver/db"
)

// 用户
type User struct {
	Id       int64     `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Password string    `json:"password,omitempty"`
	Ctime    time.Time `json:"ctime,omitempty"`
}

// 房间信息
type RoomInfo struct {
	Id                int64      `json:"id,omitempty"`
	Uid               int64      `json:"uid,omitempty" sql:"index:ri_uid"`            // 房间uid
	RoomName          string     `json:"roomName,omitempty" sql:"index:ri_room_name"` // 房间名称
	ParticipantLimits int        `json:"participantLimits,omitempty"`                 // 房间最高参会人数
	AllowAnonymous    bool       `json:"allowAnonymous,omitempty"`                    // 是否允许匿名用户创建会议
	Config            RoomConfig `json:"roomConfig,omitempty"`                        // 房间配置
	Ctime             time.Time  `json:"ctime,omitempty"  sql:"index:ri_ctime"`       // 创建时间
}

// 房间配置
type RoomConfig struct {
	Resolution            int    `json:"resolution,omitempty"`            // 分辨率，360，480，720，1080
	Subject               string `json:"subject,omitempty"`               // 主题
	LockPassword          string `json:"lockPassword,omitempty"`          // 进入密码
	RequireDisplayName    bool   `json:"requireDisplayName,omitempty"`    // 是否提示参会者输入名字
	StartWithAudioMuted   bool   `json:"startWithAudioMuted,omitempty"`   // 是否加入会议室时不开启音频
	StartWithVideoMuted   bool   `json:"startWithVideoMuted,omitempty"`   // 是否加入会议室时不开启视频
	FileRecordingsEnabled *bool  `json:"fileRecordingsEnabled,omitempty"` // 是否允许服务器录制
	LiveStreamingEnabled  *bool  `json:"liveStreamingEnabled,omitempty"`  // 是否允许直播
	Bandwidth             int    `json:"bandwidth,omitempty"`             // 设置参会者最大上行比特率，默认各分辨率对应的比特率，360：800，480：1000，720：1500，1080：3000
}

func (config RoomConfig) Value() (driver.Value, error) {
	data, _ := json.Marshal(config)

	return string(data), nil
}

func (config *RoomConfig) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for RoomConfig")
	}

	return json.Unmarshal(source, config)
}

// 会议室信息，会议室表示正在开会的房间
type ConferenceInfo struct {
	Id              int64       `json:"id,omitempty"`
	Uid             int64       `json:"uid,omitempty" sql:"index:ci_uid"`            // 会议uid
	RoomName        string      `json:"roomName,omitempty" sql:"index:ci_room_name"` // 房间名称
	Participants    int         `json:"participants,omitempty"`                      // 当前人数
	MaxParticipants int         `json:"maxParticipants,omitempty"`                   // 最高人数
	IsRecording     bool        `json:"isRecording,omitempty"`                       // 是否正在录制
	IsStreaming     bool        `json:"isStreaming,omitempty"`                       // 是否正在直播
	ApiEnabled      bool        `json:"apiEnabled,omitempty"`                        // 是否是使用API接入的会议室
	LockPassword    string      `json:"lockPassword,omitempty"`                      // 进入密码
	Locked          bool        `json:"locked,omitempty"`                            // 是否锁定
	Ctime           time.Time   `json:"ctime,omitempty" sql:"index:ci_ctime"`        // 开始时间
	Etime           db.NullTime `json:"ctime,omitempty" sql:"index:ci_ctime"`        // 结束时间
}

//*****************************************会议回看定义*********************************************************/
// 会议回看信息
type RecordInfo struct {
	Id           int64     `json:"id,omitempty"`
	ConferenceId int64     `json:"conferenceId,omitempty"` // 会议室id
	RoomName     string    `json:"roomName,omitempty"`     // 会议室名称
	Duration     int64     `json:"duration,omitempty"`     // 录制时长
	Size         int64     `json:"size,omitempty"`         // 文件大小
	DownloadUrl  string    `json:"downloadUrl,omitempty"`  // 录像 url 地址
	StreamingUrl string    `json:"streamingUrl,omitempty"` // 推流 url 地址
	Ctime        time.Time `json:"ctime,omitempty"`        // 开始时间
}

// 会议回看表对应的字符串
const (
	RecordTableName   = "record"
	RecordIdCol       = "id"
	RecordRoomNameCol = "room_name"
	RecordCtimeCol    = "ctime"
	RecordDurationCol = "duration"
	RecordSizeCol     = "size"
	RecordUrlCol      = "download_url"
	WhereRecordID     = "id=?"
)
