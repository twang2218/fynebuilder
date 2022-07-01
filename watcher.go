package fynebuilder

import (
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

func NewWatcher(file string, res ResourceDict, handler func(ObjectDict)) *fsnotify.Watcher {
	//	create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("NewWatcher failed: %v", err)
	}

	//	watch the given file
	watcher.Add(file)

	f := func() {
		objs, err := Load(file, res)
		if err != nil {
			log.Errorf("Load(%s) error: %v", file, err)
		}
		//	update ObjectDict
		handler(objs)
	}

	//	update ObjectDict if file changed
	go func() {
		//	initial load the UI
		f()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				t := time.Now()
				// log.Tracef("%s %s", event.Name, event.Op)
				if event.Op == fsnotify.Write {
					f()
					//	show reloaded time
					log.Debugf("Reloaded %q in %v.", file, time.Since(t))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("watcher.Errors: %v", err)
			}
		}
	}()

	return watcher
}
