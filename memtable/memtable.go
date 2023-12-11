package memtable

import (
	"encoding/json"
	"sstable/filesystem"
	"sstable/util"
	"strings"
)

type MemoryTable struct {
	file    filesystem.FileOperation
	records map[string]any
}

func NewMemoryTable(file filesystem.FileOperation) MemoryTable {
	records := make(map[string]any)
	return MemoryTable{file: file, records: records}
}

func (memtable *MemoryTable) Read(key string) any {
	value, ok := memtable.records[key]
	if ok {
		return value
	}

	return nil
}

func (memtable *MemoryTable) Write(key string, value any) error {
	keyValue := util.KeyValueObject{Key: key, Value: value}
	bytes, err := json.Marshal(keyValue)

	if err != nil {
		return err
	}

	line := string(bytes) + "\n"

	if err = memtable.file.AppendBytes([]byte(line)); err != nil {
		return err
	}

	memtable.records[key] = value

	return nil
}

func (memtable *MemoryTable) LoadMemoryTable() error {
	fileOp := memtable.file
	bytes, error := fileOp.ReadAll()
	if error != nil {
		return error
	}

	content := string(bytes)
	memtable.enrichRecordsFromContent(content)

	return nil
}

func (memtable *MemoryTable) enrichRecordsFromContent(content string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		var keyValue *util.KeyValueObject
		err := json.Unmarshal([]byte(line), &keyValue)
		if err == nil {
			memtable.records[keyValue.Key] = keyValue.Value
		}
	}

}
