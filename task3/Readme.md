go mod init github.com/local/go-etherum-studey/task3 

go get -v github.com/blocto/solana-go-sdk 

go get -v github.com/joho/godotenv
 
go mod tidy 

# Solana交易生命周期流
## 阶段一：创建与签名 (Client Side)
1. 获取 Recent Blockhash： 客户端（如您的 Go 程序）首先向 RPC节点请求一个最近区块哈希（Recent Blockhash）。这个哈希是交易的“有效期”证明。
2. 构建交易： 客户端将指令、签名账户（如 feePayer 和 alice）和 Recent Blockhash 组装成一个完整的交易对象。
3. 签名： 客户端使用私钥（例如从 Backpack 导入的私钥）对交易进行签名。签名是授权执行交易的证明。
##  阶段二：提交与转发 (RPC Node / Relayer)
1. 提交至 RPC： 客户端将签名后的交易通过 sendTransaction 方法发送给 RPC节点。
2. 转发至 TPU： RPC节点会将该交易转发到当前 Slot 的领导者验证者（Leader Validator）的交易处理单元（TPU）。这个 TPU 是验证者接收交易的入口点。
##  阶段三：处理流水线 (Validator TPU)
领导者验证者接收到交易后，会通过其 TPU 的多个阶段进行流水线处理：
1. 抓取 (Fetch): TPU 监听网络上的交易。
2. 签名验证 (SigVerify): 验证交易的数字签名是否有效。
3. 银行 (Banking): 这是核心执行阶段。
  
    a. 检查 Recent Blockhash 是否仍然有效。  
    b. 验证费用支付方 (feePayer) 是否有足够的 SOL 支付交易费用。  
    c. 执行指令： 按照交易中包含的指令（如 system.Transfer）的顺序执行，更新账户状态。
4. 写入 PoH (Proof-of-History): 交易被执行后，其结果与 Slot 编号一起被写入 Solana的历史证明 (PoH) 序列中。
## 阶段四：区块发布与共识 (Network Consensus)
1.区块生成与发布： 在当前 Slot 结束时，领导者验证者将该 Slot 期间成功处理的所有交易打包成一个区块，并将其与 PoH 条目一起广播给所有其他验证者。

2.验证与投票： 其他验证者接收到这个区块后，会重新执行交易以验证其正确性。如果验证通过，它们将使用 Tower BFT 共识机制对该区块进行投票。
##  阶段五：确认与最终性 (Finality)
客户端可以通过 RPC节点查询交易状态，状态通常分为两个等级：

`Confirmed`	交易已被领导者验证者处理，并被大多数（但不一定是压倒性多数）的验证者投票确认。速度快，但仍有轻微回滚可能。  
`Finalized`	交易已经得到了绝大多数验证者的压倒性投票（2/3 以上）确认。	交易已达成最终性，不可逆转。 这是应用程序通常需要等待的状态。

总结： Solana采用 TPU 流水线和 PoH 时间戳，将交易处理、共识和网络转发的工作并行化，从而实现了远超传统区块链的速度。

# BPF加载器工作原理
BPF 是 Berkeley Packet Filter 的缩写
Solana 使用 BPF 加载器（BPF Loader） 来管理链上的程序（智能合约）的部署和执行。尽管 Solana 的运行时现在被称为 Solana BPF（SBF），但原理上仍与 eBPF 类似。  
BPF 加载器的核心是作为一个中介，管理链上的程序代码并将其连接到 Solana 运行时（Runtime）。
## 阶段一：程序部署（加载程序）
此阶段描述如何将编译好的程序（SBF 字节码）上传到 Solana 链上。

|步骤|	流程描述|	涉及账户|
| ---------- |-----------------| ---------- |
|1.| 创建可执行账户|	客户端（开发者）通过交易调用 BPF 加载器，指示它创建一个新的程序账户（Program Account）。这个账户的 Owner 被设置为 BPF 加载器。	新程序账户 P
|2.| 存储程序数据|	客户端将程序的 SBF 字节码分批上传到另一个特殊的程序数据账户（Program Data Account）中。	程序数据账户 D
|3.| 链接（设置 Program Account）|	上传完成后，客户端发起一个交易，将步骤 1 中创建的程序账户 P 标记为“可执行”，并使其指向存储了实际代码的程序数据账户 D。	程序账户 P
|4.| 最终确定（Finalize）	|BPF 加载器完成所有设置。程序账户 P 的公钥（PublicKey）现在成为应用程序调用智能合约时使用的 Program ID。	程序账户 P

## 阶段二：程序执行（调用程序）
此阶段描述当用户发起交易调用程序时，BPF 加载器如何启用程序代码。 

|步骤|	流程描述|	涉及账户|
| ---------- |-----------------| ---------- |
|5. |交易发起 (Instruction Call)|	用户发起一个包含指令的交易。该指令的目标是步骤 4 中确定的程序账户 P 的 Program ID。	用户交易
|6. |BPF 加载器中介|	Solana 运行时（Runtime）检测到该指令是调用程序账户 P。由于 程序账户 P 的 Owner 是 BPF 加载器，运行时将执行权交给 BPF 加载器。	BPF 加载器
|7. |获取代码	BPF| 加载器通过程序账户 P 的内部引用，找到实际存储 SBF 字节码的程序数据账户 D。	程序数据账户 D
|8.| 执行与隔离|	BPF 加载器将程序代码和指令中传入的所有账户数据加载到 Solana 的 SBF 虚拟机 (VM) 中。VM 在一个安全、隔离的环境中执行字节码。	SBF 虚拟机
|9.| 状态更新|	如果程序执行成功，Solana 运行时将把 VM 中对账户数据的修改应用到链上账本，并返回控制权，完成交易。	Solana 运行时

# Solana账户存储模型对比（vs EVM）
|特性| 	Solana (账户模型)  |	EVM (单体合约模型)
| ---------- |-----------------| ---------- |
|存储结构|	数据与代码分离。拥有两个主要实体：Program Accounts（代码）和 Data Accounts（数据）。	|数据与代码耦合。一个智能合约账户内部同时包含代码和私有存储空间（Storage Trie）。
|程序状态|	程序是无状态的（Stateless）。程序本身不存储数据，它通过读写传入的数据账户来改变状态。	|合约是有状态的（Stateful）。合约的数据存储在自己的内部存储空间中。
|数据访问|	显式传入（Pull Model）。交易必须明确列出程序需要读取或写入的所有数据账户。	|隐式访问（Push Model）。合约可以直接访问并修改自身的内部存储空间。
|所有权|	每个数据账户都有一个 Owner（即一个 Program ID）。只有该 Owner Program 才能修改其拥有的数据账户。	|合约是其内部数据的唯一 Owner。
|存储成本/持久性|	租金 (Rent) 机制。账户必须维护足够的 SOL 余额（Rent）才能保留在链上。如果余额低于租金，账户可能被清除（Purged）。	|预付 Gas 机制。通过 Gas 支付一次性存储费用 (SSTORE 操作)。一旦存储，除非明确删除，否则数据永久存在（Storage Slots）。
