package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

/**************************************************************************************************/
/*!
 *  エントリポイント
 */
/**************************************************************************************************/

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/images/", imagesHandler)
	http.HandleFunc("/upload", uploadHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
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

	// 被らなそうな名前をつけて、ファイル生成
	baseName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
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
	fmt.Fprintf(w, "http://%s/images/%s", r.Host, baseName)
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
