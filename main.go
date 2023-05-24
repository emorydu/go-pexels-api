package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func (c *Client) SearchPhotos(query string, perPage, page int32) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) CuratedPhotos(perPage, page int32) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CuratedResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Photo
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int31n(1001)
	result, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}

func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token)
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return resp, nil
	} else {
		c.RemainingTimes = int32(times)
	}

	return resp, nil
}

func (c *Client) SearchVideo(query string, perPage, page int32) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result VideoSearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) PopularVideo(perPage, page int32) (*PopularVideo, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PopularVideo
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetRandomVideo() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int31n(1001)
	result, err := c.PopularVideo(1, randNum)
	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}
	return nil, err
}

func (c *Client) GetRemainingRequestInThisMonth() int32 {
	return c.RemainingTimes
}

func NewClient(token string) *Client {
	h := http.Client{}

	return &Client{
		Token: token,
		hc:    h,
	}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	URL             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerURL string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}

type Video struct {
	Id            int32          `json:"id"`
	Width         int32          `json:"width"`
	Height        int32          `json:"height"`
	URL           string         `json:"url"`
	Image         string         `json:"image"`
	FullRes       any            `json:"full_res"`
	Duration      float64        `json:"duration"`
	VideoFiles    []VideoFile    `json:"video_files"`
	VideoPictures []VideoPicture `json:"video_pictures"`
}

type PopularVideo struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	URL          string  `json:"url"`
	Videos       []Video `json:"videos"`
}

type VideoFile struct {
	Id       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

type VideoPicture struct {
	Id      int32  `json:"id"`
	Picture string `json:"picture"`
	Nr      int32  `json:"nr"`
}

func main() {

	os.Setenv("PexelsToken", "WtPWUzdzSYXkK8R3UGb2xiJskBQe09ewrz5uq5ycFq9crCcSPw9fGhg4")

	var token = os.Getenv("PexelsToken")

	var c = NewClient(token)

	result, err := c.SearchPhotos("waves", 80, 1)
	if err != nil {
		_ = fmt.Errorf("search error: %v", err)
	}
	if result.Page == 0 {
		_ = fmt.Errorf("search result wrong")
	}
	rj, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	os.WriteFile("test.json", rj, 0666)
	//var wg sync.WaitGroup
	//now := time.Now()
	//for _, v := range result.Photos {
	//	wg.Add(1)
	//	go getImage(path.Base(v.Src.Original), v.Src.Original, &wg)
	//}
	//wg.Wait()
	//fmt.Printf("%f\n", time.Since(now).Minutes())

}

func getImage(filename string, url string, wg *sync.WaitGroup) error {
	fmt.Println("prepare downloads file: ", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if body, err := io.ReadAll(resp.Body); err != nil {
		return err
	} else {

		if err := os.WriteFile(filename, body, 0666); err != nil {
			return err
		}
		wg.Done()
	}
	return nil
}
