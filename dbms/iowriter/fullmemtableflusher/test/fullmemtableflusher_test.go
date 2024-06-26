package fullmemtableflusher_test

import (
	"sstable/test/util/mockmemtable"
	testdbms "sstable/test/util/testdbms"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFulledMemtable(t *testing.T) {
	//arrange

	memtable := mockmemtable.NewAlmostFullMemtable()
	memtable.Write("name", "deea")
	dbms := testdbms.NewDummyDbms(nil)
	testdbms.InitializeDbmsPartially(dbms)
	testdbms.AddFullMemtable(dbms, memtable)

	//act
	databaseIO := dbms.DatabaseIO
	fullMemtableFlusher := dbms.FullMemtableFlusher

	fullMemtableFlusher.FlushFulledSStable()
	state := dbms.StateManagement.GetState()
	readResult, _ := databaseIO.Read("name")

	//assert
	assert.Equal(t, 0, len(state.FulledMemoryTables))
	assert.Equal(t, 1, len(state.Metadata.MemtableToSSTable))
	assert.Equal(t, "deea", readResult)
}
