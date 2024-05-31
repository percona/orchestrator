package inst

import (
	"testing"

	test "github.com/openark/golib/tests"
)

func TestNewOracleGtidSetEntry(t *testing.T) {
	{
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)
		test.S(t).ExpectEquals(entry.UUID, "00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(entry.DefaultIv, "1-7")
	}
	{
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7:10-20"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)
		test.S(t).ExpectEquals(entry.UUID, "00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(entry.DefaultIv, "1-7:10-20")
	}
	{
		uuidSet := "00020194-3333-3333-3333-333333333333"
		_, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNotNil(err)
	}

	// Negative tests for tagged gtids
	{
		// Empty tag name
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7:10-20:"
		_, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNotNil(err)
	}
	{
		// Invalid tag name
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7:10-20:123tag1:1"
		_, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNotNil(err)
	}
	{
		// Valid tag name but no interval after tag name
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7:10-20:tag1"
		_, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNotNil(err)
	}
	{
		// Valid tag name but empty interval after tag name
		uuidSet := "00020194-3333-3333-3333-333333333333:tag1:"
		_, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNotNil(err)
	}
	{
		// Valid tag name but invalid interval after tag name
		uuidSet := "00020194-3333-3333-3333-333333333333:tag1:random"
		_, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNotNil(err)
	}

	// Positive tests for tagged gtids
	{
		// Single valued tagged interval
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7:10-20:tag1:1"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)
		test.S(t).ExpectEquals(entry.UUID, "00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(entry.DefaultIv, "1-7:10-20")

		// There is only one tagged interval
		test.S(t).ExpectEquals(len(entry.TaggedIv), 1)

		// There is only tagged interval range
		test.S(t).ExpectEquals(len(entry.TaggedIv[0].Interval), 1)

		test.S(t).ExpectEquals(entry.TaggedIv[0].Tag, "tag1")
		test.S(t).ExpectEquals(entry.TaggedIv[0].Interval[0], "1")
	}
	{
		// Multi valued tagged interval
		uuidSet := "00020194-3333-3333-3333-333333333333:1-7:10-20:tag1:1-56"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)
		test.S(t).ExpectEquals(entry.UUID, "00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(entry.DefaultIv, "1-7:10-20")

		// There is only one tagged interval
		test.S(t).ExpectEquals(len(entry.TaggedIv), 1)

		// There are only tagged interval range
		test.S(t).ExpectEquals(len(entry.TaggedIv[0].Interval), 1)

		test.S(t).ExpectEquals(entry.TaggedIv[0].Tag, "tag1")
		test.S(t).ExpectEquals(entry.TaggedIv[0].Interval[0], "1-56")
	}
	{
		// No default interval
		uuidSet := "00020194-3333-3333-3333-333333333333:domain.com:1:51-56"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)
		test.S(t).ExpectEquals(entry.UUID, "00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(entry.DefaultIv, "")

		// There is only one tagged interval
		test.S(t).ExpectEquals(len(entry.TaggedIv), 1)

		// There are two tagged interval ranges
		test.S(t).ExpectEquals(len(entry.TaggedIv[0].Interval), 2)

		test.S(t).ExpectEquals(entry.TaggedIv[0].Tag, "domain.com")
		test.S(t).ExpectEquals(entry.TaggedIv[0].Interval[0], "1")
		test.S(t).ExpectEquals(entry.TaggedIv[0].Interval[1], "51-56")
	}
	{
		// No default interval, multiple tagged intervals
		uuidSet := "00020194-3333-3333-3333-333333333333:domain1.com:1:51-56:domain2.com:1-78:151-256:514"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)
		test.S(t).ExpectEquals(entry.UUID, "00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(entry.DefaultIv, "")

		// There are two tagged intervals
		test.S(t).ExpectEquals(len(entry.TaggedIv), 2)

		// tag 1 (domain1.com) has two interval ranges
		test.S(t).ExpectEquals(len(entry.TaggedIv[0].Interval), 2)

		test.S(t).ExpectEquals(entry.TaggedIv[0].Tag, "domain1.com")
		test.S(t).ExpectEquals(entry.TaggedIv[0].Interval[0], "1")
		test.S(t).ExpectEquals(entry.TaggedIv[0].Interval[1], "51-56")

		// tag 2 (domain2.com) has three interval ranges
		test.S(t).ExpectEquals(len(entry.TaggedIv[1].Interval), 3)

		test.S(t).ExpectEquals(entry.TaggedIv[1].Tag, "domain2.com")
		test.S(t).ExpectEquals(entry.TaggedIv[1].Interval[0], "1-78")
		test.S(t).ExpectEquals(entry.TaggedIv[1].Interval[1], "151-256")
		test.S(t).ExpectEquals(entry.TaggedIv[1].Interval[2], "514")
	}

}

func TestExplode(t *testing.T) {
	{
		uuidSet := "00020194-3333-3333-3333-333333333333:7"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)

		exploded := entry.Explode()
		test.S(t).ExpectEquals(len(exploded), 1)
		test.S(t).ExpectEquals(exploded[0].String(), "00020194-3333-3333-3333-333333333333:7")
	}
	{
		uuidSet := "00020194-3333-3333-3333-333333333333:1-3"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)

		exploded := entry.Explode()
		test.S(t).ExpectEquals(len(exploded), 3)
		test.S(t).ExpectEquals(exploded[0].String(), "00020194-3333-3333-3333-333333333333:1")
		test.S(t).ExpectEquals(exploded[1].String(), "00020194-3333-3333-3333-333333333333:2")
		test.S(t).ExpectEquals(exploded[2].String(), "00020194-3333-3333-3333-333333333333:3")
	}
	{
		uuidSet := "00020194-3333-3333-3333-333333333333:1-3:6-7"
		entry, err := NewOracleGtidSetEntry(uuidSet)
		test.S(t).ExpectNil(err)

		exploded := entry.Explode()
		test.S(t).ExpectEquals(len(exploded), 5)
		test.S(t).ExpectEquals(exploded[0].String(), "00020194-3333-3333-3333-333333333333:1")
		test.S(t).ExpectEquals(exploded[1].String(), "00020194-3333-3333-3333-333333333333:2")
		test.S(t).ExpectEquals(exploded[2].String(), "00020194-3333-3333-3333-333333333333:3")
		test.S(t).ExpectEquals(exploded[3].String(), "00020194-3333-3333-3333-333333333333:6")
		test.S(t).ExpectEquals(exploded[4].String(), "00020194-3333-3333-3333-333333333333:7")
	}
	{
		gtidSetVal := "00020192-1111-1111-1111-111111111111:29-30, 00020194-3333-3333-3333-333333333333:7-8"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		exploded := gtidSet.Explode()
		test.S(t).ExpectEquals(len(exploded), 4)
		test.S(t).ExpectEquals(exploded[0].String(), "00020192-1111-1111-1111-111111111111:29")
		test.S(t).ExpectEquals(exploded[1].String(), "00020192-1111-1111-1111-111111111111:30")
		test.S(t).ExpectEquals(exploded[2].String(), "00020194-3333-3333-3333-333333333333:7")
		test.S(t).ExpectEquals(exploded[3].String(), "00020194-3333-3333-3333-333333333333:8")
	}

	// Explode tests for tagged GTIDS
	{
		// Single valued tagged gtid with default interval
		gtidSetVal := "00020192-1111-1111-1111-111111111111:29-30:tag1:1"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		exploded := gtidSet.Explode()
		test.S(t).ExpectEquals(len(exploded), 3)
		test.S(t).ExpectEquals(exploded[0].String(), "00020192-1111-1111-1111-111111111111:29")
		test.S(t).ExpectEquals(exploded[1].String(), "00020192-1111-1111-1111-111111111111:30")
		test.S(t).ExpectEquals(exploded[2].String(), "00020192-1111-1111-1111-111111111111:tag1:1")

	}
	{
		// Multi valued tagged gtid with default interval
		gtidSetVal := "00020192-1111-1111-1111-111111111111:29-30:tag1:1-4"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		exploded := gtidSet.Explode()
		test.S(t).ExpectEquals(len(exploded), 6)
		test.S(t).ExpectEquals(exploded[0].String(), "00020192-1111-1111-1111-111111111111:29")
		test.S(t).ExpectEquals(exploded[1].String(), "00020192-1111-1111-1111-111111111111:30")
		test.S(t).ExpectEquals(exploded[2].String(), "00020192-1111-1111-1111-111111111111:tag1:1")
		test.S(t).ExpectEquals(exploded[3].String(), "00020192-1111-1111-1111-111111111111:tag1:2")
		test.S(t).ExpectEquals(exploded[4].String(), "00020192-1111-1111-1111-111111111111:tag1:3")
		test.S(t).ExpectEquals(exploded[5].String(), "00020192-1111-1111-1111-111111111111:tag1:4")
	}
	{
		// Gtid sets with multiple uuids
		gtidSetVal :=
			"00020192-1111-1111-1111-111111111111:29-30:tag1:123," +
				"00020194-3333-3333-3333-333333333333:414-416:tag2:7-8," +
				"00020199-2222-2222-2222-111111111111:tag3:1-2:5"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		exploded := gtidSet.Explode()
		test.S(t).ExpectEquals(len(exploded), 11)
		test.S(t).ExpectEquals(exploded[0].String(), "00020192-1111-1111-1111-111111111111:29")
		test.S(t).ExpectEquals(exploded[1].String(), "00020192-1111-1111-1111-111111111111:30")
		test.S(t).ExpectEquals(exploded[2].String(), "00020192-1111-1111-1111-111111111111:tag1:123")
		test.S(t).ExpectEquals(exploded[3].String(), "00020194-3333-3333-3333-333333333333:414")
		test.S(t).ExpectEquals(exploded[4].String(), "00020194-3333-3333-3333-333333333333:415")
		test.S(t).ExpectEquals(exploded[5].String(), "00020194-3333-3333-3333-333333333333:416")
		test.S(t).ExpectEquals(exploded[6].String(), "00020194-3333-3333-3333-333333333333:tag2:7")
		test.S(t).ExpectEquals(exploded[7].String(), "00020194-3333-3333-3333-333333333333:tag2:8")
		test.S(t).ExpectEquals(exploded[8].String(), "00020199-2222-2222-2222-111111111111:tag3:1")
		test.S(t).ExpectEquals(exploded[9].String(), "00020199-2222-2222-2222-111111111111:tag3:2")
		test.S(t).ExpectEquals(exploded[10].String(), "00020199-2222-2222-2222-111111111111:tag3:5")
	}

}

func TestNewOracleGtidSet(t *testing.T) {
	{
		gtidSetVal := "00020191-1111-1111-1111-111111111111:20-30," +
			"00020192-3333-3333-3333-333333333333:7-8:tag1:1," +
			"00020193-3333-3333-3333-333333333333:tag2:110-789:tag3:11447"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 3)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020191-1111-1111-1111-111111111111:20-30")
		test.S(t).ExpectEquals(gtidSet.GtidEntries[1].String(), "00020192-3333-3333-3333-333333333333:7-8:tag1:1")
		test.S(t).ExpectEquals(gtidSet.GtidEntries[2].String(), "00020193-3333-3333-3333-333333333333:tag2:110-789:tag3:11447")
	}
	{
		gtidSetVal := "00020192-1111-1111-1111-111111111111:20-30, 00020194-3333-3333-3333-333333333333:7-8"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020192-1111-1111-1111-111111111111:20-30")
		test.S(t).ExpectEquals(gtidSet.GtidEntries[1].String(), "00020194-3333-3333-3333-333333333333:7-8")
	}
	{
		gtidSetVal := "   ,,, , , 00020192-1111-1111-1111-111111111111:20-30,,,, 00020194-3333-3333-3333-333333333333:7-8,,  ,,"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020192-1111-1111-1111-111111111111:20-30")
		test.S(t).ExpectEquals(gtidSet.GtidEntries[1].String(), "00020194-3333-3333-3333-333333333333:7-8")
	}
	{
		gtidSetVal := "   ,,, , ,,  ,,"
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 0)
		test.S(t).ExpectTrue(gtidSet.IsEmpty())
	}
}

func TestRemoveUUID(t *testing.T) {
	gtidSetVal := "00020192-1111-1111-1111-111111111111:20-30, 00020194-3333-3333-3333-333333333333:7-8"
	taggedGtidSetVal := "00020192-1111-1111-1111-111111111111:20-30:tag1:100-114:tag2:1-89," +
		"00020194-3333-3333-3333-333333333333:7-8:tag1:51-200:tag2:44"
	{
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		gtidSet.RemoveUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020192-1111-1111-1111-111111111111:20-30")

		removed := gtidSet.RemoveUUID(`230ea8ea-81e3-11e4-972a-e25ec4bd140a`)
		test.S(t).ExpectFalse(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020192-1111-1111-1111-111111111111:20-30")
	}
	{
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)

		gtidSet.RemoveUUID("00020192-1111-1111-1111-111111111111")
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8")

		gtidSet.RemoveUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectTrue(gtidSet.IsEmpty())
	}

	// Test for tagged gtids
	{
		gtidSet, err := NewOracleGtidSet(taggedGtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)

		gtidSet.RemoveUUID("00020192-1111-1111-1111-111111111111")
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8:tag1:51-200:tag2:44")

		gtidSet.RemoveUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectTrue(gtidSet.IsEmpty())
	}
}

func TestRetainUUID(t *testing.T) {
	gtidSetVal := "00020192-1111-1111-1111-111111111111:20-30, 00020194-3333-3333-3333-333333333333:7-8"
	taggedGtidSetVal := "00020192-1111-1111-1111-111111111111:20-30:tag1:100-114:tag2:1-89," +
		"00020194-3333-3333-3333-333333333333:7-8:tag1:51-200:tag2:44"
	{
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		removed := gtidSet.RetainUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8")

		removed = gtidSet.RetainUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectFalse(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8")

		removed = gtidSet.RetainUUID("230ea8ea-81e3-11e4-972a-e25ec4bd140a")
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 0)
	}
	{
		gtidSet, err := NewOracleGtidSet(taggedGtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		removed := gtidSet.RetainUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8:tag1:51-200:tag2:44")

		removed = gtidSet.RetainUUID("00020194-3333-3333-3333-333333333333")
		test.S(t).ExpectFalse(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8:tag1:51-200:tag2:44")

		removed = gtidSet.RetainUUID("230ea8ea-81e3-11e4-972a-e25ec4bd140a")
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 0)
	}
}

func TestRetainUUIDs(t *testing.T) {
	gtidSetVal := "00020192-1111-1111-1111-111111111111:20-30, 00020194-3333-3333-3333-333333333333:7-8"
	taggedGtidSetVal := "00020192-1111-1111-1111-111111111111:20-30:tag1:100-114:tag2:1-89," +
		"00020193-3333-3333-3333-333333333333:7-8:tag1:51-200:tag2:44," +
		"00020194-3333-3333-3333-333333333333:7-8:tag4:151-209:tag2:44-47"
	{
		gtidSet, err := NewOracleGtidSet(gtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		removed := gtidSet.RetainUUIDs([]string{"00020194-3333-3333-3333-333333333333", "00020194-5555-5555-5555-333333333333"})
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8")

		removed = gtidSet.RetainUUIDs([]string{"00020194-3333-3333-3333-333333333333", "00020194-5555-5555-5555-333333333333"})
		test.S(t).ExpectFalse(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8")

		removed = gtidSet.RetainUUIDs([]string{"230ea8ea-81e3-11e4-972a-e25ec4bd140a"})
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 0)
	}
	{
		gtidSet, err := NewOracleGtidSet(taggedGtidSetVal)
		test.S(t).ExpectNil(err)

		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 3)
		removed := gtidSet.RetainUUIDs([]string{"00020192-1111-1111-1111-111111111111", "00020194-3333-3333-3333-333333333333"})
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 2)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020192-1111-1111-1111-111111111111:20-30:tag1:100-114:tag2:1-89")
		test.S(t).ExpectEquals(gtidSet.GtidEntries[1].String(), "00020194-3333-3333-3333-333333333333:7-8:tag4:151-209:tag2:44-47")

		removed = gtidSet.RetainUUIDs([]string{"00020194-3333-3333-3333-333333333333", "00020195-5555-5555-5555-333333333333"})
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 1)
		test.S(t).ExpectEquals(gtidSet.GtidEntries[0].String(), "00020194-3333-3333-3333-333333333333:7-8:tag4:151-209:tag2:44-47")

		removed = gtidSet.RetainUUIDs([]string{"230ea8ea-81e3-11e4-972a-e25ec4bd140a"})
		test.S(t).ExpectTrue(removed)
		test.S(t).ExpectEquals(len(gtidSet.GtidEntries), 0)
	}
}

func TestSharedUUIDs(t *testing.T) {
	gtidSetVal := "00020192-1111-1111-1111-111111111111:20-30, 00020194-3333-3333-3333-333333333333:7-8"
	gtidSet, err := NewOracleGtidSet(gtidSetVal)
	test.S(t).ExpectNil(err)
	{
		otherSet, err := NewOracleGtidSet("00020194-3333-3333-3333-333333333333:7-8,230ea8ea-81e3-11e4-972a-e25ec4bd140a:1-2")
		test.S(t).ExpectNil(err)
		{
			shared := gtidSet.SharedUUIDs(otherSet)
			test.S(t).ExpectEquals(len(shared), 1)
			test.S(t).ExpectEquals(shared[0], "00020194-3333-3333-3333-333333333333")
		}
		{
			shared := otherSet.SharedUUIDs(gtidSet)
			test.S(t).ExpectEquals(len(shared), 1)
			test.S(t).ExpectEquals(shared[0], "00020194-3333-3333-3333-333333333333")
		}
	}
	{
		otherSet, err := NewOracleGtidSet("00020194-4444-4444-4444-333333333333:7-8,230ea8ea-81e3-11e4-972a-e25ec4bd140a:1-2")
		test.S(t).ExpectNil(err)
		{
			shared := gtidSet.SharedUUIDs(otherSet)
			test.S(t).ExpectEquals(len(shared), 0)
		}
		{
			shared := otherSet.SharedUUIDs(gtidSet)
			test.S(t).ExpectEquals(len(shared), 0)
		}
	}
	{
		otherSet, err := NewOracleGtidSet("00020194-3333-3333-3333-333333333333:7-8,00020192-1111-1111-1111-111111111111:1-2")
		test.S(t).ExpectNil(err)
		{
			shared := gtidSet.SharedUUIDs(otherSet)
			test.S(t).ExpectEquals(len(shared), 2)
		}
		{
			shared := otherSet.SharedUUIDs(gtidSet)
			test.S(t).ExpectEquals(len(shared), 2)
		}
	}
}
