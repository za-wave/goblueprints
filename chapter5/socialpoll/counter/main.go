package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var fatalErr error

func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

const updateDuration = 1 * time.Second

func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	log.Println("Connecting to DB...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("Closing DB...")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	var countsLock sync.Mutex
	var counts map[string]int

	log.Println("Conecting to NSQ...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		votes := string(m.Body)
		counts[votes]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}

	log.Println("Wating for votes on NSQ...")
	var updater *time.Timer
	updater = time.AfterFunc(updateDuration, func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("No new votes found. Skip to update DB.")
		} else {
			log.Println("Update DB.")
			log.Println("counts")
			ok := true
			for option, count := range counts {
				sel := bson.M{"options": bson.M{"$in": []string{option}}}
				up := bson.M{"$inc": bson.M{"results." + option: count}}
				if _, err := pollData.UpdateAll(sel, up); err != nil {
					log.Println("Faild to update DB:", err)
					ok = false
					continue
				}
				counts[option] = 0
			}
			if ok {
				log.Println("Success to update DB.")
				counts = nil // reset counts
			}
		}
		updater.Reset(updateDuration)
	})

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-termChan:
			updater.Stop()
			q.Stop()
		case <-q.StopChan:
			//Success!
			return
		}
	}
}
