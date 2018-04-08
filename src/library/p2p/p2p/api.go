package p2p

// FileDict 一个文件的元数据信息
type FileDict struct {
	Length int64  `json:"length"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	Sum    string `json:"sum" `
}

// MetaInfo 一个任务内所有文件的元数据信息
type MetaInfo struct {
	Length   int64       `json:"length"`
	PieceLen int64       `json:"PieceLen"`
	Pieces   []byte      `json:"pieces"`
	Files    []*FileDict `json:"files"`
}

// DispatchTask 下发给Agent的分发任务
type DispatchTask struct {
	TaskID    string     `json:"taskId"`
	MetaInfo  *MetaInfo  `json:"metaInfo"`
	LinkChain *LinkChain `json:"linkChain"`
	Speed     int64      `json:"speed"`
}

// StartTask 下发给Agent的分发任务
type StartTask struct {
	TaskID    string     `json:"taskId"`
	LinkChain *LinkChain `json:"linkChain"`
}

// LinkChain 分发路径
type LinkChain struct {
	// 软件分发的路径，要求服务端的地址排在第一个
	DispatchAddrs []string `json:"dispatchAddrs"`
	// 服务端管理接口，用于上报状态
	ServerAddr string `json:"serverAddr"`
}

// PHeader 连接认证消息头
type PHeader struct {
	Len      int32
	TaskID   string
	Username string
	Password string
	Salt     string
}

// StatusReport Agent分发状态上报
type StatusReport struct {
	TaskID          string  `json:"taskId"`
	IP              string  `json:"ip"`
	PercentComplete float32 `json:"percentComplete"`
}
