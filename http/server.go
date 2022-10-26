package http

import (
	"embed"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"reader/data"
	"reader/http/assets"
)

const templatesDir = "templates"

var (
	//go:embed templates/* templates/layouts/*
	files     embed.FS
	templates map[string]*template.Template
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router
	sc     *securecookie.SecureCookie

	// Bind address
	Addr string

	// Keys used for secure cookie encryption.
	HashKey  string
	BlockKey string

	// Data services
	Feeds       data.Feeds
	UnreadItems data.UnreadItems
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}

	if err := LoadTemplates(); err != nil {
		panic(fmt.Sprintf("Failed to load templates: %s", err))
	}

	s.server.Handler = http.HandlerFunc(s.router.ServeHTTP)

	s.router.Use(loggingMiddleware)

	s.router.PathPrefix("/stylesheets/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/javascripts/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/dist/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/fonts/").
		Handler(http.FileServer(http.FS(assets.FS)))
	s.router.PathPrefix("/img/").
		Handler(http.FileServer(http.FS(assets.FS)))

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["index.tmpl"]
		if !ok {
			log.Printf("Failed to load template index.tmpl")
			w.WriteHeader(500)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Printf("Failed executing template index.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	}).Methods("GET")

	s.router.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["news.tmpl"]
		if !ok {
			log.Printf("Failed to load template news.tmpl")
			w.WriteHeader(500)
			return
		}

		viewData := struct {
			UnreadStories []string
			StoriesJson   template.JS
		}{
			[]string{"hi"},
			`[{"id":229542,"title":"Nvidia investigating reports of RTX 4090 power cables burning or melting","permalink":"https://www.theverge.com/2022/10/25/23422349/nvidia-rtx-4090-power-cables-connectors-melting-burning","body":"  \n\n    <figure>\n      <img alt=\"nvidia, nvidiastock\" src=\"https://cdn.vox-cdn.com/thumbor/3p89qY0UnqAYOPzyhkZ9bke-WrY=/0x0:2040x1360/1310x873/cdn.vox-cdn.com/uploads/chorus_image/image/71539868/DSC00697.0.jpg\">\n        <figcaption>Photo by Sam Byford / The Verge</figcaption>\n    </figure>\n\n  <p id=\"DqzyJe\">Nvidia says it’s investigating two reports of RTX 4090 cards that have had power cables burn or melt. Reddit user reggie_gakil was the <a href=\"https://www.reddit.com/r/nvidia/comments/yc6g3u/rtx_4090_adapter_burned/\">first to post</a> details about their Gigabyte RTX 4090 issues yesterday, showing burn damage on the new 12VHPWR adapter cable that Nvidia ships with the RTX 4090. The connection on the actual card was also damaged and melted.</p>\n<p id=\"lI9CsK\">A second <a href=\"https://www.reddit.com/r/nvidia/comments/yc6g3u/comment/itmjt9d/?utm_source=reddit&amp;utm_medium=web2x&amp;context=3\">Reddit poster</a> replied to the same thread showing off similar damage to an Asus RTX 4090 graphics card power connector. </p>\n<p id=\"5jOgfg\">“We are investigating the reports,” says Nvidia spokesperson Bryan Del Rizzo in a statement to <em>The Verge</em>. “We are in contact with the first owner and will be reaching out to the other for additional information.”</p>\n  <figure class=\"e-image\">\n        \n      <cite>Image: reggie_gakil (<a class=\"ql-link\" href=\"https://www.reddit.com/r/nvidia/comments/yc6g3u/rtx_4090_adapter_burned/\" target=\"_blank\">Reddit</a>)</cite>\n...</figure>\n  <p>\n    <a href=\"https://www.theverge.com/2022/10/25/23422349/nvidia-rtx-4090-power-cables-connectors-melting-burning\">Continue reading…</a>\n  </p>\n\n","feed_id":20,"created_at":"2022-10-25T09:22:23.661Z","updated_at":"2022-10-25T09:22:23.661Z","published":"2022-10-25T09:14:15.000Z","is_read":false,"keep_unread":false,"is_starred":false,"entry_id":"https://www.theverge.com/2022/10/25/23422349/nvidia-rtx-4090-power-cables-connectors-melting-burning","headline":"Nvidia investigating reports of RTX 4090 power cab","lead":"  \n\n    \n      \n        Photo by Sam Byford / The Verge\n    \n\n  Nvidia says it’s investigating two r","source":"The Verge -  All Posts","pretty_date":"Oct 25, 09:14"},{"id":229543,"title":"Flight Simulator Focuses on the Other Side of the Cockpit Door","permalink":"https://hackaday.com/2022/10/25/flight-simulator-focuses-on-the-other-side-of-the-cockpit-door/","body":"<div><img width=\"800\" height=\"450\" src=\"https://hackaday.com/wp-content/uploads/2022/10/1920x1080-vtime1_54-take2022-10-24-15.35.26.png?w=800\" class=\"attachment-large size-large wp-post-image\" style=\"margin:0 auto;margin-bottom:15px;\" data-attachment-id=\"559615\" data-permalink=\"https://hackaday.com/2022/10/25/flight-simulator-focuses-on-the-other-side-of-the-cockpit-door/1920x1080-vtime1_54-take2022-10-24-15-35-26/\" data-orig-file=\"https://hackaday.com/wp-content/uploads/2022/10/1920x1080-vtime1_54-take2022-10-24-15.35.26.png\" data-orig-size=\"1920,1080\" data-comments-opened=\"1\" data-image-meta='{\"aperture\":\"0\",\"credit\":\"\",\"camera\":\"\",\"caption\":\"\",\"created_timestamp\":\"0\",\"copyright\":\"\",\"focal_length\":\"0\",\"iso\":\"0\",\"shutter_speed\":\"0\",\"title\":\"\",\"orientation\":\"0\"}' data-image-title=\"[1920×1080] vtime=[1_54], take=[2022-10-24 15.35.26]\" data-medium-file=\"https://hackaday.com/wp-content/uploads/2022/10/1920x1080-vtime1_54-take2022-10-24-15.35.26.png?w=400\" data-large-file=\"https://hackaday.com/wp-content/uploads/2022/10/1920x1080-vtime1_54-take2022-10-24-15.35.26.png?w=800\"></div><p>When one thinks of getting into a flight simulator, one assumes that it’ll be from the pilot’s point of view. But <a href=\"https://alexshakespeare.com/2022/10/alternative-flight-simulator/\" target=\"_blank\">this alternative flight simulator</a> takes a different tack, by letting you live out your air travel fantasies from the passenger’s point of view.</p>\n<p>Those of you looking for a full-motion simulation of the passenger cabin experience will be disappointed, as [Alex Shakespeare] — we assume no relation — has built a minimal airliner cabin for this simulator. That makes sense, though; ideally, an airline pilot aims to provide passengers with as dull a ride as possible. Where a flight is at its most exciting, and what [Alex] captures nicely here, is the final approach to your destination, when the airport and its surrounding environs finally come into view after a long time staring at clouds. This is done by mounting an LCD monitor outside the window of a reasonable facsimile of an airliner cabin, complete with a row of seats. A control panel that originally lived in an airliner cockpit serves to select video of approaches to airports in various exotic destinations, like Las Vegas. The video is played by a Pi Zero, while an ESP32 takes care of controlling the lights, fans, and attendant call buttons in the quite realistic-looking overhead panel. Extra points for the button that plays the Ryanair arrival jingle.</p>\n<p>[Alex]’s simulator is impressively complete, if somewhat puzzling in conception. We don’t judge, though, and it looks like it might be fun for visitors, especially when <a href=\"https://hackaday.com/2021/07/28/workshop-tools-are-available-in-first-class/\">the drinks cart</a> comes by.</p>\n<p><span id=\"more-559580\"></span></p>\n<p></p>\n","feed_id":16,"created_at":"2022-10-25T09:52:24.677Z","updated_at":"2022-10-25T09:52:24.677Z","published":"2022-10-25T08:00:18.000Z","is_read":false,"keep_unread":false,"is_starred":false,"entry_id":"https://hackaday.com/?p=559580","headline":"Flight Simulator Focuses on the Other Side of the ","lead":"When one thinks of getting into a flight simulator, one assumes that it’ll be from the pilot’s point","source":"Hack a Day","pretty_date":"Oct 25, 08:00"},{"id":229541,"title":"WhatsApp is down in a major outage","permalink":"https://www.theverge.com/2022/10/25/23422343/whatsapp-down-outage","body":"  \n\n    <figure>\n      <img alt=\"Illustration of a number of green WhatsApp logos in black circles floating across a blue background\" src=\"https://cdn.vox-cdn.com/thumbor/J4pd_ATrmv4dY5s_Ug0V2jWVWFQ=/0x0:2040x1360/1310x873/cdn.vox-cdn.com/uploads/chorus_image/image/71539795/acastro_210119_1777_whatsapp_0003.0.jpg\">\n        <figcaption>Illustration by Alex Castro / The Verge</figcaption>\n    </figure>\n\n  <p id=\"qM36sn\">WhatsApp is down in a major outage affecting the communications app. The service started experiencing issues at around 3AM ET, with users greeted with a “connecting” message. If you’re attempting to use WhatsApp web, you’ll see a “Make sure your computer has an active internet connection” error. </p>\n<p id=\"aDnBgu\"><a href=\"https://downdetector.co.uk/status/whatsapp/\">Downdetector</a> has more than 60,000 reports of issues with the service, but WhatsApp parent company Meta hasn’t issued a comment on the problems yet. The outage appears to be affecting users globally.</p>\n<p id=\"GMdBIZ\">This is the first major WhatsApp outage since the service went down as part of a massive outage that also took down Instagram, Messenger, Oculus, and Facebook last year. That outage took nearly six hours before it was resolved and WhatsApp was back...</p>\n  <p>\n    <a href=\"https://www.theverge.com/2022/10/25/23422343/whatsapp-down-outage\">Continue reading…</a>\n  </p>\n\n","feed_id":20,"created_at":"2022-10-25T08:02:41.285Z","updated_at":"2022-10-25T08:02:41.285Z","published":"2022-10-25T07:49:35.000Z","is_read":false,"keep_unread":false,"is_starred":false,"entry_id":"https://www.theverge.com/2022/10/25/23422343/whatsapp-down-outage","headline":"WhatsApp is down in a major outage","lead":"  \n\n    \n      \n        Illustration by Alex Castro / The Verge\n    \n\n  WhatsApp is down in a major ","source":"The Verge -  All Posts","pretty_date":"Oct 25, 07:49"},{"id":229540,"title":"Luxury surveillance","permalink":"https://flowingdata.com/2022/10/25/luxury-surveillance/","body":"<p>Chris Gilliard, for The Atlantic, describes self-surveillance that people pay for <a href=\"https://www.theatlantic.com/technology/archive/2022/10/amazon-tracking-devices-surveillance-state/671772/\">in exchange for small conveniences at the expense of privacy</a>: </p>\n<blockquote><p>The conveniences promised by Amazon’s suite of products may seem divorced from this context: I am here to tell you that they’re not. These “smart” devices all fall under the umbrella of what the digital-studies scholar David Golumbia and I call “<a href=\"https://reallifemag.com/luxury-surveillance/\">luxury surveillance</a>“—that is, surveillance that people pay for and whose tracking, monitoring, and quantification features are understood by the user as benefits. These gadgets are analogous to the surveillance technologies deployed in Detroit and many other cities across the country in that they are best understood as mechanisms of control: They gather data, which are then used to affect behavior. Stripped of their gloss, these devices are similar to the ankle monitors and surveillance apps such as SmartLINK that are forced on people on parole or immigrants awaiting hearings. As the author and activist James Kilgore writes, “The ankle monitor—which for almost two decades was simply an analog device that informed authorities if the wearer was at home—has now grown into a sophisticated surveillance tool via the use of GPS capacity, biometric measurements, cameras, and audio recording.”</p></blockquote>\n<p><strong>Tags:</strong> <a href=\"https://flowingdata.com/tag/amazon/\" rel=\"tag\">Amazon</a>, <a href=\"https://flowingdata.com/tag/atlantic/\" rel=\"tag\">Atlantic</a>, <a href=\"https://flowingdata.com/tag/chris-gilliard/\" rel=\"tag\">Chris Gilliard</a>, <a href=\"https://flowingdata.com/tag/privacy/\" rel=\"tag\">privacy</a></p>","feed_id":6,"created_at":"2022-10-25T07:12:30.214Z","updated_at":"2022-10-25T07:12:30.214Z","published":"2022-10-25T07:08:04.000Z","is_read":false,"keep_unread":false,"is_starred":false,"entry_id":"https://flowingdata.com/?p=69439","headline":"Luxury surveillance","lead":"Chris Gilliard, for The Atlantic, describes self-surveillance that people pay for in exchange for sm","source":"FlowingData","pretty_date":"Oct 25, 07:08"},{"id":229539,"title":"One Of The Worst Keyboards Ever, Now An Arduino Peripheral","permalink":"https://hackaday.com/2022/10/24/one-of-the-worst-keyboards-ever-now-an-arduino-peripheral/","body":"<div><img width=\"800\" height=\"450\" src=\"https://hackaday.com/wp-content/uploads/2022/10/zxkeys-arduino-featured.jpg?w=800\" class=\"attachment-large size-large wp-post-image\" style=\"margin:0 auto;margin-bottom:15px;\" data-attachment-id=\"559589\" data-permalink=\"https://hackaday.com/2022/10/24/one-of-the-worst-keyboards-ever-now-an-arduino-peripheral/zxkeys-arduino-featured/\" data-orig-file=\"https://hackaday.com/wp-content/uploads/2022/10/zxkeys-arduino-featured.jpg\" data-orig-size=\"800,450\" data-comments-opened=\"1\" data-image-meta='{\"aperture\":\"0\",\"credit\":\"\",\"camera\":\"\",\"caption\":\"\",\"created_timestamp\":\"0\",\"copyright\":\"\",\"focal_length\":\"0\",\"iso\":\"0\",\"shutter_speed\":\"0\",\"title\":\"\",\"orientation\":\"0\"}' data-image-title=\"zxkeys-arduino-featured\" data-medium-file=\"https://hackaday.com/wp-content/uploads/2022/10/zxkeys-arduino-featured.jpg?w=400\" data-large-file=\"https://hackaday.com/wp-content/uploads/2022/10/zxkeys-arduino-featured.jpg?w=800\"></div><p>For British kids of a certain age, their first experience of a computer was very likely to have been in front of a Sinclair ZX81. The lesser-known predecessor to the wildly-successful ZX Spectrum, it came in at under £100 and sported a Z80 processor and a whopping 1k of memory. In the long tradition of Sinclair products it had a few compromises to achieve that price point, the most obvious of which was a 40-key membrane keyboard. Those who learned to code on its frustrating lack of tactile feedback may be surprised to see <a href=\"https://create.arduino.cc/projecthub/sl001/read-a-zx81-keyboard-with-arduinos-and-build-things-with-it-0189bd\" target=\"_blank\">an Arduino project presenting it as the perfect way to easily hook up a keyboard to an Arduino</a>.</p>\n<p>Like many retrocomputing parts, the ZX81 ‘board has been re-manufactured, to the joy of many a Sinclair enthusiast. It’s thus readily available and relatively cheap (we think they can be found for less than the stated 20 euros!), so surprisingly it’s a reasonable choice for an Arduino project. The task of trying to define by touch the imperceptible difference in thickness of a ZX81 key will bring a true retrocomputing experience to a new generation. Perhaps <a href=\"https://hackaday.com/2014/12/22/zx81-emulated-on-an-mbed/\">if it can be done on an Mbed</a> then someone might even make a ZX81 emulator on the Arduino.</p>\n<p><a href=\"https://hackaday.com/2020/04/01/accurate-dispensing-of-toilet-paper-will-get-us-through-the-crisis/\">We’re great fans of the ZX81 here at Hackaday</a>, for some of us it was that first computer. Long may it continue to delight its fans!</p>\n","feed_id":16,"created_at":"2022-10-25T05:22:25.800Z","updated_at":"2022-10-25T05:22:25.800Z","published":"2022-10-25T05:00:20.000Z","is_read":false,"keep_unread":false,"is_starred":false,"entry_id":"https://hackaday.com/?p=559575","headline":"One Of The Worst Keyboards Ever, Now An Arduino Pe","lead":"For British kids of a certain age, their first experience of a computer was very likely to have been","source":"Hack a Day","pretty_date":"Oct 25, 05:00"}]`,
		}

		//if storiesJson, err := json.Marshal(viewData.UnreadStories); err != nil {
		//	log.Printf("Error building stories JSON: %s", err)
		//	w.WriteHeader(500)
		//	return
		//} else {
		//	viewData.StoriesJson = template.JS(storiesJson)
		//}

		if err := t.Execute(w, viewData); err != nil {
			log.Printf("Failed executing template news.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	}).Methods("GET")

	s.router.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["feeds.tmpl"]
		if !ok {
			log.Printf("Failed to load template feeds.tmpl")
			w.WriteHeader(500)
			return
		}

		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				log.Printf("POST /form: Unable to parse request body: %s", err)
				w.WriteHeader(500)
				return
			}

			// TODO going to have to figure out how to get a flash working w/ the cookie
			u, err := url.Parse(r.FormValue("feed_url"))
			if err != nil {
				log.Printf("POST /feed: error parsing feed url: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			_, err = s.Feeds.AddFeed("test", *u)
			if err != nil {
				log.Printf("POST /form: error adding feed to DB: %s", err)
				return
			}
		}

		feeds, err := s.Feeds.GetFeedList()
		if err != nil {
			log.Printf("Error getting feed list from DB: %s", err)
			w.WriteHeader(500)
			return
		}

		templateData := struct {
			Test  string
			Feeds []data.Feed
		}{
			Test:  "Feeds list!",
			Feeds: feeds,
		}

		if err := t.Execute(w, templateData); err != nil {
			log.Printf("Failed executing template feeds.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	}).Methods("GET", "POST")

	// TODO handle POST /feeds

	s.router.HandleFunc("/feeds/new", func(w http.ResponseWriter, r *http.Request) {
		t, ok := templates["add_feed.tmpl"]
		if !ok {
			log.Printf("Failed to load template add_feed.tmpl")
			w.WriteHeader(500)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Printf("Failed executing template add_feed.tmpl: %s", err)
			w.WriteHeader(500)
			return
		}
	})

	return s
}

func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	go func() {
		err := s.server.Serve(s.ln)
		if err != nil {
			log.Fatalf("Unable to start HTTP server: %s", err)
		}
	}()

	log.Printf("Server opened on %q", s.Addr)

	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), templatesDir+"/layouts/*", templatesDir+"/partials/*")
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}
