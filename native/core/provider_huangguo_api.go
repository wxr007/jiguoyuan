package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const huangguoAIEdgePath = "/api/edge/request"

func (d *Downloader) huangguoAIAPI(ctx context.Context, path string, out any) error {
	envelope, err := json.Marshal(map[string]string{
		"method":       http.MethodGet,
		"path":         path,
		"body":         "",
		"content_type": "",
	})
	if err != nil {
		return err
	}
	base := strings.TrimRight(d.providerBaseURL(sourceHuangguoAI), "/")
	body, err := d.fetchProviderJSON(ctx, base+huangguoAIEdgePath, base+"/", envelope)
	if err != nil {
		return err
	}
	var reply struct {
		Status int             `json:"status"`
		Msg    string          `json:"msg"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &reply); err != nil {
		return fmt.Errorf("黄果接口返回格式无效：%w", err)
	}
	if reply.Status != 1 {
		return fmt.Errorf("黄果接口返回失败：%s", firstNonEmpty(reply.Msg, "未知错误"))
	}
	if len(reply.Data) == 0 || string(reply.Data) == "null" {
		return errors.New("黄果接口未返回数据")
	}
	if err := json.Unmarshal(reply.Data, out); err != nil {
		return fmt.Errorf("黄果接口数据格式已变化：%w", err)
	}
	return nil
}

func (d *Downloader) fetchHuangguoAIDetailAPI(ctx context.Context, sourceID string) (Drama, []Chapter, error) {
	var detail map[string]any
	if err := d.huangguoAIAPI(ctx, "/api/videos/"+url.PathEscape(sourceID), &detail); err != nil {
		return Drama{}, nil, err
	}
	base := strings.TrimRight(d.providerBaseURL(sourceHuangguoAI), "/")
	drama, ok := dramaFromAIMap(detail, base+"/", "")
	if !ok || drama.SourceID != sourceID {
		return Drama{}, nil, errors.New("黄果详情与请求剧集不符")
	}
	rows, _ := detail["episodes"].([]any)
	var chapters []Chapter
	for _, value := range rows {
		row, _ := value.(map[string]any)
		episode, err := strconv.Atoi(mapString(row, "ep_num"))
		if err != nil || episode < 1 {
			continue
		}
		chapters = append(chapters, huangguoAIChapter(base, sourceID, episode, mapString(row, "title")))
	}
	if len(chapters) == 0 {
		return Drama{}, nil, errors.New("黄果详情未返回分集")
	}
	sortProviderChapters(chapters)
	return drama, uniqueChapters(chapters), nil
}

func huangguoAIChapter(base, sourceID string, episode int, title string) Chapter {
	key, pageURL := "auto-1", base+"/video/"+sourceID+"/"
	if episode > 1 {
		key = "ep-" + strconv.Itoa(episode)
		pageURL += key + "/"
	}
	return Chapter{
		ID:             providerChapterID(sourceHuangguoAI, sourceID, key),
		Source:         sourceHuangguoAI,
		Title:          firstNonEmpty(strings.TrimSpace(title), fmt.Sprintf("第%d集", episode)),
		PageURL:        pageURL,
		CurrentEpisode: rawEpisode(episode),
	}
}

func (d *Downloader) fetchHuangguoAIPlayURL(ctx context.Context, sourceID string, episode int) (string, error) {
	if !rankingSourceID.MatchString(sourceID) || episode < 1 {
		return "", errors.New("黄果集数无效")
	}
	base := strings.TrimRight(d.providerBaseURL(sourceHuangguoAI), "/")
	path := fmt.Sprintf("/api/videos/%s/play?ep=%d&page_texts=0", url.PathEscape(sourceID), episode)
	var directErr error
	if body, err := d.fetchProviderText(ctx, base+path, base+"/"); err != nil {
		directErr = err
	} else if address, parseErr := huangguoAIPlayAddress(body, base); parseErr == nil {
		return address, nil
	} else {
		directErr = parseErr
	}
	var play struct {
		VideoURL string `json:"video_url"`
		PlayURL  string `json:"play_url"`
		URL      string `json:"url"`
	}
	if err := d.huangguoAIAPI(ctx, path, &play); err != nil {
		if directErr != nil {
			return "", fmt.Errorf("黄果播放地址获取失败：%v；边缘接口同样失败：%w", directErr, err)
		}
		return "", err
	}
	raw := firstNonEmpty(strings.TrimSpace(play.VideoURL), strings.TrimSpace(play.PlayURL), strings.TrimSpace(play.URL))
	if raw == "" {
		if directErr != nil {
			return "", fmt.Errorf("黄果播放地址获取失败：%v", directErr)
		}
		return "", errors.New("黄果播放接口未返回播放地址")
	}
	address := normalizeProviderMediaURL(base+"/", raw)
	if address == "" && isProviderHTTPMediaURL(raw) {
		address = raw
	}
	if address == "" {
		return "", errors.New("黄果播放接口未返回有效播放地址")
	}
	return address, nil
}

func huangguoAIPlayAddress(body, base string) (string, error) {
	var reply struct {
		Status   *int            `json:"status"`
		Msg      string          `json:"msg"`
		Data     json.RawMessage `json:"data"`
		VideoURL string          `json:"video_url"`
	}
	if err := json.Unmarshal([]byte(body), &reply); err != nil {
		return "", fmt.Errorf("黄果播放接口返回格式无效：%w", err)
	}
	if reply.Status != nil && *reply.Status != 1 {
		return "", fmt.Errorf("黄果播放接口返回失败：%s", firstNonEmpty(reply.Msg, "未知错误"))
	}
	raw := strings.TrimSpace(reply.VideoURL)
	if len(reply.Data) > 0 {
		var data struct {
			VideoURL string `json:"video_url"`
			PlayURL  string `json:"play_url"`
			URL      string `json:"url"`
		}
		if json.Unmarshal(reply.Data, &data) == nil {
			raw = firstNonEmpty(strings.TrimSpace(data.VideoURL), strings.TrimSpace(data.PlayURL), strings.TrimSpace(data.URL), raw)
		}
	}
	if raw == "" {
		return "", errors.New("黄果播放接口未返回播放地址")
	}
	address := normalizeProviderMediaURL(base+"/", raw)
	if address == "" && isProviderHTTPMediaURL(raw) {
		address = raw
	}
	if address == "" {
		return "", errors.New("黄果播放接口未返回有效播放地址")
	}
	return address, nil
}
