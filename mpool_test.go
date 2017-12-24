package Mpool

import (
	"fmt"
	"github.com/flyingtimes/Mpool"
	// 读取配置文件工具
	"github.com/go-ini/ini"
	// 日志工具
	"github.com/mkideal/log"
	"os"
	"runtime"
	"strconv"
	"time"
	"testing"
)

// job1 for simple test
type Job1 struct {
	Name           string
	NextDispatcher *Mpool.Dispatcher
}

func (j Job1) Run(nextDispatcher *Mpool.Dispatcher) {
	// do your work here
	fmt.Printf("Processing job [%s]\n",j.Name)
	time.Sleep(time.Second*1)
}
func (j Job1) GetName() string {
	return j.Name
}
func (j Job1) GetNextDispatcher() *Mpool.Dispatcher {
	return j.NextDispatcher
}


//job2 for complex test
type Job2 struct {
	Name           string
	NextDispatcher *Mpool.Dispatcher
}

func (j Job2) Run(nextDispatcher *Mpool.Dispatcher) {
	// do your work here
	fmt.Printf("Processing job [%s]\n",j.Name)
	time.Sleep(time.Second*1)
	// if want to dispatch to next dispatcher ,do below
	if (nextDispatcher!=nil){
		nextDispatcher.AddTask(Job2{
			fmt.Sprintf("%s->工序2",j.Name),
			nil,
		})
	}

}
func (j Job2) GetName() string {
	return j.Name
}
func (j Job2) GetNextDispatcher() *Mpool.Dispatcher {
	return j.NextDispatcher
}


func Test_simple(t *testing.T){
	runtime.GOMAXPROCS(4)
	dispacher:= Mpool.NewDispatcher("01", 4,false)
	dispacher.Run()
	for i := 0; i < 20; i++ {
		dispacher.AddTask(Job1{
			fmt.Sprintf("任务-[%s]", strconv.Itoa(i)),
			nil,
		})
	}
	dispacher.Close()
	t.Log("测试通过")
}
func Test_all(t *testing.T) {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("找不到配置文件：", err)
		os.Exit(1)
	}
	section, err := cfg.GetSection("main")
	if err != nil {
		fmt.Println("找不到main的配置信息：", err)
		os.Exit(1)
	}
	key, err := section.GetKey("logfile")
	if err != nil {
		fmt.Println("找不到logfile的配置信息：", err)
		os.Exit(1)
	}
	log_file_name := key.String()
	defer log.Uninit(log.InitFile(log_file_name))
	log.Info("Main started.")

	runtime.GOMAXPROCS(4)

	dispacher1 := Mpool.NewDispatcher("01", 4,true)
	dispacher2 := Mpool.NewDispatcher("02", 4,true)
	dispacher1.Run()
	dispacher2.Run()
	for i := 0; i < 30; i++ {
		dispacher1.AddTask(Job2{
			fmt.Sprintf("工序1-[%s]", strconv.Itoa(i)),
			dispacher2,
		})

	}
	dispacher1.Close()
	dispacher2.Close()
	//close(dispacher1.JobQueue)
	//close(dispacher2.JobQueue)
	log.Info("Main exit normally.")
	t.Log("测试通过")
}
