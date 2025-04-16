package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// โครงสร้าง RSS Feed
type RSS struct {
    Channel Channel `xml:"channel"`
}

type Channel struct {
    Items []Item `xml:"item"`
}

type Item struct {
    Title       string `xml:"title"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
    Link        string `xml:"link"`
}

func fetchEarthquakeData() ([]Item, error) {
    // ดึงข้อมูล RSS
    resp, err := http.Get("https://earthquake.tmd.go.th/feed/rss_tmd.xml")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // parse XML
    var rss RSS
    err = xml.Unmarshal(body, &rss)
    if err != nil {
        return nil, err
    }

    // filter เอาเฉพาะที่ไม่ใช่เมียนมา
    var filtered []Item
    for _, item := range rss.Channel.Items {
        if !strings.Contains(item.Title, "ประเทศเมียนมา") && !strings.Contains(item.Title, "Myanmar") {
            filtered = append(filtered, item)
        }
    }

    return filtered, nil
}

func fetchFloodData() ([]Item, error) {
    // ดึงข้อมูล RSS
    resp, err := http.Get("https://disaster.gistda.or.th/api/1.0/documents/flood/1day/wms")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // parse XML
    var rss RSS
    err = xml.Unmarshal(body, &rss)
    if err != nil {
        return nil, err
    }

    // filter เอาเฉพาะที่ไม่ใช่เมียนมา
    var filtered []Item
    for _, item := range rss.Channel.Items {
        if !strings.Contains(item.Title, "ประเทศเมียนมา") && !strings.Contains(item.Title, "Myanmar") {
            filtered = append(filtered, item)
        }
    }

    return filtered, nil
}


func main() {
    r := gin.Default()

    r.GET("/earthquakes", func(c *gin.Context) {
        data, err := fetchEarthquakeData()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        alert := "พบเหตุการณ์แผ่นดินไหวที่ไม่ใช่ประเทศเมียนมา จำนวน " + fmt.Sprint(len(data)) + " รายการ"

        c.JSON(http.StatusOK, gin.H{
            "alert": alert,
            "data":  data,
        })
    })

	r.GET("/flood", func(c *gin.Context) {
        data, err := fetchFloodData()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        alert := "พบเหตุการณ์แผ่นดินไหวที่ไม่ใช่ประเทศเมียนมา จำนวน " + fmt.Sprint(len(data)) + " รายการ"

        c.JSON(http.StatusOK, gin.H{
            "alert": alert,
            "data":  data,
        })
    })

    r.Run(":8080") // API ใช้ที่ http://localhost:8080/earthquakes
}
