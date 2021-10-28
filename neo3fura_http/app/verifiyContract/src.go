package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/nef"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 定义主网和测试网节点常量，之后根据部署主网测试网选择响应的结点
const RPCNODEMAIN = "https://neofura.ngd.network:1927"

const RPCNODETEST = "https://testneofura.ngd.network:444"

//定义主网和测试往数据库结构
type Config struct {
	Database_main struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_main"`
	Database_test struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_test"`
}
//定义http应答返回格式
type jsonResult struct {
	Code int
	Msg string
}
//定义插入VerifiedContract表的数据格式, 记录被验证的合约
type insertVerifiedContract struct {
	Hash string
	Id int
	Updatecounter int

}
//定义插入ContractSourceCode表的数据格式，记录被验证的合约源代码
type insertContractSourceCode struct {
	Hash string
	Updatecounter int
	FileName string
	Code string
}

func multipleFile(w http.ResponseWriter, r *http.Request) {
	//定义value 为string 类型的字典，用来存合约hash,合约编译器，文件名字
	var m1 = make(map[string]string)
	//定义value 为int 类型的字典，用来存合约更新次数，合约id
	var m2 = make(map[string]int)
	//声明一个http数据接收器
	reader, err := r.MultipartReader()
	//根据当前时间戳来创建文件夹，用来存放合约作者要上传的合约源文件
	pathFile:=createDateDir("./")
	if err != nil {
		fmt.Println("stop here")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	// 读取作者上传的文件以及ContractHash,CompilerVersion等数据，并保存在map中。
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		//fmt.Printf("FileName =[%S], FormName=[%s]\n", part.FileName(), part.FormName())

		if part.FileName()== "" {
			data, _ := ioutil.ReadAll(part)
			//fmt.Printf("FormName=[%s] FormData=[%s]\n",part.FormName(), string(data))
			//fmt.Println(part.FormName())
			if part.FormName() ==  "Contract" {
				m1[part.FormName()] = string(data)
				//fmt.Println(m1)
			} else if part.FormName() == "Version" {
				m1[part.FormName()] = string(data)
				//fmt.Println(m1)
			} else {
				//fmt.Println("map storage error")
			}
		} else {
			//dst,_ :=os.Create("./"+part.FileName()
			dst,_:= os.OpenFile(pathFile+"/"+part.FileName(),os.O_WRONLY|os.O_CREATE,0666)
			defer dst.Close()
			io.Copy(dst,part)
			fileExt := path.Ext(pathFile+"/"+part.FileName())
			if fileExt == ".csproj" {
				point := strings.Index(part.FileName(),".")
				tmp := part.FileName()[0:point]
				m1["Filename"] = tmp
			}


		}

	}

	//编译用户上传的合约源文件，并返回编译后的.nef数据
	chainNef:=execCommand(pathFile,w,m1)
	//如果编译出错，程序不向下执行
	if chainNef == "0"||chainNef=="1"||chainNef =="2"{
		return

	}
	//向链上结点请求合约的状态，返回请求到的合约nef数据
	version,sourceNef:= getContractState(pathFile,w,m1,m2)
	//如果请求失败，程序不向下执行
	if sourceNef == "3"||sourceNef=="4"{
		return
	}
	//比较用户上传的源代码编译的.nef文件与链上存储的合约.nef数据是否相等，如果相等的话，向数据库插入数据
	if sourceNef==chainNef {
		//打开数据库配置文件
		cfg, err := OpenConfigFile()
		if err != nil {
			log.Fatal(" open file error")
		}
		//连接数据库
		ctx := context.TODO()
		co,_:=intializeMongoOnlineClient(cfg, ctx)
		rt := os.ExpandEnv("${RUNTIME}")
		//查询当前合约是否已经存在于VerifiedContract表中，参数为合约hash，合约更新次数
		filter:= bson.M{"hash":getContract(m1),"updatecounter":getUpdateCounter(m2)}
		var result *mongo.SingleResult
		if rt=="mainnet"{
			result=co.Database("neofura").Collection("VerifyContractModel").FindOne(ctx,filter)
		} else {
			result=co.Database("testneofura").Collection("VerifyContractModel").FindOne(ctx,filter)
		}

		//如果合约不存在于VerifiedContract表中，验证成功
		if result.Err() != nil {
			//在VerifyContract表中插入该合约信息
			verified:= insertVerifiedContract{getContract(m1),getId(m2),getUpdateCounter(m2)}
			var insertOne *mongo.InsertOneResult
			if rt== "mainnet" {
				insertOne, err = co.Database("neofura").Collection("VerifyContractModel").InsertOne(ctx,verified)
				fmt.Println("Connect to mainnet database")
			} else {
				insertOne, err = co.Database("testneofura").Collection("VerifyContractModel").InsertOne(ctx,verified)
				fmt.Println("connect to testnet database")
			}

			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted a verified Contract in verifyContractModel collection in"+rt+" database",insertOne.InsertedID)
			//在ContractSourceCode表中，插入上传的合约源代码。
			rd, err:= ioutil.ReadDir(pathFile+"/")
			if err != nil {
				fmt.Println(err)
			}
			for _, fi := range rd {
				if fi.IsDir(){
					continue
				} else {
					fmt.Println(fi.Name())
					file,err:= os.Open(pathFile+"/"+fi.Name())
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					fileinfo, err := file.Stat()
					if err != nil {
						log.Fatal(err)
					}
					filesize := fileinfo.Size()
					buffer := make([]byte,filesize)
					_, err = file.Read(buffer)
					if err != nil {
						log.Fatal(err)

					}

					var insertOneSourceCode *mongo.InsertOneResult
					sourceCode := insertContractSourceCode{getContract(m1),getUpdateCounter(m2),fi.Name(),string(buffer)}
					if rt=="mainnet"{
						insertOneSourceCode, err = co.Database("neofura").Collection("ContractSourceCode").InsertOne(ctx, sourceCode)
					} else {
						insertOneSourceCode, err = co.Database("testneofura").Collection("ContractSourceCode").InsertOne(ctx, sourceCode)
					}

					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Inserted a contract source code in contractSourceCode collection in "+ rt+"database",insertOneSourceCode.InsertedID)

					//fmt.Println(" registed buffer",buffer)
					//fmt.Println("bytes read :",bytesread)
					//fmt.Println("bytestream to string", string(buffer))
				}
			}
			fmt.Println("=================Insert verified contract in database===============")
			msg, _ :=json.Marshal(jsonResult{5,"Verify done and record verified contract in database!"})
			w.Header().Set("Content-Type","application/json")
			os.Rename(pathFile,getContract(m1))
			w.Write(msg)
			//如果合约存在于VerifiedContract表中，说明合约已经被验证过，不会存新的数据
		} else {
			fmt.Println("=================This contract has already been verified===============")
			msg, _ :=json.Marshal(jsonResult{6,"This contract has already been verified"})
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type","application/json")
			os.RemoveAll(pathFile)
			w.Write(msg)


		}

		////比较用户上传的源代码编译的.nef文件与链上存储的合约.nef数据是否相等，如果不等的话，返回以下内容
	} else {
		fmt.Println(getVersion(m1))
		if version != getVersion(m1) {
			fmt.Println("=================Please change your compiler version and try again===============")
			msg, _ :=json.Marshal(jsonResult{7,"Compiler version error, Compiler verison shoud be "+version})
			w.Header().Set("Content-Type","application/json")
			os.RemoveAll(pathFile)
			w.Write(msg)
		} else {
			fmt.Println("=================Your source code doesn't match the contract on bloackchain===============")
			msg, _ :=json.Marshal(jsonResult{8,"Contract Source Code Verification error!"})
			w.Header().Set("Content-Type","application/json")
			os.RemoveAll(pathFile)
			w.Write(msg)
		}


	}

}
// 根据上传文件的时间戳来命名新生成的文件夹
func createDateDir(basepath string) string  {
	folderName := time.Now().Format("20060102150405")
	fmt.Println("Create folder "+ folderName)
	folderPath := filepath.Join(basepath, folderName)
	if _,err := os.Stat(folderPath);os.IsNotExist(err){
		os.Mkdir(folderPath,0777)
		os.Chmod(folderPath,0777)
	}
	return folderPath

}
//编译用户上传的合约源码
func execCommand(pathFile string,w http.ResponseWriter,m map[string] string) string{
	//cmd := exec.Command("ls")
	//根据用户上传参数选择对应的编译器
	cmd:=exec.Command("echo")
	if getVersion(m)=="Neo.Compiler.CSharp 3.0.0"{
		cmd= exec.Command("/Users/qinzilie/flamingo-contract-swap/Swap/flamingo-contract-swap/c/nccs")
		fmt.Println("use 3.0.0 compiler")
	} else if getVersion(m)=="Neo.Compiler.CSharp 3.0.2"{
		cmd = exec.Command("/Users/qinzilie/flamingo-contract-swap/Swap/flamingo-contract-swap/b/nccs")
		fmt.Println("use 3.0.2 compiler")
	} else if getVersion(m)=="Neo.Compiler.CSharp 3.0.3" {
		cmd = exec.Command("/Users/qinzilie/flamingo-contract-swap/Swap/flamingo-contract-swap/a/nccs")
		fmt.Println("use 3.0.3 compiler")
	} else {
		fmt.Println("===============Compiler version doesn't exist==============")
		msg, _ :=json.Marshal(jsonResult{0,"Compiler version doesn't exist, please choose Neo.Compiler.CSharp 3.0.0/Neo.Compiler.CSharp 3.0.2/Neo.Compiler.CSharp 3.0.3 version"})
		w.Header().Set("Content-Type","application/json")
		w.Write(msg)
		os.RemoveAll(pathFile)
		return "0"
	}

	cmd.Dir = pathFile+"/"
	stdout,err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		fmt.Println("=============== Cmd execution failed==============")
		msg, _ :=json.Marshal(jsonResult{1,"Cmd execution failed "})
		w.Header().Set("Content-Type","application/json")
		w.Write(msg)
		os.RemoveAll(pathFile)
		return "1"

	}

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(opBytes))
	}
	_, err = os.Lstat(pathFile+"/" + "bin/sc/" + m["Filename"] + ".nef")
	if !os.IsNotExist(err) {
		f, err := ioutil.ReadFile(pathFile+"/"+"bin/sc/"+m["Filename"]+".nef")
		if err != nil {
			log.Fatal(err)
		}
		res,err :=nef.FileFromBytes(f)
		if err != nil {
			log.Fatal("error")
		}
		//fmt.Println(res.Script)
		var result = base64.StdEncoding.EncodeToString(res.Script)


		fmt.Println("===========Now is soucre code============")
		fmt.Println(result)
		return result

	} else {
		fmt.Println("============.nef file doesn't exist===========", err)
		msg, _ :=json.Marshal(jsonResult{2,".nef file doesm't exist "})
		w.Header().Set("Content-Type","application/json")
		w.Write(msg)
		os.RemoveAll(pathFile)
		return "2"

	}
	//	fmt.Println(res.Magic)
	//	fmt.Println(res.Compiler)
	//	fmt.Println(res.Header)
	//	fmt.Println(res.Tokens)
	//	fmt.Println(res.Script)
}
// 向链上结点请求合约的nef数据
func getContractState(pathFile string,w http.ResponseWriter,m1 map[string] string,m2 map[string] int) (string,string) {
	rt := os.ExpandEnv("${RUNTIME}")
	var resp *http.Response
	payload, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "getcontractstate",
		"params":  []interface{}{
			getContract(m1),
		},
		"id": 1,
	})
	if rt !="mainnet" && rt!="testnet"{
		rt = "mainnet"
	}
	fmt.Println("RPC params: ContractHash:"+getContract(m1))
	switch rt {
	case "mainnet":
		resp, err = http.Post(RPCNODEMAIN, "application/json", bytes.NewReader(payload))
		fmt.Println("Runtime is:"+rt)
	case "testnet":
		resp, err = http.Post(RPCNODETEST, "application/json", bytes.NewReader(payload))
		fmt.Println("Runtime is:"+rt)
	}

	if err != nil {
		fmt.Println("=================RPC Node doesn't exsite===============")
		msg, _ :=json.Marshal(jsonResult{3,"RPC Node doesn't exsite! "})
		w.Header().Set("Content-Type","application/json")
		w.Write(msg)
		os.RemoveAll(pathFile)
		return "","3"
	}
	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//
	//fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("response Body:", string(body))
	if gjson.Get(string(body),"error").Exists() {
		message:=gjson.Get(string(body),"error.message").String()
		fmt.Println("================="+message+"===============")
		msg, _ :=json.Marshal(jsonResult{4,message})
		w.Header().Set("Content-Type","application/json")
		w.Write(msg)
		os.RemoveAll(pathFile)
		return "","4"
	}

	nef := gjson.Get(string(body),"result.nef.script")
	version:=gjson.Get(string(body),"result.nef.compiler").String()
	updateCounter := gjson.Get(string(body),"result.updatecounter").String()
	id := gjson.Get(string(body),"result.id").String()
	m2["id"],_ =strconv.Atoi(id)
	m2["updateCounter"],_ = strconv.Atoi(updateCounter)
	//fmt.Println(base64.StdEncoding.DecodeString(sourceNef))
	fmt.Println("===============Now is ChainNode nef===============")
	fmt.Println(nef.String())
	return version,nef.String()

}

func OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("config.yml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}
//链接主网和测试网数据库
func intializeMongoOnlineClient(cfg Config, ctx context.Context) (*mongo.Client, string) {
	rt := os.ExpandEnv("${RUNTIME}")
	var clientOptions *options.ClientOptions
	var dbOnline string
	if rt != "mainnet" && rt !="testnet"{
		rt = "mainnet"
	}
	switch rt {
	case "mainnet":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_main.User + ":" + cfg.Database_main.Pass + "@" + cfg.Database_main.Host + ":" + cfg.Database_main.Port + "/" + cfg.Database_main.Database)
		dbOnline = cfg.Database_main.Database
	case "testnet":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_test.User + ":" + cfg.Database_test.Pass + "@" + cfg.Database_test.Host + ":" + cfg.Database_test.Port + "/" + cfg.Database_test.Database)
		dbOnline = cfg.Database_test.Database
	}


	clientOptions.SetMaxPoolSize(50)
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("momgo connect error")
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		log.Fatal("ping mongo error")
	}
	fmt.Println("Connect mongodb success")
	return co, dbOnline
}

func getContract(m map[string] string) string {
	return m["Contract"]
}

func getVersion(m map[string] string)  string{
	return m["Version"]
}

func getUpdateCounter(m map[string] int)  int{
	return m["updateCounter"]
}

func getId(m map[string] int)  int{
	return m["id"]
}
//监听127.0.0.1:1926端口
func main() {

	fmt.Println("Server start")
	fmt.Println("YOUR ENV IS " +os.ExpandEnv("${RUNTIME}"))
	mux := http.NewServeMux()
	mux.HandleFunc("/upload",func(writer http.ResponseWriter, request *http.Request){
		multipleFile(writer,request)
	})
	mux.Handle("/metrics", promhttp.Handler())
	handler := cors.Default().Handler(mux)
	err := http.ListenAndServe("127.0.0.1:1926", handler)
	if err != nil {
		fmt.Println("listen and server error")
	}
}
