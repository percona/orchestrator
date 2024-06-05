/*
   Copyright 2015 Shlomi Noach, courtesy Booking.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package inst

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regex Pattern to match a Singe Value Interval e.g. 1 in uuid:1
	singleValueInterval = regexp.MustCompile("^([0-9]+)$")

	// Regex Pattern to match a Multi Value Interval e.g. 1-5 in uuid:1-5
	multiValueInterval = regexp.MustCompile("^([0-9]+)[-]([0-9]+)$")

	// Regex Pattern to match GTID tags
	// Tag must start with a letter e.g. tag1 in uuid:tag1:1-5
	tagRegex = regexp.MustCompile("^[a-z_][a-z0-9_]{0,31}")
)

type TagInterval struct {
	Tag      string   // tag name
	Interval []string // intervals
}

// OracleGtidSetEntry represents an entry in a set of GTID ranges,
// for example,
//
//	Valid Formats:
//		- "316d193c-70e5-11e5-adb2-ecf4bb2262ff:1-8935:8984-6124596" (may include gaps)
//		- "321f5c0d-70e5-11e5-adb2-ecf4bb2262ff:1-56457" (no gaps)
//		- "230ea8ea-81e3-11e4-972a-e25ec4bd140a:1-10539:tag1:1-2474" (tagged intervals)
//		- "230ea8ea-81e3-11e4-972a-e25ec4bd140a:1-5139:5780-6317:tag1:1-2474:3201-4157" (tagged intervals and may have gaps)
type OracleGtidSetEntry struct {
	UUID      string
	DefaultIv string        // default (untagged) interval
	TaggedIv  []TagInterval // tagged intervals
}

func ParseOracleGtidSetEntry(gtidRangeString string) (*OracleGtidSetEntry, error) {

	// Split the string into two parts: UUID part and the non-UUID part
	gtid_str := strings.SplitN(gtidRangeString, ":", 2)

	// Sanity check
	if len(gtid_str) != 2 {
		return nil, fmt.Errorf("Cannot parse OracleGtidSetEntry from %s", gtidRangeString)
	}

	if gtid_str[0] == "" {
		return nil, fmt.Errorf("Unexpected UUID: %s", gtid_str[0])
	}

	if gtid_str[1] == "" {
		return nil, fmt.Errorf("Unexpected GTID range: %s", gtid_str[1])
	}

	// UUID is the first part
	uuid := gtid_str[0]

	// Split the non-UUID parts into multiple blocks
	s := strings.SplitN(gtid_str[1], ":", -1)

	// Initialize the tag and interval
	var default_iv string     // Default interval
	var tag_ivs []TagInterval // Full tagged interval
	var ti TagInterval        // Current tag interval

	for i := range s {

		// If it is a GTID tag
		if tagRegex.MatchString(s[i]) {

			if (ti.Tag != "") && (len(ti.Interval) == 0) {
				// If the tag is already set and we got another tag
				return nil, fmt.Errorf("Invalid format")
			} else if (ti.Tag == "") && (len(ti.Interval) != 0) {
				// If the tag is not set and we already have the interval set
				return nil, fmt.Errorf("Invalid format")
			} else {
				// Now process the new tag
				ti.Tag = s[i]
				// Reset iv for the current tag
				ti.Interval = nil
				// Append the new tag to tag_ivs
				tag_ivs = append(tag_ivs, ti)
			}
		} else {
			// If it is an GTID interval
			if singleValueInterval.MatchString(s[i]) || multiValueInterval.MatchString(s[i]) {
				// If it is an empty tag, add it to default interval
				if len(ti.Tag) == 0 {
					default_iv += ":" + s[i]
				} else {
					// If tag is already set, add it to the tag interval
					ti.Interval = append(ti.Interval, s[i])
					tag_ivs[len(tag_ivs)-1].Interval = append(tag_ivs[len(tag_ivs)-1].Interval, s[i])
				}
			} else {
				// Regex failed, invalid format
				return nil, fmt.Errorf("Invalid format")
			}
		}
	}

	// If the interval of the last tag is empty, then it is an invalid format
	// eg: "UUID:1-5139::tag1:"
	if (ti.Tag != "") && (len(tag_ivs[len(tag_ivs)-1].Interval) == 0) {
		return nil, fmt.Errorf("Invalid format")
	}

	// Don't append ':' for the first interval in the default set
	if len(default_iv) != 0 {
		default_iv, _ = strings.CutPrefix(default_iv, ":")
	}

	entry := OracleGtidSetEntry{UUID: uuid, DefaultIv: default_iv, TaggedIv: tag_ivs}

	return &entry, nil
}

// NewOracleGtidSetEntry parses a single entry text
func NewOracleGtidSetEntry(gtidRangeString string) (*OracleGtidSetEntry, error) {
	gtidRangeString = strings.TrimSpace(gtidRangeString)

	gtidRange, error := ParseOracleGtidSetEntry(gtidRangeString)
	if gtidRange == nil {
		return nil, error
	}

	return gtidRange, nil
}

// String() returns a user-friendly string representation of this entry
func (this *OracleGtidSetEntry) String() string {

	var res string

	// UUID is always added in the beginning of the Gtid_set
	res += this.UUID

	// Default ranges are always added immediately after the UUID
	if len(this.DefaultIv) != 0 {
		res += ":" + this.DefaultIv
	}

	// Tagged ranges are added in the end
	for _, v := range this.TaggedIv {
		res += ":" + v.Tag
		if len(v.Interval) != 0 {
			res += ":" + strings.Join(v.Interval, ":")
		}
	}
	return res
}

/*
Explode() returns a list of individual gtids that are represented by this entry.

Example:
Explode of the GTID set "48ebed33-0d12-11ef-a3ec-ac198e4551c8:1-3:7:tag1:1-2:10-12:tag2:74-75:78:81"
shall return the following

48ebed33-0d12-11ef-a3ec-ac198e4551c8:1
48ebed33-0d12-11ef-a3ec-ac198e4551c8:2
48ebed33-0d12-11ef-a3ec-ac198e4551c8:3
48ebed33-0d12-11ef-a3ec-ac198e4551c8:7
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag1:1
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag1:2
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag1:10
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag1:11
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag1:12
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag2:74
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag2:75
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag2:78
48ebed33-0d12-11ef-a3ec-ac198e4551c8:tag2:81
*/
func (this *OracleGtidSetEntry) Explode() (result [](*OracleGtidSetEntry)) {

	// Appends the default interval to the result
	var AppendDefaultInterval = func(this *OracleGtidSetEntry) {
		intervals := strings.Split(this.DefaultIv, ":")
		for _, interval := range intervals {
			// Multi-value interval
			if submatch := multiValueInterval.FindStringSubmatch(interval); submatch != nil {
				intervalStart, _ := strconv.Atoi(submatch[1])
				intervalEnd, _ := strconv.Atoi(submatch[2])
				for i := intervalStart; i <= intervalEnd; i++ {
					result = append(result, &OracleGtidSetEntry{UUID: this.UUID, DefaultIv: fmt.Sprintf("%d", i)})
				}
			} else if submatch := singleValueInterval.FindStringSubmatch(interval); submatch != nil {
				// Single-value interval
				result = append(result, &OracleGtidSetEntry{UUID: this.UUID, DefaultIv: interval})
			}
		}
	}

	// Appends tagged intervals to the result
	var AppendTaggedInterval = func(tag string, interval string) {
		intervals := strings.Split(interval, ":")
		for _, interval := range intervals {
			// Multi-value interval
			if submatch := multiValueInterval.FindStringSubmatch(interval); submatch != nil {
				intervalStart, _ := strconv.Atoi(submatch[1])
				intervalEnd, _ := strconv.Atoi(submatch[2])
				for i := intervalStart; i <= intervalEnd; i++ {

					ti := TagInterval{
						Tag:      tag,
						Interval: []string{fmt.Sprintf("%d", i)}}
					taggedIv := []TagInterval{ti}

					entry := OracleGtidSetEntry{UUID: this.UUID, TaggedIv: taggedIv}
					result = append(result, &entry)
				}
			} else if submatch := singleValueInterval.FindStringSubmatch(interval); submatch != nil {

				// Single-value interval
				ti := TagInterval{
					Tag:      tag,
					Interval: []string{interval}}
				taggedIv := []TagInterval{ti}

				entry := OracleGtidSetEntry{UUID: this.UUID, TaggedIv: taggedIv}
				result = append(result, &entry)
			}
		}
	}

	// Process default interval first
	AppendDefaultInterval(this)

	// Process tagged intervals next
	for _, v := range this.TaggedIv {
		for _, iv := range v.Interval {
			AppendTaggedInterval(v.Tag, iv)
		}
	}
	return result
}
