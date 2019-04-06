package bc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"simplebtc/util"
	"time"
)

// 区块
type Block struct {
	Timestamp int64  // 时间戳，区块产生的时间
	Number    int64  // 区块高度（索引、ID）
	Nonce     int64  //
	Parent    []byte // 父区块哈希
	Hash      []byte // 当前区块哈希
	Data      []byte // 交易数据
}

// 创建区块
func NewBlock(num int64, parentHash []byte, data []byte) *Block {
	block := Block{
		Number:    num,
		Parent:    parentHash,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	pow := NewPOW(&block)
	hash, nonce := pow.Run() // 进行工作量证明
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

// 计算区块哈希
func (b *Block) SetHash() {
	h := util.IntToHex(b.Number) // 将整数转字节数组
	t := util.IntToHex(b.Timestamp)
	data := bytes.Join([][]byte{h, t, b.Parent, b.Data}, []byte{}) // 拼接所有字节数组
	hash := sha256.Sum256(data)                                    // 进行哈希
	b.Hash = hash[:]
}

// 生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(0, nil, []byte(data))
}

// 序列化，把区块结构序列化为字节数组([]byte)
func (b *Block) Serialize() []byte {
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)          // 创建encoder对象
	if err := encoder.Encode(b); err != nil { // 编码
		log.Panicf("serialize block failed:%v\n", err)
	}
	return data.Bytes()
}

// 反序列化，把字节数组结构化为区块
func Deserialize(data []byte) *Block {
	b := Block{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&b); err != nil {
		log.Panicf("deserialize block failed: %v\n", err)
	}
	return &b
}
