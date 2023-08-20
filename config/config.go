package config

// Secret 密钥
var Secret = []byte("dancebyte")

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// sftp服务器地址
const FtpAdd = "172.16.102.82:22"
const FtpUser = "zhouyx"
const FtpPsw = "5ocia1-Z"

// PlayUrlPrefix 存储的图片和视频的链接
const PlayUrlPrefix = "http://172.16.102.82/home/zhouyx"
const CoverUrlPrefix = "http://172.16.102.82/home/zhouyximages/"

// HostSSH SSH配置
const HostSSH = "172.16.102.82"
const UserSSH = "zhouyx"
const PasswordSSH = "5ocia1-Z"
const TypeSSH = "password"
const PortSSH = 22
const MaxMsgCount = 100
const SSHHeartbeatTime = 10 * 60

const DateTime = "2006-01-02 15:04:05"

//const ChanCapacity = 10 //chan管道容量，暂时没定
