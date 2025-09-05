package app

import "flag"

// Config holds the application configuration
type Config struct {
	RecordFile string
	ReplayFile string
}

// ParseFlags parses command line flags and returns configuration
func ParseFlags() *Config {
	recordFile := flag.String("record", "", "Record conversation to file")
	replayFile := flag.String("replay", "", "Replay conversation from file")
	flag.Parse()

	return &Config{
		RecordFile: *recordFile,
		ReplayFile: *replayFile,
	}
}

// IsRecording returns true if recording is enabled
func (c *Config) IsRecording() bool {
	return c.RecordFile != ""
}

// IsReplaying returns true if replay mode is enabled
func (c *Config) IsReplaying() bool {
	return c.ReplayFile != ""
}

// IsNormalMode returns true if neither recording nor replaying
func (c *Config) IsNormalMode() bool {
	return !c.IsRecording() && !c.IsReplaying()
}