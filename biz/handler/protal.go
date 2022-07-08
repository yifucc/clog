package handler

import (
	"cc_blog/biz/service"
	"cc_blog/config"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kataras/iris/v12"
	ssh2 "golang.org/x/crypto/ssh"
	"log"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

var blogService service.BlogService = &service.BlogServiceImpl{}
var auth transport.AuthMethod

func init() {
	switch config.Conf.CodeRepoAuthType {
	case config.PrivateKeyType:
		privateKeys, _ := ssh.NewPublicKeysFromFile("git", fmt.Sprintf("%s/%s", config.GetConfigDir(), "private_key.perm"), "")
		privateKeys.HostKeyCallback = func(hostname string, remote net.Addr, key ssh2.PublicKey) error {
			return nil
		}
		auth = privateKeys
	case config.PasswordType:
		auth = &http.BasicAuth{
			Username: config.Conf.CodeRepoAuthUser,
			Password: config.Conf.CodeRepoAuthPassword,
		}
	}
}

func Site(ctx iris.Context) {
	ctx.ViewData("title", config.Conf.Name)
	ctx.ViewData("slogan", config.Conf.Slogan)
	ctx.ViewData("icp", config.Conf.Icp)
	ctx.Next()
}

func Slider(ctx iris.Context) {
	slider, _ := blogService.GetSlider()
	data, _ := json.Marshal(slider)
	ctx.ViewData("slider", string(data))
	ctx.Next()
}

func Categories(ctx iris.Context) {
	categories, _ := blogService.GetCategories()
	data, _ := json.Marshal(categories)
	ctx.ViewData("cate", string(data))
	ctx.Next()
}

func Navbar(ctx iris.Context) {
	navbar, _ := blogService.GetNavbar()
	data, _ := json.Marshal(navbar)
	ctx.ViewData("menu", string(data))
	ctx.Next()
}

func ArticleList(ctx iris.Context) {
	res, _ := blogService.GetArticles(filepath.Join(config.Conf.DocDir, ctx.Params().Get("path")))
	pageParam := ctx.URLParam("page")
	page, err := strconv.Atoi(pageParam)
	if pageParam == "" || err != nil {
		page = 1
	}
	limit := config.Conf.PageLimit
	if len(*res) > limit {
		begin := (page - 1) * limit
		end := page * limit
		if begin > (len(*res) - 1) {
			begin = len(*res)
			end = len(*res)
		}
		if end > (len(*res) - 1) {
			end = len(*res)
		}
		*res = (*res)[begin:end]
	}
	data, _ := json.Marshal(res)
	ctx.ViewData("data", string(data))
	ctx.ViewData("page", page)
	ctx.ViewData("records", len(*res))
	ctx.ViewData("perPage", limit)
	ctx.View("list.html")
}

func ArticleDetail(ctx iris.Context) {
	dir, _ := url.QueryUnescape(ctx.Params().Get("path"))
	name, _ := url.QueryUnescape(ctx.Params().Get("name"))
	path := filepath.Join(config.Conf.DocDir, dir, name)
	res, _ := blogService.GetArticleDetail(path + ".md")
	data, _ := json.Marshal(res)
	ctx.ViewData("data", string(data))
	ctx.View("marked_view.html")
}

func Webhook(ctx iris.Context) {
	time := ctx.GetHeader("X-Gitee-Timestamp")
	token := ctx.GetHeader("X-Gitee-Token")
	key := fmt.Sprintf("%s\n%s", time, config.Conf.CodeRepoSecret)
	sha := sha256.New()
	sha.Write([]byte(key))
	sign := base64.URLEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString(sha.Sum(nil))))
	if token != sign {
		ctx.StatusCode(400)
		return
	}
	updateResource()
	ctx.StatusCode(200)
}

func updateResource() {
	repo, err := git.PlainOpen(config.Conf.RootDir)
	if err != nil {
		log.Printf("open resource dir error: %s", err.Error())
		repo, err = git.PlainClone(config.Conf.RootDir, false, &git.CloneOptions{
			URL:  config.Conf.CodeRepoAddress,
			Auth: auth,
		})
		if err != nil {
			log.Printf("clone resource dir error: %s", err.Error())
		}
	} else {
		worktree, _ := repo.Worktree()
		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
			Force:      true,
			Auth:       auth,
		})
		if err != nil {
			log.Printf("pull resource dir error: %s", err.Error())
		}
	}
}

func checkSign(sign string, body []byte) bool {
	if len(sign) != 45 || !strings.HasPrefix(sign, "sha1=") {
		return false
	}
	secret := []byte(config.Conf.CodeRepoSecret)
	mac := hmac.New(sha1.New, secret)
	mac.Write(body)
	key := mac.Sum(nil)
	signature := make([]byte, 20)
	hex.Decode(signature, []byte(sign[5:]))
	return hmac.Equal(signature, key)
}
