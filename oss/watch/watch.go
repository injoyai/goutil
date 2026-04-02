package watch

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/injoyai/conv"
)

type (
	Op    = fsnotify.Op
	Event = fsnotify.Event
)

var (
	Create = fsnotify.Create
	Write  = fsnotify.Write
	Remove = fsnotify.Remove
	Rename = fsnotify.Rename
	Chmod  = fsnotify.Chmod
)

func Watch[T string | []string](path T, fn Handler) error {
	return New(path, fn).Run()
}

type Handler func(e Event)

type Watcher struct {
	dirs map[string]*files

	debounce *Debounce[string]
	cancel   context.CancelFunc

	callback Handler
}

func New[T string | []string](path T, fn Handler) *Watcher {

	dirs := map[string]*files{}

	for _, originName := range conv.Strings(path) {

		//统一格式,绝对路径
		fullName, _ := filepath.Abs(originName)

		//如果监听的是文件夹,则不再监听该目录下的文件
		if len(originName) > 0 && originName[len(originName)-1] == '/' {
			dirs[fullName] = &files{
				isDir:  true,
				origin: originName,
			}
			continue
		}

		dir := filepath.Dir(fullName)
		//如果该文件所在目录被监听,则不再监听该文件
		if val := dirs[dir]; val != nil && val.isDir {
			continue
		}

		if _, ok := dirs[dir]; !ok {
			dirs[dir] = &files{files: map[string]string{}}
		}

		//储存输入的文件原始路径
		dirs[dir].files[fullName] = originName
	}

	return &Watcher{
		dirs:     dirs,
		debounce: NewDebounce[string](300 * time.Millisecond),
		callback: fn,
	}
}

func (w *Watcher) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}

func (w *Watcher) SetAfter(d time.Duration) {
	w.debounce.SetAfter(d)
}

func (w *Watcher) Run(ctx ...context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for dir, _ := range w.dirs {
		if err := watcher.Add(dir); err != nil {
			return err
		}
	}

	_ctx := conv.Default(context.Background(), ctx...)
	_ctx, w.cancel = context.WithCancel(_ctx)
	defer w.cancel()

	for {
		select {
		case <-_ctx.Done():
			return _ctx.Err()

		case ev, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			//统一格式,绝对路径
			fullName, _ := filepath.Abs(ev.Name)

			//如果监听到文件夹创建,切在监听目录下,则需要继续监听
			if f := w.dirs[fullName]; f != nil && ev.Op == Create {
				watcher.Add(f.origin)
				continue
			}

			dir := filepath.Dir(fullName)
			f := w.dirs[dir]
			if f == nil {
				continue
			}

			if f.isDir {
				ev.Name = f.origin + filepath.Base(ev.Name)
			} else {
				filename := filepath.Clean(fullName)
				ev.Name, ok = f.files[filename]
				if !ok {
					continue
				}
			}

			w.debounce.Do(ev.Op.String()+ev.Name, func() { w.callback(ev) })

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			log.Println("watch error:", err)

		}
	}

}

type files struct {
	isDir  bool
	origin string
	files  map[string]string
}

func NewDebounce[K comparable](after time.Duration) *Debounce[K] {
	return &Debounce[K]{
		after:  after,
		timers: map[K]*time.Timer{},
	}
}

type Debounce[K comparable] struct {
	after  time.Duration
	timers map[K]*time.Timer
	mu     sync.Mutex
}

func (d *Debounce[K]) SetAfter(after time.Duration) {
	d.after = after
}

func (d *Debounce[K]) Do(k K, f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t := d.timers[k]; t != nil {
		t.Stop()
	}

	d.timers[k] = time.AfterFunc(d.after, func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("panic:", r)
			}
		}()
		f()
	})
}
