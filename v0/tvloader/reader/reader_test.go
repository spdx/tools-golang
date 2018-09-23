// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reader

import "testing"

func TestCanGetTVListWithFinalize(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	tvList, err := reader.finalize()
	if err != nil {
		t.Errorf("got error when calling finalize: %v", err)
	}
	if len(tvList) != 1 || tvList[0].tag != "Tag" || tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", tvList)
	}
}

func TestCanGetTVListIncludingMultilineWithFinalize(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:<text>value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	err = reader.readNextLine("rest of value</text>")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	tvList, err := reader.finalize()
	if err != nil {
		t.Errorf("got error when calling finalize: %v", err)
	}
	if len(tvList) != 1 || tvList[0].tag != "Tag" || tvList[0].value != "value\nrest of value" {
		t.Errorf("got invalid tag/value list: %v", tvList)
	}
}

func TestCannotFinalizeIfInMidtextState(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:<text>value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	_, err = reader.finalize()
	if err == nil {
		t.Errorf("should have gotten error when calling finalize midtext")
	}
}

func TestCurrentLineIncreasesOnEachReadCall(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}

	reader.currentLine = 23
	err = reader.readNextLine("Tag:value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if reader.currentLine != 24 {
		t.Errorf("expected %d for currentLine, got %d", 23, reader.currentLine)
	}
}

func TestReadyCanReadSingleTagValue(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestReadyCanStripWhitespaceFromValue(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:   value	  	 ")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
}

func TestReadyCannotReadLineWithNoColon(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("No colon should be an error")
	if err == nil {
		t.Errorf("should have gotten error when calling readNextLine")
	}
}

func TestReadyTextTagSwitchesToMidtext(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag: <text>This begins a multiline value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if !reader.midtext {
		t.Errorf("expected midtext to be true, got false")
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "Tag" {
		t.Errorf("expected %s for currentTag, got %s", "Tag", reader.currentTag)
	}
	if reader.currentValue != "This begins a multiline value\n" {
		t.Errorf("expected %s for currentValue, got %s", "This begins a multiline value\n", reader.currentValue)
	}
}

func TestReadyTextTagAndClosingTagInOneLineFinishesRead(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag: <text>Just one line</text>")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "Just one line" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestCanReadMultilineTextAcrossThreeLines(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag: <text>This value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	err = reader.readNextLine("is three")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	err = reader.readNextLine("lines long</text>")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "This value\nis three\nlines long" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 3 {
		t.Errorf("expected %d for currentLine, got %d", 3, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestMidtextContinuesIfNoClosingText(t *testing.T) {
	reader := &tvReader{}
	reader.midtext = true
	reader.currentLine = 1
	reader.currentTag = "Multiline"
	reader.currentValue = "First line\n"

	err := reader.readNextLine("Second line")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if !reader.midtext {
		t.Errorf("expected midtext to be true, got false")
	}
	if reader.currentLine != 2 {
		t.Errorf("expected %d for currentLine, got %d", 2, reader.currentLine)
	}
	if reader.currentTag != "Multiline" {
		t.Errorf("expected %s for currentTag, got %s", "Multiline", reader.currentTag)
	}
	if reader.currentValue != "First line\nSecond line\n" {
		t.Errorf("expected %s for currentValue, got %s", "First line\nSecond line\n", reader.currentValue)
	}
}

func TestMidtextFinishesIfReachingClosingText(t *testing.T) {
	reader := &tvReader{}
	reader.midtext = true
	reader.currentLine = 1
	reader.currentTag = "Multiline"
	reader.currentValue = "First line\n"

	err := reader.readNextLine("Second line</text>")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Multiline" || reader.tvList[0].value != "First line\nSecond line" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 2 {
		t.Errorf("expected %d for currentLine, got %d", 2, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestReadyIgnoresCommentLines(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("# this is a comment")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestMidtextIncludesCommentLines(t *testing.T) {
	reader := &tvReader{}
	reader.midtext = true
	reader.currentLine = 1
	reader.currentTag = "Multiline"
	reader.currentValue = "First line\n"

	err := reader.readNextLine("# This is part of multiline text")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if !reader.midtext {
		t.Errorf("expected midtext to be true, got false")
	}
	if reader.currentLine != 2 {
		t.Errorf("expected %d for currentLine, got %d", 2, reader.currentLine)
	}
	if reader.currentTag != "Multiline" {
		t.Errorf("expected %s for currentTag, got %s", "Multiline", reader.currentTag)
	}
	if reader.currentValue != "First line\n# This is part of multiline text\n" {
		t.Errorf("expected %s for currentValue, got %s", "First line\n# This is part of multiline text\n", reader.currentValue)
	}
}

func TestReadyIgnoresEmptyLines(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestMidtextIncludesEmptyLines(t *testing.T) {
	reader := &tvReader{}
	reader.midtext = true
	reader.currentLine = 1
	reader.currentTag = "Multiline"
	reader.currentValue = "First line\n"

	err := reader.readNextLine("")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if !reader.midtext {
		t.Errorf("expected midtext to be true, got false")
	}
	if reader.currentLine != 2 {
		t.Errorf("expected %d for currentLine, got %d", 2, reader.currentLine)
	}
	if reader.currentTag != "Multiline" {
		t.Errorf("expected %s for currentTag, got %s", "Multiline", reader.currentTag)
	}
	if reader.currentValue != "First line\n\n" {
		t.Errorf("expected %s for currentValue, got %s", "First line\n\n", reader.currentValue)
	}
}

func TestReadyIgnoresWhitespaceOnlyLines(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("   \t\t\t ")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestMidtextIncludesWhitespaceOnlyLines(t *testing.T) {
	reader := &tvReader{}
	reader.midtext = true
	reader.currentLine = 1
	reader.currentTag = "Multiline"
	reader.currentValue = "First line\n"

	err := reader.readNextLine("     \t\t ")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if !reader.midtext {
		t.Errorf("expected midtext to be true, got false")
	}
	if reader.currentLine != 2 {
		t.Errorf("expected %d for currentLine, got %d", 2, reader.currentLine)
	}
	if reader.currentTag != "Multiline" {
		t.Errorf("expected %s for currentTag, got %s", "Multiline", reader.currentTag)
	}
	if reader.currentValue != "First line\n     \t\t \n" {
		t.Errorf("expected %s for currentValue, got %s", "First line\n     \t\t \n", reader.currentValue)
	}
}

func TestReadyIgnoresSpacesBeforeTag(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("    \t Tag:value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestReadyIgnoresSpacesBeforeCommentLines(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("     \t\t  # this is a comment")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}

	if len(reader.tvList) != 0 {
		t.Errorf("expected empty tag/value list, got %v", reader.tvList)
	}
	if reader.midtext {
		t.Errorf("expected midtext to be false, got true")
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestReadyIgnoresSpacesBetweenTagAndColon(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag   \t :value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestReadyIgnoresSpacesBetweenColonAndValue(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:    \t value")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}

func TestReadyIgnoresSpacesAfterEndOfValue(t *testing.T) {
	reader := &tvReader{}
	err := reader.readNextLine("Tag:value   \t  ")
	if err != nil {
		t.Errorf("got error when calling readNextLine: %v", err)
	}
	if len(reader.tvList) != 1 || reader.tvList[0].tag != "Tag" || reader.tvList[0].value != "value" {
		t.Errorf("got invalid tag/value list: %v", reader.tvList)
	}
	if reader.currentLine != 1 {
		t.Errorf("expected %d for currentLine, got %d", 1, reader.currentLine)
	}
	if reader.currentTag != "" {
		t.Errorf("expected empty string for currentTag, got %s", reader.currentTag)
	}
	if reader.currentValue != "" {
		t.Errorf("expected empty string for currentValue, got %s", reader.currentValue)
	}
}
