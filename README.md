基于论文《[In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)》实现，基本可用，但不支持集群成员变更、日志压缩备份。

加入了DEBUG模式，在DEBUG模式下，leader选举以及follower端的日志接收，可通过键盘控制(输入`N`表示`否`，直接回车或其他任意字符表示`是`)

**注意**：仅测试了节点宕机的情况，尚未测试网络断开的情况
