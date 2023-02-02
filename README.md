# goStealer

- container
  - gsarray
  - gslist
  - gsmap
  - gspriority_queue
  - gsqueue
  - gsset
  - gstire
  - gstype
- encoding
  - gsbase64
  - gsbinary
- errors
  - gscode
  - gserror
- net
  - gsipv4
  - gsservice
  - gsssh
  - gstcp
  - gsudp
- os
  - gsenv
  - gstimer
    - 实现
      - 自身有个定时器是100ms就当前计数（currentTicks）+1
      - 若加入定时任务1s一次的话，则会加入优先队列（小根堆）【优先级为 currentTicks+10】
      - 自身的定时器会定时看堆顶，若currentTicks>=根顶的优先级 就弹出  调用定时任务
    - 任何的定时任务都是有误差的，在定时间隔比较大，或者并发量大，负载较高的系统中尤其明显，具体请参考：https://github.com/golang/go/issues/14410
    - 定时间隔不会考虑任务的执行时间。例如，如果一项工作需要`3`分钟才能执行完成，并且计划每隔`5`分钟运行一次，那么每次任务之间只有`2`分钟的空闲时间。
    - 需要注意的是**单例模式**运行的定时任务，任务的执行时间会影响该任务下一次执行的**开始时间**。例如：一个每间隔`1`秒执行的任务，运行耗时为`1`秒，那么在**第1秒**开始运行后，下一次任务将会在**第3秒**开始执行。因为中间有一次运行检查时发现有当前任务正在进行，因此退出等待下一次执行检查。
