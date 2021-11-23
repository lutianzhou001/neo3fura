package verify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"neo3fura_http/lib/cli"
	log2 "neo3fura_http/lib/log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/smartcontract/nef"
	"github.com/tidwall/gjson"
)

// T...
type T struct {
	Client *cli.T
}

//定义http应答返回格式
type jsonResult struct {
	Code int
	Msg  string
}

//定义插入VerifiedContract表的数据格式, 记录被验证的合约
type insertVerifiedContract struct {
	Hash          string
	Id            int
	Updatecounter int
}

//定义插入ContractSourceCode表的数据格式，记录被验证的合约源代码
type insertContractSourceCode struct {
	Hash          string
	Updatecounter int
	FileName      string
	Code          string
}

func (me *T) MultipleFile(w http.ResponseWriter, r *http.Request) {
	//定义value 为string 类型的字典，用来存合约hash,合约编译器，文件名字
	var m1 = make(map[string]string)
	//定义value 为int 类型的字典，用来存合约更新次数，合约id
	var m2 = make(map[string]int)
	//声明一个http数据接收器
	reader, err := r.MultipartReader()
	//根据当前时间戳来创建文件夹，用来存放合约作者要上传的合约源文件
	pathFile := createDateDir("./")
	if err != nil {
		log2.Info("Stop here")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 读取作者上传的文件以及ContractHash,CompilerVersion等数据，并保存在map中。
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			data, _ := ioutil.ReadAll(part)
			if part.FormName() == "Contract" {
				m1[part.FormName()] = string(data)
			} else if part.FormName() == "Version" {
				m1[part.FormName()] = string(data)
			} else if part.FormName() == "CompileCommand"{
				m1[part.FormName()] = string(data)
			}
		} else {
			dst, _ := os.OpenFile(pathFile+"/"+part.FileName(), os.O_WRONLY|os.O_CREATE, 0666)
			defer func(dst *os.File) {
				err := dst.Close()
				if err != nil {
					log2.Fatalf("Closing file error: %v", err)
				}
			}(dst)
			_, err := io.Copy(dst, part)
			if err != nil {
				log2.Fatalf("Copy error: %v", err)
			}
			fileExt := path.Ext(pathFile + "/" + part.FileName())
			if fileExt == ".csproj" {
				point := strings.Index(part.FileName(), ".")
				tmp := part.FileName()[0:point]
				m1["Filename"] = tmp
			}
		}
	}

	//编译用户上传的合约源文件，并返回编译后的.nef数据
	chainNef := execCommand(pathFile, w, m1)
	//如果编译出错，程序不向下执行
	if chainNef == "0" || chainNef == "1" || chainNef == "2" {
		return
	}
	//向链上结点请求合约的状态，返回请求到的合约nef数据
	version, sourceNef := getContractState(pathFile, w, m1, m2)
	//如果请求失败，程序不向下执行
	if sourceNef == "3" || sourceNef == "4" {
		return
	}
	//比较用户上传的源代码编译的.nef文件与链上存储的合约.nef数据是否相等，如果相等的话，向数据库插入数据
	if sourceNef == chainNef {
		rt := os.ExpandEnv("${RUNTIME}")
		//查询当前合约是否已经存在于VerifiedContract表中，参数为合约hash，合约更新次数
		filter := bson.M{"hash": getContract(m1), "updatecounter": getUpdateCounter(m2)}
		var result *mongo.SingleResult
		if rt == "staging" {
			result = me.Client.C_online.Database("neofura").Collection("VerifyContractModel").FindOne(context.TODO(), filter)
		} else {
			result = me.Client.C_online.Database("testneofura").Collection("VerifyContractModel").FindOne(context.TODO(), filter)
		}
		//如果合约不存在于VerifiedContract表中，验证成功
		if result.Err() != nil {
			//在VerifyContract表中插入该合约信息
			verified := insertVerifiedContract{getContract(m1), getId(m2), getUpdateCounter(m2)}
			if rt == "staging" {
				_, err := me.Client.C_online.Database("neofura").Collection("VerifyContractModel").InsertOne(context.TODO(), verified)
				if err != nil {
					log2.Fatalf("Insert to online database error: %v", err)
				}
			} else {
				_, err := me.Client.C_online.Database("testneofura").Collection("VerifyContractModel").InsertOne(context.TODO(), verified)
				if err != nil {
					log2.Fatalf("Insert to online database eror: %v", err)
				}
			}
			log2.Infof("Inserted a verified Contract in verifyContractModel collection in" + rt + " database")
			//在ContractSourceCode表中，插入上传的合约源代码。
			rd, err := ioutil.ReadDir(pathFile + "/")
			if err != nil {
				log2.Infof("ReadFile error: %v", err)
			}
			for _, fi := range rd {
				if fi.IsDir() {
					continue
				} else {
					log2.Infof("File name is: %v", fi.Name())
					file, err := os.Open(pathFile + "/" + fi.Name())
					if err != nil {
						log2.Fatalf("Open file err: %v", err)
					}
					defer func(file *os.File) {
						err := file.Close()
						if err != nil {
							log2.Fatalf("Closing file error: %v", err)
						}
					}(file)
					fileInfo, err := file.Stat()
					if err != nil {
						log2.Fatalf("Stat file err: %v", err)
					}
					fileSize := fileInfo.Size()
					buffer := make([]byte, fileSize)
					_, err = file.Read(buffer)
					if err != nil {
						log2.Fatalf("Read file err: %v", err)

					}
					sourceCode := insertContractSourceCode{getContract(m1), getUpdateCounter(m2), fi.Name(), string(buffer)}
					if rt == "staging" {
						_, err := me.Client.C_online.Database("neofura").Collection("ContractSourceCode").InsertOne(context.TODO(), sourceCode)
						if err != nil {
							log2.Fatalf("Insert to online database error: %v", err)
						}
					} else if rt == "test" {
						_, err := me.Client.C_online.Database("testneofura").Collection("ContractSourceCode").InsertOne(context.TODO(), sourceCode)
						if err != nil {
							log2.Fatalf("Insert to online database error: %v", err)
						}
					}
					if err != nil {
						log2.Fatalf("Insert database error: %v", err)
					}
					log2.Infof("Inserted a contract source code in contractSourceCode collection in " + rt + "database")
				}
			}
			log2.Info("Insert verified contract in database")
			msg, _ := json.Marshal(jsonResult{5, "Verify done and record verified contract in database!"})
			w.Header().Set("Content-Type", "application/json")
			err = os.Rename(pathFile, getContract(m1))
			if err != nil {
				log2.Fatalf("Rename file error: %v", err)
			}
			_, err = w.Write(msg)
			if err != nil {
				log2.Fatalf("Writing message error: %v", err)
			}
			//如果合约存在于VerifiedContract表中，说明合约已经被验证过，不会存新的数据
		} else {
			log2.Info("This contract has already been verified")
			msg, _ := json.Marshal(jsonResult{6, "This contract has already been verified"})
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			err := os.RemoveAll(pathFile)
			if err != nil {
				log2.Fatalf("Remove file error: %v", err)
			}
			_, err = w.Write(msg)
			if err != nil {
				log2.Fatalf("Write message error: %v", err)
			}
		}

		////比较用户上传的源代码编译的.nef文件与链上存储的合约.nef数据是否相等，如果不等的话，返回以下内容
	} else {
		log2.Infof("version: %v", getVersion(m1))
		if version != getVersion(m1) {
			log2.Info("Please change your compiler version and try again")
			msg, _ := json.Marshal(jsonResult{7, "Compiler version error, Compiler verison shoud be " + version})
			w.Header().Set("Content-Type", "application/json")
			err := os.RemoveAll(pathFile)
			if err != nil {
				log2.Fatalf("Remove file error: %v", err)
			}
			_, err = w.Write(msg)
			if err != nil {
				log2.Fatalf("Write message error: %v", err)
			}
		} else {
			log2.Info("Your source code doesn't match the contract on blockchain")
			msg, _ := json.Marshal(jsonResult{8, "Contract Source Code Verification error!"})
			w.Header().Set("Content-Type", "application/json")
			err := os.RemoveAll(pathFile)
			if err != nil {
				log2.Fatalf("Remove file error: %v", err)
			}
			_, err = w.Write(msg)
			if err != nil {
				log2.Fatalf("Write message error: %v", err)
			}
		}
	}
}

func createDateDir(basepath string) string {
	folderName := time.Now().Format("20060102150405")
	log2.Infof("Create folder: %v", folderName)
	folderPath := filepath.Join(basepath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0777)
		if err != nil {
			log2.Fatalf("Create dir error: %v", err)
		}
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			log2.Fatalf("Chmod error: %v", err)
		}
	}
	return folderPath
}

//编译用户上传的合约源码
func execCommand(pathFile string, w http.ResponseWriter, m map[string]string) string {
	//cmd := exec.Command("ls")
	//根据用户上传参数选择对应的编译器
	cmd := exec.Command("echo")
	if getVersion(m) == "Neo.Compiler.CSharp 3.0.0" {
		if getCompileCommand(m)=="nccs --no-optimize" {
			cmd= exec.Command("/go/application/compiler/a/nccs","--no-optimize")
			log2.Infof("Compiler: Neo.Compiler.CSharp 3.0.0, Command: nccs --no-optimize")
		}
		if getCompileCommand(m)=="nccs" {
			cmd= exec.Command("/go/application/compiler/a/nccs")
			log2.Infof("Compiler: Neo.Compiler.CSharp 3.0.0, Command: nccs")
		}
	} else if getVersion(m) == "Neo.Compiler.CSharp 3.0.2" {
		if getCompileCommand(m)=="nccs --no-optimize" {
			cmd= exec.Command("/go/application/compiler/c/nccs","--no-optimize")
			log2.Infof("Compiler: Neo.Compiler.CSharp 3.0.2, Command: nccs --no-optimize")
		}
		if getCompileCommand(m)=="nccs" {
			cmd= exec.Command("/go/application/compiler/c/nccs")
			log2.Infof("Compiler: Neo.Compiler.CSharp 3.0.2, Command: nccs")
		}
	} else if getVersion(m) == "Neo.Compiler.CSharp 3.0.3" {
		if getCompileCommand(m)=="nccs --no-optimize" {
			cmd= exec.Command("/go/application/compiler/b/nccs","--no-optimize")
			log2.Infof("Compiler: Neo.Compiler.CSharp 3.0.3, Command: nccs --no-optimize")
		}
		if getCompileCommand(m)=="nccs" {
			cmd= exec.Command("/go/application/compiler/b/nccs")
			log2.Infof("Compiler: Neo.Compiler.CSharp 3.0.3, Command: nccs")
		}
	} else {
		log2.Fatalf("Compiler version doesn't exist")
		msg, _ := json.Marshal(jsonResult{0, "Compiler version doesn't exist, please choose Neo.Compiler.CSharp 3.0.0/Neo.Compiler.CSharp 3.0.2/Neo.Compiler.CSharp 3.0.3 version"})
		w.Header().Set("Content-Type", "application/json")
		err := os.RemoveAll(pathFile)
		if err != nil {
			log2.Fatalf("Remove file error: %v", err)
		}
		_, err = w.Write(msg)
		if err != nil {
			log2.Fatalf("Write message error: %v", err)
		}
		return "0"
	}

	cmd.Dir = pathFile + "/"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log2.Fatalf("StdoutPipe error: %v", err)
	}
	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
			log2.Fatalf("Closing reader error: %v", err)
		}
	}(stdout)

	err = cmd.Start()
	if err != nil {
		log2.Infof("Cmd execution failed: %v", err)
		msg, _ := json.Marshal(jsonResult{1, "Cmd execution failed "})
		w.Header().Set("Content-Type", "application/json")
		// err := os.RemoveAll(pathFile)
		// if err != nil {
		// log2.Fatalf("Remove file error: %v", err)
		// }
		_, err = w.Write(msg)
		if err != nil {
			log2.Fatalf("Write message error: %v", err)
		}
		return "1"
	}

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log2.Fatalf("Read error: %v", err)
	} else {
		log2.Info(string(opBytes))
	}
	_, err = os.Lstat(pathFile + "/" + "bin/sc/" + m["Filename"] + ".nef")
	if !os.IsNotExist(err) {
		f, err := ioutil.ReadFile(pathFile + "/" + "bin/sc/" + m["Filename"] + ".nef")
		if err != nil {
			log2.Fatalf("Read file error: %v", err)
		}
		res, err := nef.FileFromBytes(f)
		if err != nil {
			log2.Fatalf("File from bytes error: %v", err)
		}
		var result = base64.StdEncoding.EncodeToString(res.Script)
		return result
	} else {
		log2.Fatalf(".nef file doesn't exist: %v", err)
		msg, _ := json.Marshal(jsonResult{2, ".nef file doesn't exist "})
		w.Header().Set("Content-Type", "application/json")
		//err := os.RemoveAll(pathFile)
		//if err != nil {
		//	log2.Fatalf("Remove file error: %v", err)
		//}
		_, err = w.Write(msg)
		if err != nil {
			log2.Fatalf("Write message error: %v", err)
		}
		return "2"
	}
}

func getContract(m map[string]string) string {
	return m["Contract"]
}

func getVersion(m map[string]string) string {
	return m["Version"]
}

func getUpdateCounter(m map[string]int) int {
	return m["updateCounter"]
}

func getId(m map[string]int) int {
	return m["id"]
}
func getCompileCommand(m map[string] string) string{
	return m["CompileCommand"]
}

// 向链上结点请求合约的nef数据
func getContractState(pathFile string, w http.ResponseWriter, m1 map[string]string, m2 map[string]int) (string, string) {
	rt := os.ExpandEnv("${RUNTIME}")
	var resp *http.Response
	params := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "getcontractstate",
		"params": []interface{}{
			getContract(m1),
		},
		"id": 1,
	}
	payload, err := json.Marshal(params)
	log2.Infof("RPC params: ContractHash: %v", getContract(m1))
	switch rt {
	case "staging":
		resp, err = http.Post("https://neofura.ngd.network", "application/json", bytes.NewReader(payload))
	case "test":
		resp, err = http.Post("https://testneofura.ngd.network:444", "application/json", bytes.NewReader(payload))
	}

	if err != nil {
		log2.Infof("RPC node error :%v", err)
		msg, _ := json.Marshal(jsonResult{3, "RPC Node error "})
		w.Header().Set("Content-Type", "application/json")
		err := os.RemoveAll(pathFile)
		if err != nil {
			log2.Fatalf("Remove file error: %v", err)
		}
		_, err = w.Write(msg)
		if err != nil {
			log2.Fatalf("Write message error: %v", err)
		}
		return "", "3"
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log2.Fatalf("Closing reader error: %v", err)
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	if gjson.Get(string(body), "error").Exists() {
		message := gjson.Get(string(body), "error.message").String()
		log2.Infof("")
		msg, _ := json.Marshal(jsonResult{4, message})
		w.Header().Set("Content-Type", "application/json")
		err := os.RemoveAll(pathFile)
		if err != nil {
			log2.Fatalf("Remove file error: %v", err)
		}
		_, err = w.Write(msg)
		if err != nil {
			log2.Fatalf("Write message error: %v", err)
		}
		return "", "4"
	}

	nef := gjson.Get(string(body), "result.nef.script")
	version := gjson.Get(string(body), "result.nef.compiler").String()
	updateCounter := gjson.Get(string(body), "result.updatecounter").String()
	id := gjson.Get(string(body), "result.id").String()
	m2["id"], _ = strconv.Atoi(id)
	m2["updateCounter"], _ = strconv.Atoi(updateCounter)
	log2.Infof("Now is chain node nef :%v", nef.String())
	return version, nef.String()
}
