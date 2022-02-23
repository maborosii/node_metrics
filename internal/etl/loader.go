package etl

import "fmt"

type StoreResult struct {
	ip             string
	cpuAvg         string
	cpuMax         string
	memAvg         string
	memMax         string
	diskUsage      string
	diskReadAvg    string
	diskReadMax    string
	diskWriteAvg   string
	diskWriteMax   string
	networkIn      string
	networkOut     string
	contextSwitchs string
	socketNums     string
}

// 用于灵活构建StoreResult
type Option func(*StoreResult)

func WithCpuAvg(cpuAvg string) Option {
	return func(sr *StoreResult) {
		sr.cpuAvg = cpuAvg
	}
}
func WithCpuMax(cpuMax string) Option {
	return func(sr *StoreResult) {
		sr.cpuMax = cpuMax
	}
}
func WithMemAvg(memAvg string) Option {
	return func(sr *StoreResult) {
		sr.memAvg = memAvg
	}
}
func WithMemMax(memMax string) Option {
	return func(sr *StoreResult) {
		sr.memMax = memMax
	}
}
func WithDiskUsage(diskUsage string) Option {
	return func(sr *StoreResult) {
		sr.diskUsage = diskUsage
	}
}
func WithDiskReadAvg(diskReadAvg string) Option {
	return func(sr *StoreResult) {
		sr.diskReadAvg = diskReadAvg
	}
}
func WithDiskWriteAvg(diskWriteAvg string) Option {
	return func(sr *StoreResult) {
		sr.diskWriteAvg = diskWriteAvg
	}
}
func WithDiskReadMax(diskReadMax string) Option {
	return func(sr *StoreResult) {
		sr.diskReadMax = diskReadMax
	}
}
func WithDiskWriteMax(diskWriteMax string) Option {
	return func(sr *StoreResult) {
		sr.diskWriteMax = diskWriteMax
	}
}
func WithNetworkIn(networkIn string) Option {
	return func(sr *StoreResult) {
		sr.networkIn = networkIn
	}
}
func WithNetworkOut(networkOut string) Option {
	return func(sr *StoreResult) {
		sr.networkOut = networkOut
	}
}
func WithContextSwitchs(contextSwitchs string) Option {
	return func(sr *StoreResult) {
		sr.contextSwitchs = contextSwitchs
	}
}
func WithSocketNums(socketNums string) Option {
	return func(sr *StoreResult) {
		sr.socketNums = socketNums
	}
}

func NewStoreResult(ip string, options ...Option) *StoreResult {
	sr := &StoreResult{ip: ip}
	for _, option := range options {
		option(sr)
	}
	return sr
}

func (sr *StoreResult) GetIp() string {
	return sr.ip
}
func (sr *StoreResult) Print() string {
	return fmt.Sprintf("## ip: %s,cpuAvg: %s,cpuMax: %s,memAvg: %s,memMax: %s,diskUsage: %s,diskReadAvg: %s,diskReadMax: %s,diskWriteAvg: %s,diskWriteMax: %s,networkIn: %s,networkOut: %s,contextSwitchs: %s,socketNums: %s,", sr.ip, sr.cpuAvg, sr.cpuMax, sr.memAvg, sr.memMax, sr.diskUsage, sr.diskReadAvg, sr.diskReadMax, sr.diskWriteAvg, sr.diskWriteMax, sr.networkIn, sr.networkOut, sr.contextSwitchs, sr.socketNums)
}
func (sr *StoreResult) ModifyStoreResult(options ...Option) {
	for _, option := range options {
		option(sr)
	}
}
func (sr *StoreResult) ConvertToSlice() []string {
	// 这里可以使用反射，但为了保证生成的slice的有序性
	return []string{
		sr.ip,
		sr.cpuAvg,
		sr.cpuMax,
		sr.memAvg,
		sr.memMax,
		sr.diskUsage,
		sr.diskReadAvg,
		sr.diskReadMax,
		sr.diskWriteAvg,
		sr.diskWriteMax,
		sr.networkIn,
		sr.networkOut,
		sr.contextSwitchs,
		sr.socketNums,
	}
}

type StoreResults []*StoreResult

func NewStoreResults() StoreResults {
	return []*StoreResult{}
}

func (srs StoreResults) FindIp(ip string) (bool, int) {
	for i, sr := range srs {
		if sr.GetIp() == ip {
			return true, i
		}
	}
	return false, -1
}

func (srs *StoreResults) CreateOrModifyStoreResults(ip string, option Option) {
	ok, index := (*srs).FindIp(ip)
	if ok {
		(*srs)[index].ModifyStoreResult(option)
	} else {
		sr := NewStoreResult(ip, option)
		*srs = append(*srs, sr)
	}
}
