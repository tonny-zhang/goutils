package network

import "testing"

func TestDownload(t *testing.T) {
	url := "http://audio04.dmhmusic.com/71_53_T10044825974_128_4_1_0_sdk-cpm/cn/0209/M00/5A/2D/ChR47FsFJAiAAu9nADjHrjjEjno357.mp3?xcode=e8b03d86ad74d8bb6c943b6746663d69032da4b"
	Download(url, ".", ".mp3")
}
