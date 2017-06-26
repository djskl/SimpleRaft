基于论文《[In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)》实现，可以保证**Election Safety**, **Leader Append-Only**, **Log Matching**, **Leader Completeness**, **State Machine Safety**，但不支持集群成员变更、日志压缩备份。

**注意** : Follower中加入了随机延迟来模拟网络延迟
