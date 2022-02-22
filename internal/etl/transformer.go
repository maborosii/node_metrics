package etl

import (
	. "node_metrics_go/pkg/log"
)

func ShuffleResult(series int, storeResults *StoreResults) {
	defer WgReceiver.Done()
	for i := 0; i < series; i++ {
		queryResult := <-metricsChan
		results := queryResult.CleanValue(Pattern)
		for _, result := range results {
			switch queryResult.GetLabel() {
			case "cpu_usage_avg_percents":
				storeResults.CreateOrModifyStoreResults(result[0], WithCpuAvg(result[1]+"%"))
			case "cpu_usage_max_percents":
				storeResults.CreateOrModifyStoreResults(result[0], WithCpuMax(result[1]+"%"))
			case "mem_usage_avg_percents":
				storeResults.CreateOrModifyStoreResults(result[0], WithMemAvg(result[1]+"%"))
			case "mem_usage_max_percents":
				storeResults.CreateOrModifyStoreResults(result[0], WithMemMax(result[1]+"%"))
			case "rootdir_disk_usage_percents":
				storeResults.CreateOrModifyStoreResults(result[0], WithDiskUsage(result[1]+"%"))
			case "disk_read_speed_avg_KB_persec":
				storeResults.CreateOrModifyStoreResults(result[0], WithDiskReadAvg(result[1]+"KB/s"))
			case "disk_read_speed_max_KB_persec":
				storeResults.CreateOrModifyStoreResults(result[0], WithDiskReadMax(result[1]+"KB/s"))
			case "disk_write_speed_avg_KB_persec":
				storeResults.CreateOrModifyStoreResults(result[0], WithDiskWriteAvg(result[1]+"KB/s"))
			case "disk_write_speed_max_KB_persec":
				storeResults.CreateOrModifyStoreResults(result[0], WithDiskWriteMax(result[1]+"KB/s"))
			case "network_in_speed_MB_persec":
				storeResults.CreateOrModifyStoreResults(result[0], WithNetworkIn(result[1]+"MB/s"))
			case "network_out_speed_MB_persec":
				storeResults.CreateOrModifyStoreResults(result[0], WithNetworkOut(result[1]+"MB/s"))
			case "context_switches_persec_K":
				storeResults.CreateOrModifyStoreResults(result[0], WithContextSwitchs(result[1]+"K"))
			case "socket_nums_K":
				storeResults.CreateOrModifyStoreResults(result[0], WithSocketNums(result[1]+"K"))
			default:
				Log.Info("Default")
			}
		}
	}
	close(notifyChan)
}
