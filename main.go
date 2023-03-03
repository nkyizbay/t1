package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type UserInfo struct {
	X, Height string
}

var userid string

func main() {

	router := gin.Default()
	mrouter := melody.New()

	router.GET("/ws", func(c *gin.Context) {
		mrouter.HandleRequest(c.Writer, c.Request)
	})

	router.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "test.html")
	})

	mrouter.HandleConnect(func(s *melody.Session) {

		ss, _ := mrouter.Sessions() // array of sessions
		// fmt.Println(len(ss))        // 0

		if len(ss) == 0 {

			userid = "1"
			s.Set("info", &UserInfo{"56", "150px"}) // UserInfo{X, "350px"}

		} else {

			userid = "2"
			s.Set("info", &UserInfo{"", "150px"}) // UserInfo{X, "350px"}
			o := ss[0]
			info := o.MustGet("info").(*UserInfo) // there is key info in map; the value of info key is empty interface
			s.Write([]byte("setX " + info.X))     // set X 150px
			s.Write([]byte("setH " + info.X))     // set X 150px

		}

		s.Write([]byte("iam " + userid)) // writes message to session: "iam 1"

	})

	mrouter.HandleDisconnect(func(s *melody.Session) {
		fmt.Println("Session disconnected.")
	})

	mrouter.HandleMessage(func(s *melody.Session, msg []byte) { // "e.pageX 350px"
		p := strings.Split(string(msg), " ") // []string{"e.pageX", "350px"}
		if len(p) == 2 {
			info := s.MustGet("info").(*UserInfo) // UserInfo{"56", "150px"}

			info.X = p[0]
			info.Height = p[1]
			mrouter.BroadcastOthers([]byte("setX "+info.X), s)      // text to all sessions except s
			mrouter.BroadcastOthers([]byte("setH "+info.Height), s) // text to all sessions except s
		}
	})

	router.Run(":8080")
}
