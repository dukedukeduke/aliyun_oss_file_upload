package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"strings"
)

func main(){
	logFile, err  := os.OpenFile("aliyunOssUpload.log", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil{
		fmt.Fprintf(os.Stderr,"创建日志失败")
	}

	logger := log.New(logFile, "", log.LstdFlags|log.Llongfile)

	domain := flag.String("domain", "", "aliyun oss domain")
	bucket := flag.String("bucket", "", "aliyun oss bucket")
	key := flag.String("key", "", "aliyun ram key")
	secret := flag.String("secret", "", "aliyun ram secret")
	firebase := flag.String("firebase", "", "firebase storage url")
	local := flag.String("local", "", "local file absolute path")
	flag.Parse()

	if *domain == "" || *bucket == "" || *key == "" || *secret== "" || *firebase=="" || *local == ""{
		fmt.Fprintf(os.Stderr,fmt.Sprintf("Parameter error for domain:%v, bucket:%v, local file:%v\n", *domain, *bucket, * local))
		logger.Fatalln(fmt.Sprintf("Parameter error for domain:%v, bucket:%v, local file:%v", *domain, *bucket, * local))
	}
	if _, err := os.Stat(*local); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr,fmt.Sprintf("Local file not exists for domain:%v, bucket:%v, local file:%v\n", *domain, *bucket, * local))
		logger.Fatalln(fmt.Sprintf("Local file not exists for domain:%v, bucket:%v, local file:%v", *domain, *bucket, * local))
	}
	client, err := oss.New(*domain, *key, *secret)
	if err != nil {
		fmt.Fprintf(os.Stderr,fmt.Sprintf("Get new client failed for domain:%v\n", *domain))
		logger.Fatalln(fmt.Sprintf("Get new client failed for domain:%v", *domain))
	}
	bucketTmp, err := client.Bucket(*bucket)
	if err != nil{
		fmt.Fprintf(os.Stderr,fmt.Sprintf("Get bucket object failed for domain:%v, bucket:%v\n", *domain, *bucket))
		logger.Fatalln(fmt.Sprintf("Get bucket object failed for domain:%v, bucket:%v", *domain, *bucket))
	}

	destDirSplitByO := strings.Split(*firebase, "/o/")[1:]
	var destDir string
	for i :=0;i < len(destDirSplitByO); i++{
		destDir += destDirSplitByO[i]
	}

	destDir = strings.Replace(destDir, "%2F","/", -1)

	err = bucketTmp.PutObjectFromFile(destDir, *local)
	if err != nil {
		fmt.Fprintf(os.Stderr,fmt.Sprintf("Upload file to aliyun oss failed for domain:%v, local:%v, error:%v\n", *domain, *local, err))
		logger.Fatalln(fmt.Sprintf("Upload file to aliyun oss failed for domain:%v, local:%v, error:%v", *domain, *local, err))
	}
	logger.Println("Upload success")
	fmt.Fprintf(os.Stdout,"Upload success\n")
}
