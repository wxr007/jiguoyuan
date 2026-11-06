package core

import (
	"strings"
	"testing"
)

func TestHuangguoAIPlayAddressAcceptsEnvelopeAndFlatPayload(t *testing.T) {
	base := "https://huangguoai.com"
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "edge envelope with relative path",
			body: `{"status":1,"msg":"","data":{"video_url":"/hls/12/ep1.m3u8"}}`,
			want: "https://huangguoai.com/hls/12/ep1.m3u8",
		},
		{
			name: "flat payload with absolute url",
			body: `{"video_url":"https://cdn.example.com/a/b.m3u8?token=1"}`,
			want: "https://cdn.example.com/a/b.m3u8?token=1",
		},
		{
			name: "play url alias",
			body: `{"data":{"play_url":"https://cdn.example.com/a/c.mp4"}}`,
			want: "https://cdn.example.com/a/c.mp4",
		},
		{
			name: "extensionless absolute url",
			body: `{"data":{"video_url":"https://cdn.example.com/stream/abc"}}`,
			want: "https://cdn.example.com/stream/abc",
		},
	}
	for _, item := range cases {
		address, err := huangguoAIPlayAddress(item.body, base)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", item.name, err)
		}
		if address != item.want {
			t.Fatalf("%s: got %q want %q", item.name, address, item.want)
		}
	}
}

func TestHuangguoAIChapterKeepsLegacyKeys(t *testing.T) {
	base := "https://huangguoai.com"
	first := huangguoAIChapter(base, "12", 1, "")
	if first.ID != providerChapterID(sourceHuangguoAI, "12", "auto-1") ||
		first.PageURL != "https://huangguoai.com/video/12/" ||
		first.Title != "第1集" {
		t.Fatalf("first chapter: %+v", first)
	}
	third := huangguoAIChapter(base, "12", 3, "第 3 集")
	if third.ID != providerChapterID(sourceHuangguoAI, "12", "ep-3") ||
		third.PageURL != "https://huangguoai.com/video/12/ep-3/" ||
		third.Title != "第 3 集" ||
		string(third.CurrentEpisode) != "3" {
		t.Fatalf("third chapter: %+v", third)
	}
}

func TestHuangguoAIPlayAddressReportsSiteFailures(t *testing.T) {
	base := "https://huangguoai.com"
	if _, err := huangguoAIPlayAddress(`{"status":0,"msg":"播放已下架"}`, base); err == nil ||
		!strings.Contains(err.Error(), "播放已下架") {
		t.Fatalf("expected site message, got %v", err)
	}
	if _, err := huangguoAIPlayAddress(`{"status":1,"data":{}}`, base); err == nil {
		t.Fatal("expected error for empty payload")
	}
	if _, err := huangguoAIPlayAddress(`<html></html>`, base); err == nil {
		t.Fatal("expected error for non-json payload")
	}
}
