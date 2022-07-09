package service

import (
	"bufio"
	"cc_blog/biz/common"
	"cc_blog/biz/model"
	"cc_blog/biz/util"
	"cc_blog/config"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const SegLine = "<br line/>"

type BlogServiceImpl struct {
}

func (b *BlogServiceImpl) Refresh() error {
	return util.ClearAllCache()
}

func (b *BlogServiceImpl) GetArticles(path string) (*model.Articles, error) {
	list := model.Articles{}
	err := util.GetCache(common.Articles+path, &list)
	if err == nil {
		return &list, nil
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			article, err := b.GetArticleBase(path)
			if err != nil {
				return err
			}
			list = append(list, article)
		}
		return nil
	})
	if err != nil {
		return &list, err
	}
	sort.Sort(&list)
	_ = util.SetCache(common.Articles+path, list)
	return &list, nil
}

func (b *BlogServiceImpl) GetCategories() (*model.Categories, error) {
	categories := model.Categories{}
	err := util.GetCache(common.Categories, &categories)
	if err == nil {
		return &categories, nil
	}
	err = filepath.Walk(config.Conf.DocDir, func(path string, info os.FileInfo, err error) error {
		if path == config.Conf.DocDir {
			return nil
		}
		if info.IsDir() {
			categories = append(categories, &model.Category{Name: info.Name(), Path: info.Name(), Number: 0})
		}
		if filepath.Ext(path) == ".md" {
			categories[len(categories)-1].Number++
		}
		return nil
	})
	_ = util.SetCache(common.Categories, categories)
	return &categories, err
}

func (b *BlogServiceImpl) GetNavbar() (*model.Navbar, error) {
	navbar := &model.Navbar{}
	err := util.GetCache(common.Navbar, navbar)
	if err == nil {
		return navbar, nil
	}
	file, err := os.Open(filepath.Join(config.Conf.RootDir, "navbar.json"))
	if err != nil {
		return navbar, err
	}
	content, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(content, navbar)
	_ = util.SetCache(common.Navbar, navbar)
	return navbar, nil
}

func (b *BlogServiceImpl) GetSlider() (*model.Slider, error) {
	slider := &model.Slider{}
	err := util.GetCache(common.Slider, slider)
	if err == nil {
		return slider, nil
	}
	file, err := os.Open(filepath.Join(config.Conf.RootDir, "app.json"))
	if err != nil {
		return slider, err
	}
	content, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(content, slider)
	_ = util.SetCache(common.Slider, slider)
	return slider, nil
}

func (b *BlogServiceImpl) GetArticleBase(name string) (*model.Article, error) {
	result := &model.Article{}
	if file, err := os.Open(name); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		hasBegin := false
		blockNum := 0
		blockContent := ""
		line := ""
		for {
			line, err = reader.ReadString('\n')
			if err != nil {
				break
			}
			if line != "" && !hasBegin {
				hasBegin = true
				str := strings.TrimSpace(line)
				ok, err := regexp.MatchString("^(#{1,3})\\s+(.*)", str)
				if ok && err == nil {
					fly := regexp.MustCompile("^(#{1,3})\\s+")
					titles := fly.FindStringSubmatch(str)
					result.Title = strings.TrimSpace(string(str)[len(titles[1]):])
					continue
				}
			}
			if strings.TrimSpace(line) == SegLine {
				blockNum++
				if blockNum == 2 {
					break
				}
				continue
			}
			if blockNum == 1 {
				blockContent += line
			}
		}
		if result.Title == "" {
			suffix := filepath.Ext(name)
			result.Title = strings.TrimSuffix(filepath.Base(name), suffix)
		}
		result.Description = blockContent
		fileInfo, err := file.Stat()
		result.CreatedTime = util.GetFileCreatedTime(fileInfo)
		if err != nil {
			result.UpdatedTime = time.Now()
		} else {
			result.UpdatedTime = fileInfo.ModTime()
		}
		dirs := strings.Split(name, "/")
		result.Category = dirs[len(dirs)-2]
		if strings.Contains(name, config.Conf.DocDir) {
			result.Path = string([]byte(name)[len(config.Conf.DocDir) : len(name)-3])
			result.Type = 1
		}
		return result, err
	} else {
		return result, err
	}
}

func (b *BlogServiceImpl) GetArticleDetail(name string) (*model.Article, error) {
	result := &model.Article{}
	err := util.GetCache(common.Article+name, result)
	if err == nil {
		return result, nil
	}
	if file, err := os.Open(name); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		hasBegin := false
		blockNum := 0
		descContent := ""
		body := ""
		line := ""
		for {
			line, err = reader.ReadString('\n')
			if err != nil {
				break
			}
			if line != "" && !hasBegin {
				hasBegin = true
				str := strings.TrimSpace(line)
				ok, err := regexp.MatchString("^(#{1,3})\\s+(.*)", str)
				if ok && err == nil {
					fly := regexp.MustCompile("^(#{1,3})\\s+")
					titles := fly.FindStringSubmatch(str)
					result.Title = strings.TrimSpace(str[len(titles[1]):])
					continue
				}
			}
			if strings.TrimSpace(line) == SegLine {
				blockNum++
				continue
			}
			if blockNum == 1 {
				descContent += line
			}
			body += line
		}
		if result.Title == "" {
			suffix := filepath.Ext(name)
			result.Title = strings.TrimSuffix(filepath.Base(name), suffix)
		}
		result.Description = descContent
		result.Body = body
		fileInfo, err := file.Stat()
		result.CreatedTime = util.GetFileCreatedTime(fileInfo)
		if err != nil {
			result.UpdatedTime = time.Now()
		} else {
			result.UpdatedTime = fileInfo.ModTime()
		}
		dirs := strings.Split(name, "/")
		result.Category = dirs[len(dirs)-2]
		if strings.Contains(name, config.Conf.DocDir) {
			result.Path = string([]byte(name)[len(config.Conf.DocDir) : len(name)-3])
			result.Type = 1
		}
		_ = util.SetCache(common.Article+name, result)
		return result, err
	} else {
		return result, err
	}
}
