package main

import (
	"errors"
	"fmt"
	"github.com/bingerambo/go-file-json-server/utils"
	"github.com/fsnotify/fsnotify"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/**
 通过上传服务API(/upload, /file)成功上传的json文件。会触发目录文件扫描，把最新的文件名更新到Cache
Cache.data 用于server动态构建相应API进行处理
*/
type Cache struct {
	// thread safe
	data       *utils.DataSet
	watch      *fsnotify.Watcher
	watchRoot string
}

func NewCache(watchRoot string) *Cache {
	wat, _ := fsnotify.NewWatcher()
	return &Cache{
		data:       utils.NewDataSet("file_cache"),
		watch:      wat,
		watchRoot: watchRoot,
	}
}

func (c *Cache) Boot() {
	c.watchFS(c.watchRoot)
}


//const(
//	CREATE = "CREATE"
//	DELETE = "DELETE"
//)
// Op describes a set of file operations.
type CacheOp uint32

const (
	CREATE CacheOp = 1 << iota
	DELETE
)

//监控目录
func (c *Cache) watchDir(dir string) {
	//通过Walk来遍历目录下的所有子目录
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir(){
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = c.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("start to monitor: ", path)
		} else {
			// 同步处理
			c.syncCacheServer(path, CREATE)
		}

		return nil
	})
	go func() {
		for {
			select {
			case ev := <-c.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("create file: ", ev.Name)
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							c.watch.Add(ev.Name)
							fmt.Println("add to monitor: ", ev.Name)
						}

						// 同步更新处理
						c.syncCacheServer(ev.Name, CREATE)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("write file: ", ev.Name)
						// 同步更新处理
						c.syncCacheServer(ev.Name, CREATE)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("delete file: ", ev.Name)
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							c.watch.Remove(ev.Name)
							fmt.Println("delete to monitor: ", ev.Name)
						}
						// 同步更新处理
						c.syncCacheServer(ev.Name, DELETE)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("rename file: ", ev.Name)
						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						c.watch.Remove(ev.Name)
						// 同步更新处理
						c.syncCacheServer(ev.Name, DELETE)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						fmt.Println("chmod file: ", ev.Name)
					}
				}
			case err := <-c.watch.Errors:
				{
					fmt.Println("error: ", err)
					return
				}
			}
		}
	}()
}

func (c *Cache) watchFS(dir string) error {

	if !utils.Exists(dir) {
		return errors.New("WatchFS dir is error: not exist")
	}
	c.watchDir(dir)
	return nil

}

// /aaa/bbb/ccc/xxx.txt -> xxx
func (c *Cache) parseFileName(path_file string) (string, bool) {
	filenameWithSuffix := path.Base(path_file)
	fileSuffix := path.Ext(filenameWithSuffix)
	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	if "" != filenameOnly {
		return filenameOnly, true
	} else {
		return filenameOnly, false
	}
}

// /aaa/bbb/ccc/xxx.json -> xxx
const JSON_FILESUFFIX = ".json"

func (c *Cache) parseJsonFileName(path_file string) (string, bool) {

	//fmt.Println("fileSuffix =", fileSuffix)
	filenameWithSuffix := filepath.Base(path_file)
	fileSuffix := filepath.Ext(filenameWithSuffix)
	//fileSuffix = strings.TrimPrefix(fileSuffix,".")
	var filenameOnly string
	if JSON_FILESUFFIX == fileSuffix {
		filenameOnly = strings.TrimSuffix(filenameWithSuffix, JSON_FILESUFFIX)
		filenameOnly = strings.TrimSpace(filenameOnly)
		return filenameOnly, true
	} else {
		return "", false
	}

}

func (c *Cache) syncCacheServer(path string, op CacheOp) {
	fmt.Println(" start syncCacheServer...")
	// 添加要缓存文件
	filename, isJSON := c.parseJsonFileName(path)
	// 避免多次对相同的file重复注册api
	if op == CREATE && isJSON && !c.data.Set().Contains(filename) {
		// 是json文件，缓存处理
		c.data.Add(filename)
		fmt.Println("update files[CREATE] : ", filename)
	} else if op == DELETE {
		c.data.Remove(filename)
		fmt.Println("update files[DELETE] : ", filename)
	}
}
