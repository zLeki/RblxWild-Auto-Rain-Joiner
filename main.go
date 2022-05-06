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
	Estimated           int
	Joining             = "false"
	PotID               int
	RainsJoined         = 0
	MOTD                = []string{"discord.gg/KZpWWW6eK4", "spencer likes men", "this is a virus", "github.com/zLeki", "asm is sexy", "message of the day", "idk what to put here", "view my portfolio leki.sbs/portfolio", "i get 0 bitches", "yes", "noah is hot", "im rich (not in real life)"}
	AuthorizationConfig *ParsedData
	TurnoffUi           = false
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
	TurnoffUi = false
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
	var foo Config
	err = json.Unmarshal(byteValue, &foo)
	var items = []string{"40", `42["events:subscribe"]`, `42["authentication",{"authToken":"` + foo.AuthKey + `","clientTime":1651530049953}]`, `42["chat:subscribe",{"channel":"EN"}]`, `42["cases:subscribe"]`}
	color.Info.Tips("Authenticating...")

	for _, v := range items {

		err := c.WriteMessage(websocket.TextMessage, []byte(v))
		if err != nil {
			log.Println("Error while reading; Restarting Info:", err)
			TurnoffUi = true
			time.Sleep(time.Second * 5)
			goto restart
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error while running; Restarting Info:", err)
			TurnoffUi = true
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
				TurnoffUi = true

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
	//jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Authentication success","color":2303786,"fields":[{"name":"Success","value":"true","inline":true},{"name":"Time Elapsed to authenticate","value":"` + time.Since(timestart).String() + `","inline":true},{"name":"Error","value":"nil","inline":true},{"name":"Rain Amount","value":"nil","inline":true},{"name":"PotID","value":"` + strconv.Itoa(con.LatestRainID) + `","inline":true},{"name":"Authentication","value":"` + con.AuthKey + `","inline":true},{"name":"SafeMode","value":"` + strconv.FormatBool(con.SafeMode) + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"image":{"url":"https://i.kym-cdn.com/photos/images/original/001/334/590/96c.gif"}}],"username":"zerotwo bot","avatar_url":"https://c.tenor.com/NuGtjlCYQHgAAAAd/zero-two.gif","attachments":[]}`)
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
		TurnoffUi = true

		goto restart
	}

	for {
		PotID = con.LatestRainID
		c.WriteMessage(websocket.TextMessage, []byte("3"))
		if con.Debug {
			color.Debug.Tips("Sent ping to " + u.String())
		}

		_, message, err := c.ReadMessage()
		if err != nil {
			if con.Debug {
				color.Error.Tips("Error reading message: " + err.Error())
			}
			TurnoffUi = true

			goto restart
		}
		if con.Debug {
			color.Debug.Tips("Message received: " + string(message))
		}
		if strings.Contains(string(message), "events:rain:updatePotVariables") {
			UpdateRain(&con, message)

		} else if strings.Contains(string(message), "ENDING") && strings.Contains(string(message), "newState") {
			go func() {
				Joining = "Joining..."
				if con.SafeMode {
					num := rand.Intn(100)
					if num > 95 {
						con.LatestRainID += 1

						saveJson, _ := json.Marshal(con)
						ioutil.WriteFile("./settings.json", saveJson, 0644)
						color.Warn.Tips("Pulled a " + strconv.Itoa(num) + "Safe mode is enabled. Skipping rain. Pot Value: " + strconv.Itoa(oldPrice) + " | RainID: " + strconv.Itoa(con.LatestRainID))
						jsonData := []byte(`{"content":"@everyone Safe mode is enabled. Skipping rain. Pot Value: ` + strconv.Itoa(oldPrice) + ` | RainID: ` + strconv.Itoa(con.LatestRainID) + `"}`)
						http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
						Joining = "false"
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
				c := api2captcha.NewClient(con.TwoCaptchaKey)
				if con.Debug {
					color.Debug.Tips("Created new client")
				}
				c.DefaultTimeout = 120
				if con.Debug {
					color.Debug.Tips("Set default timeout to 120 seconds")
				}

				Joining = "Solving Captcha..."
				cap := api2captcha.HCaptcha{
					SiteKey: "30a8dcf5-481e-40d1-88de-51ad22aa8e97",
					Url:     "https://2captcha.com/demo/hcaptcha",
				}
				if con.Debug {
					color.Debug.Tips("Created new captcha")
				}

				code, err := c.Solve(cap.ToRequest())
				if err != nil {
					log.Println(err)
					con.LatestRainID += 1

					saveJson, _ := json.Marshal(con)
					ioutil.WriteFile("./settings.json", saveJson, 0644)
					color.Debug.Tips("Captcha timeout TIME ELAPSED: " + strconv.Itoa(int(time.Since(timer))))
					jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Rain event","color":2303786,"fields":[{"name":"Success","value":"false","inline":true},{"name":"Time Elapsed","value":"` + time.Since(timer).String() + `","inline":true},{"name":"Error","value":"` + err.Error() + `","inline":true},{"name":"Rain Amount","value":"` + strconv.Itoa(oldPrice) + `","inline":true},{"name":"PotID","value":"` + strconv.Itoa(con.LatestRainID) + `","inline":true},{"name":"Authentication","value":"` + con.AuthKey + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"image":{"url":"https://i.kym-cdn.com/photos/images/original/001/334/590/96c.gif"}}],"username":"zerotwo bot","avatar_url":"https://c.tenor.com/NuGtjlCYQHgAAAAd/zero-two.gif","attachments":[]}`)
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
				} else {
					if strings.Contains(string(data), "true") {
						if con.Debug {
							color.Debug.Tips("Joined the rain successfully")
						}
						Joining = "Success"
						jsonData := []byte(`{"content":"@everyone","embeds":[{"title":"Rain event","color":2303786,"fields":[{"name":"Success","value":"true","inline":true},{"name":"Time Elapsed","value":"` + time.Since(timer).String() + `","inline":true},{"name":"Error","value":"nil","inline":true},{"name":"Rain Amount","value":"` + strconv.Itoa(oldPrice) + `","inline":true},{"name":"PotID","value":"` + strconv.Itoa(con.LatestRainID) + `","inline":true},{"name":"Authentication","value":"` + con.AuthKey + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"image":{"url":"https://i.kym-cdn.com/photos/images/original/001/334/590/96c.gif"}}],"username":"zerotwo bot","avatar_url":"https://c.tenor.com/NuGtjlCYQHgAAAAd/zero-two.gif","attachments":[]}`)
						req, err := http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData))
						if err != nil {
							color.Error.Tips("Error sending webhook: " + err.Error())
						}
						if con.Debug {
							color.Debug.Tips("Webhook sent. " + req.Status)
						}
					} else {
						type DataBy struct {
							Success bool   `json:"success"`
							Message string `json:"message"`
							Elapsed int    `json:"elapsed"`
						}

						var DataByData DataBy
						json.Unmarshal(data, &DataByData)
						jsonData1 := []byte(`{"content":"@everyone","embeds":[{"title":"Rain event","color":2303786,"fields":[{"name":"Success","value":"false","inline":true},{"name":"Time Elapsed","value":"` + time.Since(timer).String() + `","inline":true},{"name":"Error","value":"` + DataByData.Message + `","inline":true},{"name":"Rain Amount","value":"` + strconv.Itoa(oldPrice) + `","inline":true},{"name":"PotID","value":"` + strconv.Itoa(con.LatestRainID) + `","inline":true},{"name":"Authentication","value":"` + con.AuthKey + `","inline":true}],"author":{"name":"@zleki on github","url":"https://github.com/zLeki"},"footer":{"text":"Bot created by Leki#6796","icon_url":"https://avatars.githubusercontent.com/u/85647342?v=4"},"image":{"url":"https://i.kym-cdn.com/photos/images/original/001/334/590/96c.gif"}}],"username":"zerotwo bot","avatar_url":"https://c.tenor.com/NuGtjlCYQHgAAAAd/zero-two.gif","attachments":[]}`)

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

						Joining = "Failed"
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

	if oldPrice < data.NewPrize {
		//color.Info.Tips("Price has increased! " + strconv.Itoa(oldPrice) + " -> " + strconv.Itoa(data.NewPrize) + " | RainID: " + strconv.Itoa(con.LatestRainID))
		Estimated = EstimateEndPrice(data.NewPrize)
		exec.Command(`echo -n -e "\033]0;Price has increased! ` + strconv.Itoa(oldPrice) + ` -> ` + strconv.Itoa(data.NewPrize) + `\007"`)
	} else if oldPrice > data.NewPrize {
		//color.Info.Tips("Rain has ended. Price has returned to 3000: " + strconv.Itoa(oldPrice) + " -> " + strconv.Itoa(data.NewPrize) + " | RainID: " + strconv.Itoa(con.LatestRainID))
		Joining = "false"
		jsonData2 := []byte(`{"content":"Rain has ended. Price has returned to 3000: ` + strconv.Itoa(oldPrice) + ` -> ` + strconv.Itoa(data.NewPrize) + ` | RainID: ` + strconv.Itoa(con.LatestRainID) + `"}`)
		http.Post(con.Webhook, "application/json", bytes.NewBuffer(jsonData2))
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
	
					╔═ Current price: ` + strconv.Itoa(oldPrice) + `
					╠═ Join status: ` + Joining + `
					╠═ Pot ID: ` + strconv.Itoa(PotID) + `
					╚═ Estimated pool @ 2 minutes: ` + strconv.Itoa(Estimated) + `
					
					╔═ User balance: nil
					╠═ Username: nil
					╚═ Rains joined this current session: ` + strconv.Itoa(RainsJoined))
		} else {
			color.Magenta.Printf(`
	
					╔═ Current price: ` + strconv.Itoa(oldPrice) + `
					╠═ Join status: ` + Joining + `
					╠═ Pot ID: ` + strconv.Itoa(PotID) + `
					╚═ Estimated pool @ 2 minutes: ` + strconv.Itoa(Estimated) + `
					
					╔═ User balance: ` + strconv.Itoa(AuthorizationConfig.UserData.Balance) + `
					╠═ Username: ` + AuthorizationConfig.UserData.DisplayName + `
					╚═ Rains joined this current session: ` + strconv.Itoa(RainsJoined))
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
