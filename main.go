package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/2captcha/2captcha-go"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	addr = flag.String("addr", "rblxwild.com", "http service address")
	OldBalance = 0
)
type Config struct {
	AuthKey       string `json:"AuthKey"`
	TwoCaptchaKey string `json:"2CaptchaKey"`
	LatestRainID  int    `json:"LatestRainID"`
	Webhook       string `json:"Webhook"`
	Debug         bool   `json:"Debug"`
}
//wss://rblxwild.com/socket.io/?EIO=3&transport=websocket
//wss://rblxwild.com/socket.io/?EIO=3&transport=websocket
func Authentication(c *websocket.Conn) {
	var items = []string{"40", `42["chat:subscribe",{"channel":"EN"}]`, `42["cases:subscribe"]`, `42["events:subscribe"]`, "42[\"authentication\",null]"}
	color.Info.Tips("Authenticating...")

	for _,v := range items {
		err := c.WriteMessage(websocket.TextMessage, []byte(v))
		if err != nil {
			log.Fatal(err)
		}
		c.ReadMessage()
		currentPercent+=0.25
	}

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if err := tea.NewProgram(&m).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
func main() {

	timestart := time.Now()
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	jsonFile, err := os.Open("settings.json")
	if err != nil {
		os.Create("./settings.json")
		jsonFile, err = os.Open("settings.json")
		if err != nil {
			log.Fatal(err)
		}
		jsonFile.Write([]byte(`{"AuthKey":"","2CaptchaKey":"", "LatestRainID":0}`))
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var con Config
	err = json.Unmarshal(byteValue, &con)
	if con.Debug {
		color.Debug.Tips("Debug mode enabled")
		color.Debug.Tips(string(byteValue))
	}
	if err != nil {
		color.Error.Tips("Error restarting..: %s", err)
		time.Sleep(time.Second * 5)
		os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	}
	color.Notice.Tips("Configuration file loaded successfully.")
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/socket.io/", RawQuery: "EIO=4&transport=websocket"}
	if con.Debug {
		color.Debug.Tips("Connecting to %s", u.String())
	}
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)

	Authentication(c)

	if err != nil {
		log.Fatal("dial:", err, resp.Status)
	}

	defer c.Close()
	jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Bot online","color":2303786,"fields":[{"name":"Time Elapsed to Authenticate","value":"`+time.Since(timestart).String()+`"}],"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://pa1.narvii.com/6754/65dae7219636e753ca956d20cb510fb60deb0119_hq.gif"}}],"attachments":[]}`)
	req, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		color.Error.Tips("Error sending webhook: %s", err)
	}
	if con.Debug {
		color.Debug.Tips("Webhook sent. "+req.Status)
	}
	done := make(chan struct{})
	defer close(done)
	var oldPrice int
	for {
		c.WriteMessage(websocket.TextMessage, []byte("3"))
		if con.Debug{
			color.Debug.Tips("Sent ping to "+u.String())
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error while running; Restarting Info:", err)
			os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})

		}
		if strings.Contains(string(message), "updatePotVariables") {
			type Data struct {
				NewPrize              int `json:"newPrize"`
				NewJoinedPlayersCount int `json:"newJoinedPlayersCount"`
			}

			var data Data
			jsonData := strings.Split(string(message), `",`)[1];
			jsonData = strings.Replace(jsonData, "]", "", -1)
			err := json.Unmarshal([]byte(jsonData), &data)

			if err != nil {
				log.Fatalf("error: %v", err)
			}
			if con.Debug {
				color.Debug.Tips("Data structure complete")
				color.Debug.Tips("Pointer variable created for the data structure")
				color.Debug.Tips("NewPrize: %d", data.NewPrize)
				color.Debug.Tips("NewJoinedPlayersCount: %d", data.NewJoinedPlayersCount)
				color.Debug.Tips("Json unmarshalled")
				color.Debug.Tips("Updated old price variable with the new variable")
			}


			color.Info.Tips("Price has increased! " + strconv.Itoa(oldPrice) + " -> " + strconv.Itoa(data.NewPrize) + " | RainID: " + strconv.Itoa(con.LatestRainID))

			oldPrice = data.NewPrize


		} else if strings.Contains(string(message), "ENDING") && strings.Contains(string(message), "newState") {
			go func() {
				timer := time.Now()
				if con.Debug {
					color.Debug.Tips("Received ENDING message")
				}
				type Data struct {
					CaptchaToken string `json:"captchaToken"`
					PotID        int    `json:"potId"`
					Iloveu       bool   `json:"iloveu"`
				}
				if con.Debug {
					color.Debug.Tips("Pointer variable created for the data structure")
				}
				c := api2captcha.NewClient(con.TwoCaptchaKey)
				if con.Debug {
					color.Debug.Tips("Created new client")
				}
				c.DefaultTimeout = 120
				if con.Debug {
					color.Debug.Tips("Set default timeout to 120 seconds")
				}
				cap := api2captcha.HCaptcha{
					SiteKey: "30a8dcf5-481e-40d1-88de-51ad22aa8e97",
					Url:     "https://2captcha.com/demo/hcaptcha",
				}
				if con.Debug {
					color.Debug.Tips("Created new captcha")
				}
				code, err := c.Solve(cap.ToRequest())
				if err != nil {
					log.Println(err);
					if strings.Contains(err.Error(), "impossible") || strings.Contains(err.Error(), "timeout") {
						con.LatestRainID +=1
						if con.Debug {
							color.Debug.Tips("Captcha timeout TIME ELAPSED: "+strconv.Itoa(int(time.Since(timer))))
							jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Rain event","color":2303786,"fields":[{"name":"Success","value":"false","inline":true},{"name":"Time Elapsed","value":"`+time.Since(timer).String()+`","inline":true},{"name":"Error","value":"`+err.Error()+`","inline":true},{"name":"Rain Amount","value":"`+strconv.Itoa(oldPrice)+`"}],"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://pa1.narvii.com/6754/65dae7219636e753ca956d20cb510fb60deb0119_hq.gif"}}],"attachments":[]}`)
							http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
							if con.Debug {
								if data, err := ioutil.ReadAll(req.Body); err != nil {
									log.Println("Error:", err)
								} else {
									color.Debug.Tips("Webhook sent. "+req.Status+"\n"+string(data))
								}
							}
						}
					}
					return
				}

				var Dating Data
				if con.Debug {
					color.Debug.Tips("Created new data structure")
				}
				Dating.Iloveu = true
				Dating.CaptchaToken = code
				Dating.PotID = con.LatestRainID
				marshal, err := json.Marshal(Dating)
				if err != nil {
					return
				}
				req, _ := http.NewRequest("POST", "https://rblxwild.com/api/events/rain/join", strings.NewReader(string(marshal)))
				if con.Debug {
					color.Debug.Tips("Created new request")
				}
				req.Header.Set("Content-Type", "application/json")
				if con.Debug {
					color.Debug.Tips("Set content type to application/json")
				}
				req.Header.Set("X-Requested-With", "XMLHttpRequest")
				if con.Debug {
					color.Debug.Tips("Set X-Requested-With to XMLHttpRequest")
				}
				req.Header.Set("authorization", con.AuthKey)
				if con.Debug {
					color.Debug.Tips("Set authorization to " + con.AuthKey)
				}
				resp, _ := http.DefaultClient.Do(req)
				if con.Debug {
					color.Debug.Tips("Sent request")
				}
				if data, err := ioutil.ReadAll(resp.Body); err != nil {
					log.Println("Error:", err)
				} else {
					if strings.Contains(string(data), "true") {
						if con.Debug {
							color.Debug.Tips("Joined the rain successfully")
						}
						color.Success.Tips("$uccessfully joined the rain!")

						jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Rain event","color":2303786,"fields":[{"name":"Success","value":"true","inline":true},{"name":"Time Elapsed","value":"`+time.Since(timer).String()+`","inline":true},{"name":"Error","value":"nil","inline":true},{"name":"Rain Amount","value":"`+strconv.Itoa(oldPrice)+`"}],"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://pa1.narvii.com/6754/65dae7219636e753ca956d20cb510fb60deb0119_hq.gif"}}],"attachments":[]}`)
						req, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
						if err != nil {
							color.Error.Tips("Error sending webhook: "+err.Error())
						}
						if con.Debug {
							color.Debug.Tips("Webhook sent. "+req.Status)
						}
					} else {
						type DataBy struct {
							Success bool   `json:"success"`
							Message string `json:"message"`
							Elapsed int    `json:"elapsed"`
						}
						var DataByData DataBy
						json.Unmarshal(data, &DataByData)

						jsonData1 := []byte(`{"content":"@everyone","embeds":[{"title":"Rain event","color":2303786,"fields":[{"name":"Success","value":"false","inline":true},{"name":"Time Elapsed","value":"`+time.Since(timer).String()+`","inline":true},{"name":"Error","value":"`+DataByData.Message+`","inline":true},{"name":"Rain Amount","value":"`+strconv.Itoa(oldPrice)+`"}],"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://pa1.narvii.com/6754/65dae7219636e753ca956d20cb510fb60deb0119_hq.gif"}}],"attachments":[]}`)
						req1, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData1))
						log.Println(oldPrice, string(data), time.Since(timer).String())
						if err != nil {
							color.Error.Tips("Error sending webhook: "+err.Error())
							if data, err := ioutil.ReadAll(req1.Body); err != nil {
								log.Println("Error:", err)
							} else {
								color.Debug.Tips("Webhook sent. "+req1.Status+"\n"+string(data))

							}
						}





						color.Error.Tips("Failed to join the rain :( " + string(data))

					}
					con.LatestRainID += 1
				}
			}()
		}
	}
}


var (
	currentPercent = 0.00
)

const (
	padding  = 2
	maxWidth = 80
)



var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
type tickMsg time.Time

type model struct {
	progress progress.Model
}

func (_ *model) Init() tea.Cmd {
	return tickCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		cmd := m.progress.IncrPercent(currentPercent)
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (e *model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + e.progress.View() + "\n\n"


}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}