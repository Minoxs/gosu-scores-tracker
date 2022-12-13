package main

import "osu-phantom/libgo"

func main() {
	//var file, _ = ioutil.ReadFile("test.osu")
	//var req, _ = http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/diff?mods=%d", 1<<4), bytes.NewBuffer(file))
	//req.Header.Set("Content-Type", "text/osu")
	//var res, _ = http.DefaultClient.Do(req)
	//io.Copy(os.Stdout, res.Body)

	libgo.GetPP()
	//Client := &client.PhantomClient{
	//	Username: "minoxs",
	//}
	//Client.Login()
	//
	//scores := Client.GetRecentScores()
	//fmt.Println("Recent Scores: ", scores)
}
