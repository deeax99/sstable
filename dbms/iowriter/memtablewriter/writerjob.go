package memtablewriter

import (
	"sstable/dbms/components/memtablemanagement"
	"sstable/dbms/statemanagement"
	"sstable/util"
)

const receiverChanSize = 1 << 16
const receiverWorkers = 12

type MemtableWriterJob struct {
	stateManagement    *statemanagement.DatabaseManagementStateManagement
	memtableManagement *memtablemanagement.MemtableManagement
	receiverChan       chan receiverChanData
	writerChan         chan writerChanData
}

func (memtableWriter *MemtableWriterJob) Initialize() {
	for i := 0; i < receiverWorkers; i++ {
		go util.RunInLoop(memtableWriter.ReceiverExec)
	}
	go util.RunInLoop(memtableWriter.WriterExec)
}

func NewMemtableWriteJob(
	stateManagement *statemanagement.DatabaseManagementStateManagement,
	memtableManagement *memtablemanagement.MemtableManagement) *MemtableWriterJob {
	return &MemtableWriterJob{
		stateManagement:    stateManagement,
		memtableManagement: memtableManagement,
		receiverChan:       make(chan receiverChanData, receiverChanSize),
		writerChan:         make(chan writerChanData, 1),
	}
}
