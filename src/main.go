package main

import (
	"fmt"
	"io/ioutil"

	"osu-phantom/src/calculator"
)

func main() {
	//var file, _ = ioutil.ReadFile("test.osu")
	//var req, _ = http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/diff?mods=%d", 1<<4), bytes.NewBuffer(file))
	//req.Header.Set("Content-Type", "text/osu")
	//var res, _ = http.DefaultClient.Do(req)
	//io.Copy(os.Stdout, res.Body)

	var buf, _ = ioutil.ReadFile("./test.osu")

	fmt.Println(calculator.GetPPFromMap(buf, 476, 281, 37, 4, 0, calculator.HR, calculator.OSU))
	//Client := &client.PhantomClient{
	//	Username: "minoxs",
	//}
	//Client.Login()
	//
	//scores := Client.GetRecentScores()
	//fmt.Println("Recent Scores: ", scores)
}
