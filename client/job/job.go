package job

var csJob = NewConfigSyncJob()

func ConfigSync() *ConfigSyncJob {
	return csJob
}
