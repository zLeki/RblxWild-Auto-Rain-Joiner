package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	addr                = flag.String("addr", "rblxwild.com", "http service address")
	oldPrice            int
	PotID               int
	RainsJoined         = 0
	MOTD                = []string{"discord.gg/KZpWWW6eK4", "spencer likes men", "this is a virus", "github.com/zLeki", "asm is sexy", "message of the day", "idk what to put here", "view my portfolio leki.sbs/portfolio", "yes", "noah is hot", "im rich (not in real life)"}
	AuthorizationConfig *ParsedData
)

type ParsedData struct {
	UserData struct {
		DisplayName string `json:"displayName"`
		Balance     int    `json:"balance"`
	} `json:"userData"`
	Events struct {
		Rain struct {
			Pot struct {
				ID                 int    `json:"id"`
				Prize              int    `json:"prize"`
				State              string `json:"state"`
				CreatedAt          int    `json:"createdAt"`
				LastUpdateMs       int64  `json:"lastUpdateMs"`
				JoinedPlayersCount int    `json:"joinedPlayersCount"`
			} `json:"pot"`
		} `json:"rain"`
	} `json:"events"`
}
type Config struct {
	AuthKey       string `json:"AuthKey"`
	TwoCaptchaKey string `json:"2CaptchaKey"`
	LatestRainID  int    `json:"LatestRainID"`
	Webhook       string `json:"Webhook"`
	Debug         bool   `json:"Debug"`
	SafeMode      bool   `json:"SafeMode"`
}

//wss://rblxwild.com/socket.io/?EIO=3&transport=websocket
//wss://rblxwild.com/socket.io/?EIO=3&transport=websocket
func Authentication(c *websocket.Conn) {

restart:
	jsonFile, err := os.Open("./settings.json")
	if err != nil {
		os.Create("./settings.json")
		jsonFile, err = os.Open("./settings.json")
		if err != nil {
			log.Println(err)
		}
		jsonData, _ := json.Marshal(Config{
			AuthKey:       "",
			TwoCaptchaKey: "",
			LatestRainID:  0,
			Webhook:       "",
			Debug:         false,
			SafeMode:      true,
		})
		ioutil.WriteFile("./settings.json", jsonData, 0644)
		fmt.Println(err)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		goto restart
	}
	var foo Config
	err = json.Unmarshal(byteValue, &foo)
	var items = []string{"40", `42["events:subscribe"]`, `42["authentication",{"authToken":"` + foo.AuthKey + `","clientTime":1651530049953}]`, `42["chat:subscribe",{"channel":"EN"}]`, `42["cases:subscribe"]`}
	color.Info.Tips("Authenticating...")

	for _, v := range items {
		time.Sleep(time.Second * 1)
		err5 := c.WriteMessage(websocket.TextMessage, []byte(v))
		if err5 != nil {
			log.Println("Error while reading; Restarting Info:", err5)
			time.Sleep(time.Second * 5)
			goto restart
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error while running; Restarting Info:", err)

			time.Sleep(time.Second * 5)
			goto restart

		}
		if strings.Contains(string(message), "authenticationResponse") {

			message = []byte(strings.Split(string(message), "42[\"authenticationResponse\",")[1])
			message = []byte(strings.Split(string(message), "]")[0])
			var parsedData ParsedData
			json.Unmarshal(message, &parsedData)
			jsonFile, err := os.Open("./settings.json")
			if err != nil {
				log.Println("Error while opening settings.json:", err)

				goto restart
			}
			byteValue, _ := ioutil.ReadAll(jsonFile)
			var con Config
			err = json.Unmarshal(byteValue, &con)
			con.LatestRainID = parsedData.Events.Rain.Pot.ID
			AuthorizationConfig = &parsedData
			saveJson, _ := json.Marshal(con)
			ioutil.WriteFile("./settings.json", saveJson, 0644)

		}
	}

}
func GetLogo() string {
	return `

⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡎⠀⠀⠀⠀⠀⠀⣴⡏⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣾⠟⠀⠀⠀⠀⠀⣠⣾⡿⣿⣿⣿⣷⣶⣶⣶⠀
⠀⣿⣿⣷⣶⣶⣶⣿⣿⣿⣿⣟⣁⣀⣀⣀⣠⣴⣾⡿⠋⠀⠀⠈⠉⠉⠉⠉⠉⠀
⠀⠈⠉⠉⠉⠉⠉⠉⠙⠛⠻⠿⠿⠿⠿⠛⠛⢋⣁⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⣶⣶⣶⣶⣾⣿⣿⣿⣷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⣿⣿⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀⠀⠀⠀⠀⠀		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀▒█▀▀▀█ █▀▀ █▀▀█ █▀▀█ █▀▀█ █▀▀ ▒█▀▀▀█ ▀▀█▀▀ █▀▀█ █▀▀█ █▀▄▀█
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢤⣶⣶⣾⣿⣿⣿⣿⣿⣿⣟⠉⠉⠀⠀⠀⠀⠀░▀▀▀▄▄ █░░ █▄▄▀ █▄▄█ █░░█ █▀▀ ░▀▀▀▄▄ ░░█░░ █░░█ █▄▄▀ █░▀░█
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠿⣿⣿⣿⣿⣿⣿⣿⣧⡀⠀⠀⠀⠀⠀▒█▄▄▄█ ▀▀▀ ▀░▀▀ ▀░░▀ █▀▀▀ ▀▀▀ ▒█▄▄▄█ ░░▀░░ ▀▀▀▀ ▀░▀▀ ▀░░░▀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠻⢿⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀		b    y        l    e    k    i
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⠿⣿⣿⣿⣄⠀⠀⠀	   [ Message of the day "` + MOTD[rand.Intn(len(MOTD)-1)] + `" ]
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⢿⣆⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`
}

func main() {
	fmt.Printf(GetLogo() + `

⠀⠀⠀		[    P r e s s    e n t e r    t o    c o n t i n u e    ]
`)

	_, err := fmt.Scanln()
	if err != nil {
		return
	}
	fmt.Println("Starting...")
	fmt.Println("Checking for updates...")
	time.Sleep(time.Second * 2)
	fmt.Println("Starting bot...")

	go GUI()

restart:

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/socket.io/", RawQuery: "EIO=4&transport=websocket"}

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)

	Authentication(c)

	if err != nil {
		log.Println("dial:", err, resp.Status)
	}

	defer c.Close()
	//jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Authentication success","color":2303786,"fields":[{"name":"Success","value":"true","inline":true},{"name":"Time Elapsed to authenticate","value":"` + time.Since(timestart).String() + `","inline":true},{"name":"Error","value":"nil","inline":true},{"name":"Rain Amount","value":"nil","inline":true},{"name":"PotID","value":"` + strconv.Itoa(con.LatestRainID) + `","inline":true},{"name":"Authentication","value":"` + con.AuthKey + `","inline":true},{"name":"SafeMode","value":"` + strconv.FormatBool(con.SafeMode) + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"image":{"url":"https://cdn.dribbble.com/users/1314513/screenshots/3928265/storm.gif"}}],"username":"StormScraper","avatar_url":"https://cdn.dribbble.com/users/1314513/screenshots/3928265/storm.gif","attachments":[]}`)
	//req, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
	//if err != nil {
	//	color.Error.Tips("Error sending webhook: %s", err)
	//}
	//if con.Debug {
	//	color.Debug.Tips("Webhook sent. " + req.Status)
	//}
	jsonFile, err := os.Open("./settings.json")
	if err != nil {
		os.Create("./settings.json")
		jsonFile, err = os.Open("./settings.json")
		if err != nil {
			log.Println(err)
		}
		jsonData, _ := json.Marshal(Config{
			AuthKey:       "",
			TwoCaptchaKey: "",
			LatestRainID:  0,
			Webhook:       "",
			Debug:         false,
			SafeMode:      true,
		})
		ioutil.WriteFile("./settings.json", jsonData, 0644)
		fmt.Println(err)
		time.Sleep(time.Second * 5)
		os.Exit(1)
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
		if con.Debug {
			color.Error.Tips("Error restarting..: %s", err)
			time.Sleep(time.Second * 5)
		}

		goto restart
	}

	for {
		PotID = con.LatestRainID
		err := c.WriteMessage(websocket.TextMessage, []byte("3"))
		if err != nil {
			goto restart
		}
		if con.Debug {
			color.Debug.Tips("Sent ping to " + u.String())
		}

		_, message, err := c.ReadMessage()
		if err != nil {
			if con.Debug {
				color.Error.Tips("Error reading message: " + err.Error())
			}

			goto restart
		}
		if con.Debug {
			color.Debug.Tips("Message received: " + string(message))
		}
		if strings.Contains(string(message), "events:rain:updatePotVariables") {
			UpdateRain(&con, message)

		} else if strings.Contains(string(message), "updateBalance") {
			type Data struct {
				Value int `json:"value"`
			}
			//42["user:updateBalance",{"value":39,"time":16524//78879050}]	1652478879.099368
			var data Data
			f := strings.Split(string(message), `ce",`)[1]
			f = strings.Split(f, "]")[0]
			err = json.Unmarshal([]byte(f), &data)
			if err != nil {
				if con.Debug {
					color.Error.Tips("Error parsing json: " + err.Error())
				}
			}
			AuthorizationConfig.UserData.Balance = data.Value
			if AuthorizationConfig.UserData.Balance > 1000 {
				usdBalance := float64(AuthorizationConfig.UserData.Balance) / 100
				floatToString := strconv.FormatFloat(usdBalance, 'f', 2, 64)
				type SendData struct {
					Type         string `json:"type"`
					Amount       int    `json:"amount"`
					Instant      bool   `json:"instant"`
					DummyAssetID int    `json:"dummyAssetId"`
				}
				var sendData SendData
				sendData.Type = "WITHDRAW"
				sendData.Amount = AuthorizationConfig.UserData.Balance
				sendData.Instant = false
				sendData.DummyAssetID = 0
				sendDataJson, _ := json.Marshal(sendData)

				req, _ := http.NewRequest("POST", "https://rblxwild.com/api/trading/robux/request-exchange", bytes.NewBuffer(sendDataJson))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", con.AuthKey)
				resp, _ := http.DefaultClient.Do(req)

				if resp.StatusCode == 200 {

					color.Debug.Tips("Sent withdraw request")
					jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Withdraw successful","color":2303786,"fields":[{"name":"Amount withdrawed","value":"$` + floatToString + ` USDT","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://bitcoin.org/img/icons/opengraph.png?1651392467"}}],"username":"Withdraw","avatar_url":"https://bitcoin.org/img/icons/opengraph.png?1651392467","attachments":[]}`)
					http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
				} else {
					if con.Debug {
						color.Error.Tips("Error sending withdraw request: " + resp.Status)
						jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Withdraw failed","color":2303786,"fields":[{"name":"Error","value":"` + resp.Status + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://bitcoin.org/img/icons/opengraph.png?1651392467"}}],"username":"Error withdrawing","avatar_url":"https://bitcoin.org/img/icons/opengraph.png?1651392467","attachments":[]}`)
						http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
					}
				}
			}

		} else if strings.Contains(string(message), "ENDING") && strings.Contains(string(message), "newState") {
			go func() {
				if con.SafeMode {
					num := rand.Intn(100)
					if num > 95 {
						con.LatestRainID += 1
						usdBalance := float64(AuthorizationConfig.UserData.Balance) / 100
						floatToString := strconv.FormatFloat(usdBalance, 'f', 2, 64)
						saveJson, _ := json.Marshal(con)
						ioutil.WriteFile("./settings.json", saveJson, 0644)
						jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Skipping hash","color":2303786,"fields":[{"name":"Time Elapsed","value":"nil","inline":true},{"name":"Balance","value":"$` + floatToString + ` USDT","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://bitcoin.org/img/icons/opengraph.png?1651392467"}}],"username":"Skipped","avatar_url":"https://bitcoin.org/img/icons/opengraph.png?1651392467","attachments":[]}`)
						http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
						return
					}
				}
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
				captcha := api2captcha.NewClient(con.TwoCaptchaKey)
				if con.Debug {
					color.Debug.Tips("Created new client")
				}
				captcha.DefaultTimeout = 120
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

				code, err := captcha.Solve(cap.ToRequest())
				if err != nil {
					log.Println(err)
					con.LatestRainID += 1

					saveJson, _ := json.Marshal(con)
					ioutil.WriteFile("./settings.json", saveJson, 0644)
					color.Debug.Tips("Captcha timeout TIME ELAPSED: " + strconv.Itoa(int(time.Since(timer))))
					usdBalance := float64(AuthorizationConfig.UserData.Balance) / 100
					floatToString := strconv.FormatFloat(usdBalance, 'f', 2, 64)
					jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Failed to mine","color":2303786,"fields":[{"name":"Time Elapsed","value":"` + time.Since(timer).String() + `","inline":true},{"name":"Balance","value":"$` + floatToString + ` USDT","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://bitcoin.org/img/icons/opengraph.png?1651392467"}}],"username":"Failed","avatar_url":"https://bitcoin.org/img/icons/opengraph.png?1651392467","attachments":[]}`)
					req, _ := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
					if con.Debug {
						if data, err := ioutil.ReadAll(req.Body); err != nil {
							log.Println("Error:", err)
						} else {
							color.Debug.Tips("Webhook sent. " + req.Status + "\n" + string(data))
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
					return
				} else {
					if strings.Contains(string(data), "true") {
						if con.Debug {
							color.Debug.Tips("Joined the rain successfully")
						}
						usdBalance := float64(AuthorizationConfig.UserData.Balance) / 100
						floatToString := strconv.FormatFloat(usdBalance, 'f', 2, 64)
						jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Mined successfully","color":2303786,"fields":[{"name":"Time Elapsed","value":"` + time.Since(timer).String() + `","inline":true},{"name":"Balance","value":"$` + floatToString + ` USDT","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://bitcoin.org/img/icons/opengraph.png?1651392467"}}],"username":"Success","avatar_url":"https://bitcoin.org/img/icons/opengraph.png?1651392467","attachments":[]}`)
						req, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
						if err != nil {
							color.Error.Tips("Error sending webhook: " + err.Error())
						}
						if con.Debug {
							color.Debug.Tips("Webhook sent. " + req.Status)
						}
						err = c.WriteMessage(websocket.TextMessage, []byte(`42["upgrader:play",{"inputBetAmount":5,"outputBetAmount":6,"rollType":"UNDER","fastAnimation":false}]`))
						if err != nil {
							return
						}
					} else {
						type DataBy struct {
							Success bool   `json:"success"`
							Message string `json:"message"`
							Elapsed int    `json:"elapsed"`
						}

						var DataByData DataBy
						json.Unmarshal(data, &DataByData)
						jsonData1 := []byte(`{"content":"@everyone","embeds":[{"title":"Error mining","color":2303786,"fields":[{"name":"Time Elapsed","value":"` + time.Since(timer).String() + `","inline":true},{"name":"Error","value":"` + DataByData.Message + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"thumbnail":{"url":"https://bitcoin.org/img/icons/opengraph.png?1651392467"}}],"username":"Error","avatar_url":"https://bitcoin.org/img/icons/opengraph.png?1651392467","attachments":[]}`)

						req1, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData1))
						log.Println(oldPrice, string(data), time.Since(timer).String())
						if err != nil {
							color.Error.Tips("Error sending webhook: " + err.Error())
							if data, err := ioutil.ReadAll(req1.Body); err != nil {
								log.Println("Error:", err)
							} else {
								color.Debug.Tips("Webhook sent. " + req1.Status + "\n" + string(data))

							}
						}

					}

				}
				con.LatestRainID += 1
				RainsJoined += 1
				saveJson, _ := json.Marshal(con)
				ioutil.WriteFile("./settings.json", saveJson, 0644)
			}()
		}
	}
}

func UpdateRain(con *Config, message []byte) {
	type Data struct {
		NewPrize              int `json:"newPrize"`
		NewJoinedPlayersCount int `json:"newJoinedPlayersCount"`
	}

	var data Data
	jsonData := strings.Split(string(message), `",`)[1]
	jsonData = strings.Replace(jsonData, "]", "", -1)
	err := json.Unmarshal([]byte(jsonData), &data)

	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	if con.Debug {
		color.Debug.Tips("Data structure complete")
		color.Debug.Tips("Pointer variable created for the data structure")
		color.Debug.Tips("NewPrize: %d", data.NewPrize)
		color.Debug.Tips("NewJoinedPlayersCount: %d", data.NewJoinedPlayersCount)
		color.Debug.Tips("Json unmarshalled")
		color.Debug.Tips("Updated old price variable with the new variable")
	}

	oldPrice = data.NewPrize
}
func GUI() {
	for {
		cmd1 := exec.Command(`echo -n -e "\033]0;Rain ` + strconv.Itoa(oldPrice) + `\007"`)
		cmd1.Stdout = os.Stdout
		cmd1.Run()
		fmt.Printf(GetLogo())
		if AuthorizationConfig == nil {
			color.Magenta.Printf(`
					╔═ USDT balance: $0.00
					╠═ BTC balance: 0.00000000 BTC
					╚══ Username: nil`)
		} else {
			usdBalance := float64(AuthorizationConfig.UserData.Balance) / 100
			floatToString := strconv.FormatFloat(usdBalance, 'f', 2, 64)
			btcBalance := usdBalance / 29836
			floatToString2 := strconv.FormatFloat(btcBalance, 'f', 8, 64)
			color.Magenta.Printf(`

					╔═ USDT balance: $` + floatToString + `
					╠═ BTC balance: ` + floatToString2 + ` BTC
					╚══ Username: ` + AuthorizationConfig.UserData.DisplayName)
		}
		time.Sleep(time.Second * 1)
		cmd := exec.Command(`clear`)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}
func EstimateEndPrice(newPrice int) int {
	return ((newPrice - oldPrice) * 30) + oldPrice
}
