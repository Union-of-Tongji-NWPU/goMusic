package tool

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"musicDance/model"
	"os"
	"regexp"
	"strings"
)

func GetSong(song string) bool {
	song = ProcessSong(song)
	if !findSong(song) {
		return clawSong(song)
	} else {
		return true
	}
}

func findSong(song string) bool {
	file, err := os.Open(model.ResourceDir + "/" + song + ".txt")
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func ProcessSong(song string) string {
	songs := strings.FieldsFunc(song, split)
	var result string
	for i, s := range songs {
		if i != 0 {
			result += "-"
		}
		result += strings.ToLower(s)
	}
	return result
}

func clawSong(song string) bool {

	// 起始Url
	startUrl := "https://noobnotes.net/" + song + "?solfege=false"

	// 创建Collector
	collector := colly.NewCollector()
	extensions.RandomUserAgent(collector)

	//// 设置抓取频率限制
	//collector.Limit(&colly.LimitRule{
	//	DomainGlob:  "*",
	//	RandomDelay: 5 * time.Second, // 随机延迟
	//})

	// 异常处理
	var result string
	collector.OnError(func(response *colly.Response, err error) {
		log.Println(err.Error())
		result = ""
	})

	collector.OnRequest(func(request *colly.Request) {
		log.Println("start visit: ", request.URL.String())
	})

	// 解析列表
	collector.OnHTML("body", func(element *colly.HTMLElement) {
		selection := element.DOM.Find("div.post-content")
		// 依次遍历所有的p节点
		selection.Find("p").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			// 进一步处理
			if text != " " {
				result += parseText(text) + "\n"
			}
		})
		fmt.Println(result)
	})
	// 起始入口
	collector.Visit(startUrl)
	if result != "" {
		// 存储文件名
		fName := model.ResourceDir + "/" + song + ".txt"
		file, err := os.Create(fName)
		if err != nil {
			log.Fatalf("创建文件失败 %q: %s\n", fName, err)
			return false
		}
		defer file.Close()
		_, err = file.WriteString(result)
		if err != nil {
			panic(err)
			return false
		}
		return true
	} else {
		return false
	}
}
func split(r rune) bool {
	return r == ' ' || r == '-' || r == '’'
}
func deleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	regStr := "\\s{2,}"                         //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regStr)            //编译正则表达式
	s2 := make([]byte, len(s1))                 //定义字符数组切片
	copy(s2, s1)                                //将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spcIndex) > 0 {                     //找到适配项
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...) //删除多余空格
		spcIndex = reg.FindStringIndex(string(s2))           //继续在字符串中搜索
	}
	return string(s2)
}

func spaceConsecutiveLetter(s string) string {
	regStr := "[A-G]{2}"
	reg, _ := regexp.Compile(regStr) //编译正则表达式
	s1 := s
	spcIndex := reg.FindStringIndex(s) //在字符串中搜索
	for len(spcIndex) > 0 {            //找到适配项
		s1 = s[:spcIndex[0]+1] + " " + s[spcIndex[1]-1:]
		s = s1
		spcIndex = reg.FindStringIndex(s) //继续在字符串中搜索
	}
	return s1
}

func parseDu(du string) string {
	var letter, well, number string
	number = "4"
	for _, c := range du {
		if c >= 'A' && c <= 'G' {
			letter = string(c - ('A' - 'a'))
		} else if c == '#' {
			well = "#"
		} else if c == '_' {
			number = "2"
		} else if c == '.' {
			number = "3"
		} else if c == '^' {
			number = "5"
		} else if c == '*' {
			number = "6"
		}
	}
	return letter + well + number
}
func parseText(text string) string {
	du := strings.Split(text, "\n")[0]
	du = strings.ReplaceAll(du, " ", " ")
	du = strings.ReplaceAll(du, "#", "# ")
	du = strings.ReplaceAll(du, "^", " ^")
	du = deleteExtraSpace(du)
	du = spaceConsecutiveLetter(du)
	dus := strings.FieldsFunc(du, split)
	var result string
	for _, v := range dus {
		result += parseDu(v) + " "
	}
	return result
}
