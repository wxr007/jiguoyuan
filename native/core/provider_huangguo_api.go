package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const huangguoAIEdgePath = "/api/edge/request"

func (d *Downloader) fetchHuangguoAIPlayURL(ctx context.Context, sourceID string, episode int) (string, error) {
	if !rankingSourceID.MatchString(sourceID) || episode < 1 {
		return "", errors.New("黄果集数无效")
	}
	base := d.providerBaseURL(sourceHuangguoAI)
	path := fmt.Sprintf("/api/videos/%s/play?ep=%d&page_texts=0", url.PathEscape(sourceID), episode)
	var directErr error
	if body, err := d.fetchProviderText(ctx, base+path, base+"/"); err != nil {
		directErr = err
	} else if address, parseErr := huangguoAIPlayAddress(body, base); parseErr != nil {
		directErr = parseErr
	} else {
		return address, nil
	}
	envelope, err := json.Marshal(map[string]string{
		"method":       http.MethodGet,
		"path":         path,
		"body":         "",
		"content_type": "",
	})
	if err != nil {
		return "", err
	}
	body, err := d.fetchProviderJSON(ctx, base+huangguoAIEdgePath, base+"/", envelope)
	if err != nil {
		if directErr != nil {
			return "", fmt.Errorf("黄果播放地址获取失败：%v；边缘接口同样失败：%w", directErr, err)
		}
		return "", err
	}
	address, err := huangguoAIPlayAddress(body, base)
	if err != nil {
		if directErr != nil {
			return "", fmt.Errorf("黄果播放地址获取失败：%v", directErr)
		}
		return "", err
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
