package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Jacobious52/staticmbot/pkg/graph"
	"github.com/Jacobious52/staticmbot/pkg/model"
	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var sayings = []string{
	"(\\ /)\n(O.o)\n(> <) Bunny approves these changes.",
	"Yeah just like that!",
	"(c) Microsoft 1988",
	"-m 'So I hear you like commits ...'",
	"640K ought to be enough for anybody",
	"ALL SORTS OF THINGS",
	"I am sorry",
	"I'm human",
	"TODO: write meaningful commit message",
	"This is supposed to crash",
	"yolo push",
}

var tbToken = kingpin.Flag("token", "telegram bot token").Envar("TBTOKEN").Required().String()
var dbPath = kingpin.Flag("db", "path to save and read db file").Default("/usr/share/db/db.json").String()

func saveData(stats map[int64]model.Timeseries) {
	file, err := os.Create(*dbPath)
	if err != nil {
		log.Errorf("data not saved. %v", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&stats)
	if err != nil {
		log.Errorf("failed to encode json: %v", err)
		return
	}
}

func randSaying() string {
	return sayings[rand.Intn(len(sayings))]
}
func loadData() map[int64]model.Timeseries {

	stats := make(map[int64]model.Timeseries)

	file, err := os.Open(*dbPath)
	if err != nil {
		log.Warningf("data not opened. %v", err)
		return stats
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&stats)
	if err != nil {
		log.Errorf("json not parsed. %v", err)
		return stats
	}
	return stats
}

func renderGraph(grapher graph.Grapher) (*tb.Photo, error) {
	log.Debugln("Rendering png buffer")

	fname := fmt.Sprint(time.Now(), ".png")
	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	err = grapher.Render(f)
	if err != nil {
		return nil, err
	}
	f.Close()

	image := &tb.Photo{File: tb.FromDisk(fname)}
	return image, nil
}

func main() {
	kingpin.Parse()

	bot, err := tb.NewBot(tb.Settings{
		Token: *tbToken,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})

	if err != nil {
		log.Fatalf("failed to created bot: %v", err)
	}
	log.Infoln("Created bot")

	stats := loadData()

	var lock sync.RWMutex
	var changes bool
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			if changes {
				lock.RLock()
				saveData(stats)
				lock.RUnlock()
				log.Infoln("Saved data..")
				changes = false
			}
		}
	}()

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender, fmt.Sprintf("Hello, %v. Welcome to the bot!", m.Sender.FirstName))
	})

	bot.Handle("/fake", func(m *tb.Message) {
		if !m.FromGroup() {
			log.Warningln("not in group chat")
			return
		}

		graphMaker := graph.Default["alphav1scatter"]
		words := strings.Split(m.Text, " ")
		log.Debug(words)
		if len(words) == 2 {
			if g, ok := graph.Default[words[1]]; ok {
				graphMaker = g
			}
		}

		log.Infoln("generating fake stats for", m.Chat.Title)
		bot.Send(m.Chat, "Generating fake graph")
		// fakes for testing
		fakeGrapher := graphMaker(model.FakeSeries)
		var fakePhoto, err = renderGraph(fakeGrapher)
		if err != nil {
			log.Errorf("render failed: %v", err)
			return
		}
		bot.Send(m.Chat, fakePhoto)
	})

	bot.Handle("/speak", func(m *tb.Message) {
		bot.Send(m.Chat, randSaying())
	})

	bot.Handle("/stats", func(m *tb.Message) {
		if !m.FromGroup() {
			log.Warningln("not in group chat")
			return
		}

		graphMaker := graph.Default["alphav1scatter"]
		words := strings.Split(m.Text, " ")
		if len(words) == 2 {
			if g, ok := graph.Default[words[1]]; ok {
				graphMaker = g
			}
		}

		log.Infoln("generating stats for ", m.Chat.Title)
		bot.Send(m.Chat, randSaying())

		lock.RLock()
		series, ok := stats[m.Chat.ID]
		if !ok {
			log.Warningln("no series exists for this chat")
			return
		}

		grapher := graphMaker(series)
		image, err := renderGraph(grapher)
		lock.RUnlock()

		if err != nil {
			log.Errorf("render failed: %v", err)
			bot.Send(m.Chat, "Sorry, something went wrong")
			go func() {
				time.Sleep(3 * time.Second)
				bot.Send(m.Chat, "I need to be punished ;)")
			}()
			return
		}
		bot.Send(m.Chat, image)
	})

	bot.Handle("/save", func(m *tb.Message) {
		lock.RLock()
		saveData(stats)
		lock.RUnlock()
	})

	bot.Handle("/load", func(m *tb.Message) {
		lock.Lock()
		stats = loadData()
		lock.Unlock()
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		user := m.Sender.Username
		if len(user) == 0 {
			user = fmt.Sprintf("%v %v", m.Sender.FirstName, m.Sender.LastName)
		}

		log.Infof("%v[%v] %v", m.Chat.Title, user, m.Text)

		lock.RLock()
		_, ok := stats[m.Chat.ID]
		if !ok {
			stats[m.Chat.ID] = make(model.Timeseries)
		}
		series, _ := stats[m.Chat.ID]
		lock.RUnlock()

		day := m.Time().YearDay()
		if _, ok := series[day]; !ok {
			lock.Lock()
			series[day] = map[string]uint64{
				user: 1,
			}
			lock.Unlock()
		} else {
			lock.Lock()
			series[day][user]++
			lock.Unlock()
		}
		lock.Lock()
		stats[m.Chat.ID] = series
		lock.Unlock()
		changes = true
	})

	bot.Start()
}
