package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hayrat/gostr"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Okul struct {
	Sehir string `json:"Sehir"`
	Ilce  string `json:"İlce"`
	Adi   string `json:"Adi"`
}

func okulGetir(url string) ([]Okul, error) {
	var resultArray []Okul

	e := make(map[string]Okul)

	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://www.meb.gov.tr/baglantilar/okullar/index.php?ILKODU=" + strconv.Itoa(i))
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		doc.Find("#icerik-listesi > tbody >tr> td> a").Each(func(i int, selection *goquery.Selection) {
			if _, ok := selection.Attr("style"); !ok {
				r, _ := regexp.Compile("^[A-ZÇĞİÖŞÜ]*")

				data := selection.Text()
				il := r.FindString(data)
				data = data[len(il)+3:]
				ilce := r.FindString(data)
				data = data[len(ilce)+3:]
				fmt.Println(data)
				//data = strings.ToTitleSpecial(unicode.TurkishCase, data)
				//a := strings.Title(gostr.TurkceKucukHarfeCevir(data))
				//fmt.Println(a)
				data = basHarifBuyut(data)
				data = strings.ReplaceAll(data,"Iıı","III")
				data = strings.ReplaceAll(data,"Iı","II")
				data = strings.ReplaceAll(data,"İlkokuku","İlkokulu")
				data = strings.ReplaceAll(data,"İlk Okulu","İlkokulu")
				data = strings.ReplaceAll(data,"ilkokulu","İlkokulu")
				data = strings.ReplaceAll(data,"İllkokulu","İlkokulu")
				data = strings.ReplaceAll(data,"İlokulu","İlkokulu")
				data = strings.ReplaceAll(data,"İlkolkulu","İlkokulu")
				data = strings.ReplaceAll(data,"Ilkokulu","İlkokulu")
				data = strings.ReplaceAll(data,"İkokulu","İlkokulu")
				data = strings.ReplaceAll(data,"İlkkulu","İlkokulu")
				data = strings.ReplaceAll(data,"Öğretmen Evi","Öğretmenevi")
				data = strings.ReplaceAll(data,"Ögretmen Evi","Öğretmenevi")
				data = strings.ReplaceAll(data,"Ögretmenevi","Öğretmenevi")
				data = strings.ReplaceAll(data,"özel Eğitim"," Özel Eğitim")
				data = strings.ReplaceAll(data,"Ana Okulu","Anaokulu")
				data = strings.ReplaceAll(data,"Anaoku","Anaokulu")
				data = strings.ReplaceAll(data,"Ortoakulu","Ortaokulu")
				data = strings.ReplaceAll(data,"Ortokulu","Ortaokulu")
				data = strings.ReplaceAll(data,"Mes.eğt.merk.","Mesleki Eğitim Merkezi")
				data = strings.ReplaceAll(data,"Mes.eğt.m.","Mesleki Eğitim Merkezi")
				data = strings.ReplaceAll(data,"Meslekî Eğitim Merkezi","Mesleki Eğitim Merkezi")
				data = strings.ReplaceAll(data,"\n","")
				data = strings.ReplaceAll(data,"Lis.","Lisesi")
				data = strings.ReplaceAll(data,"Yusuf Beylem","Yusuf Beylem İlkokulu")
				data = strings.ReplaceAll(data,"Kız Ybo","Kız Yatılı Bölge Ortaokulu")
				data = strings.ReplaceAll(data,"Bekir - Sacide - Filiz","Bekir - Sacide - Filiz Anaokulu")

				okulAdi := selection.Text()
				okulAdiKesilmis := strings.Split(okulAdi, " - ")
				for i := 0; i < len(okulAdiKesilmis); i++ {
					a := strings.TrimSpace(okulAdiKesilmis[i])
					b := strings.Title(strings.ToLower(a))
					okulAdiKesilmis[i] = b
					okulAdiKesilmis[2] = data
				}
				e[okulAdi] = Okul{Sehir: okulAdiKesilmis[0], Ilce: okulAdiKesilmis[1], Adi: okulAdiKesilmis[2]}
				resultArray = append(resultArray, Okul{
					Sehir: okulAdiKesilmis[0],
					Ilce:  okulAdiKesilmis[1],
					Adi:   okulAdiKesilmis[2],
				})
				fmt.Println(len(resultArray))
			}
		})
	}
	return resultArray, nil
}
func basHarifBuyut(text string) string {
	if len(text) <= 1 {
		return text
	}
	sonuc := make([]string, 0)
	text = gostr.TurkceKucukHarfeCevir(text)
	kelimeler := strings.Split(text, " ")
	fmt.Println(kelimeler)
	for _, kelime := range kelimeler {
		if len(strings.TrimSpace(kelime)) == 0 {
			continue
		}
		kodlanmis := []rune(kelime)
		ilkHarf := string(kodlanmis[0])
		ilkHarfRune := []rune(ilkHarf)
		gerisi := string(kodlanmis[len(ilkHarfRune):])


		sonuc = append(sonuc, gostr.TurkceBuyukHarfeCevir(ilkHarf)+gerisi)
		fmt.Println(sonuc)
	}
	return strings.Join(sonuc, " ")
}
func main() {
	school, err := okulGetir("http://www.meb.gov.tr/baglantilar/okullar/index.php")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	byteData, _ := json.MarshalIndent(&school, "", " ")
	err = ioutil.WriteFile("okul.json", byteData, 0644)
	if err != nil {
		fmt.Println("Yazamadı +" + err.Error())
		fmt.Println("Sübhanallah")
	} else {
		fmt.Println("Elhamdülillah")
	}
}
