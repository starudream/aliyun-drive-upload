package aliyunDriveUpload

import (
	"io"
)

type progressEventType int

const (
	transferStartedEvent progressEventType = 1 + iota
	transferDataEvent
	transferCompletedEvent
	transferFailedEvent
)

type ProgressEvent struct {
	TotalBytes    int64
	ConsumedBytes int64
	RwBytes       int64
	EventType     progressEventType
}

type ProgressHook func(event *ProgressEvent)

type ProgressReader struct {
	reader io.Reader

	hook ProgressHook

	totalBytes    int64
	consumedBytes int64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	if err != nil {
		if err == io.EOF {
			pr.hook(&ProgressEvent{TotalBytes: pr.totalBytes, ConsumedBytes: pr.consumedBytes, EventType: transferCompletedEvent})
		} else {
			pr.hook(&ProgressEvent{TotalBytes: pr.totalBytes, ConsumedBytes: pr.consumedBytes, EventType: transferFailedEvent})
		}
	}
	if n > 0 {
		pr.consumedBytes += int64(n)
		pr.hook(&ProgressEvent{TotalBytes: pr.totalBytes, ConsumedBytes: pr.consumedBytes, RwBytes: int64(n), EventType: transferDataEvent})
	}
	return n, err
}

func (pr *ProgressReader) Close() error {
	if rc, ok := pr.reader.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}
