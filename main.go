package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var HOST = "localhost"
var IS_DEBUG = false;

func init() {
	// ランダムシード
	rand.Seed(time.Now().UnixNano())

	// ホスト名
	if !IS_DEBUG {
		HOST = os.Getenv("HOSTNAME")
	}
}

/**************************************************************************************************/
/*!
 *  エントリポイント
 */
/**************************************************************************************************/
func main() {
	// routing
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/images/", imagesHandler)
	http.HandleFunc("/upload", uploadHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println(err)
	}
}

/**************************************************************************************************/
/*!
 *  疎通確認
 *
 *  \param   w : Writer
 *  \param   r : リクエスト
 */
/**************************************************************************************************/
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

/**************************************************************************************************/
/*!
 *  UPLOAD処理
 *
 *  \param   w : Writer
 *  \param   r : リクエスト
 */
/**************************************************************************************************/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// ディレクトリ取得
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// imagesがなければ作成
	imageDir := path.Join(dir, "images")
	if err := os.Mkdir(imageDir, 0755); err != nil && !os.IsExist(err) {
		fmt.Fprintln(w, err)
		return
	}

	// ファイルデータを取得
	file, _, err := r.FormFile("imagedata")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	// MD5で被らなそうな名前をつける
	key := strconv.FormatInt(rand.Int63(), 10)
	timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	h := md5.New()
	io.WriteString(h, key+timeStr)
	baseName := fmt.Sprintf("%x", h.Sum(nil))

	// ファイル生成
	imageFile := path.Join(imageDir, baseName)
	out, err := os.Create(imageFile)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer out.Close()

	// ファイルにデータを流す
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// 生成されたURLを返す
	fmt.Fprintf(w, "http://%s/images/%s", HOST, baseName)
}

/**************************************************************************************************/
/*!
 *  画像表示
 *
 *  \param   w : Writer
 *  \param   r : リクエスト
 */
/**************************************************************************************************/
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	// ディレクトリ取得
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	// ファイルへのURL
	imagefile := path.Join(dir, r.URL.Path)
	http.ServeFile(w, r, imagefile)
}
