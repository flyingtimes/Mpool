package Mpool

import (
	"fmt"
	"github.com/mkideal/log"
	"strconv"
)


//任务
type RunnableTask interface {
	GetName() string
	GetNextDispatcher() *Dispatcher
	Run(pp *Dispatcher)
}


//  工人
type Worker struct {
    name string //工人的名字
    WorkerPool chan chan RunnableTask//对象池
    JobChannel chan RunnableTask //通道里面拿
    quit chan bool //
}



func (w *Worker) LoopWork(){
    //开一个新的协程
    go func(){

        for{
            //注册到对象池中,
	    log.Info("woker[%s]返回任务池等待任务",w.name)
            w.WorkerPool <-w.JobChannel
            select {
            //接收到了新的任务
            case job :=<- w.JobChannel:
		    log.Info("woker[%s]接收到了任务 [%s]",w.name,job.GetName())
		    job.Run(job.GetNextDispatcher())
		    log.Info("woker[%s]完成任务 [%s]",w.name,job.GetName())
            //接收到了任务
            case <-w.quit:
								log.Info("woker[%s]退出。",w.name)
                return
            }
        }
    }()
}

func (w Worker) Stop(){
    go func(){
        w.quit <- true
    }()
}

type Dispatcher struct {
                 //WorkerPool chan JobQueue
    name string //调度的名字
    maxWorkers int //获取 调试的大小
    WorkerPool chan chan RunnableTask //注册和工人一样的通道
		JobQueue chan RunnableTask
}

func (d *Dispatcher) Run(){
    // 开始运行
    for i :=0;i<d.maxWorkers;i++{
        worker := NewWorker(d.WorkerPool,fmt.Sprintf("%s-work-%s",d.name,strconv.Itoa(i)))
        //开始工作
        worker.LoopWork()
    }
    //监控
    go d.LoopGetTask()

}

func (d *Dispatcher) LoopGetTask()  {
    for {
        select {
        case job :=<-d.JobQueue:
		log.Info("调度者[%s][%d]接收到一个工作任务 %s ",d.name, len(d.WorkerPool),job.GetName())
            // 调度者接收到一个工作任务
            go func (job RunnableTask) {
                //从现有的对象池中拿出一个
                jobChannel := <-d.WorkerPool

                jobChannel <- job

            }(job)

        default:

            //fmt.Println("ok!!")
        }

    }
}

// 新建一个工人
func NewWorker(workerPool chan chan RunnableTask,name string) Worker{
    log.Info("创建了一个worker:%s \n",name);
    return Worker{
        name:name,//工人的名字
        WorkerPool: workerPool, //工人在哪个对象池里工作,可以理解成部门
        JobChannel:make(chan RunnableTask),//工人的任务
        quit:make(chan bool),
    }
}

func NewDispatcher(dname string,maxWorkers int) *Dispatcher {
	  jq := make(chan RunnableTask,maxWorkers)
    pool :=make(chan chan RunnableTask,maxWorkers)
		log.Info("调度者(%s) 初始化完毕.", dname)
    return &Dispatcher{
        WorkerPool:pool,// 将工人放到一个池中,可以理解成一个部门中
        name:dname,//调度者的名字
        maxWorkers:maxWorkers,//这个调度者有好多个工人
				JobQueue:jq,
    }
}

