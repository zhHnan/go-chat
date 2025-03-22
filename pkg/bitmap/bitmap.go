package bitmap

/*
该文件实现了一个基于位图（Bitmap）的数据结构，用于高效地存储和查询布尔值（即某个元素是否存在）。通过将数据压缩到字节数组的每一位，位图可以显著节省内存空间，适合处理大规模数据集。
主要功能包括：
创建位图：通过 NewBitmap 初始化一个指定大小的位图。
设置位：通过 Set 方法将字符串 id 映射到位图中的某一位，并将其设置为 1。
检查位：通过 IsSet 方法检查字符串 id 是否在位图中标记为 1。
导出与导入：通过 Export 和 Load 方法支持位图的序列化和反序列化。
哈希计算：通过 hash 函数将字符串 id 转换为整数哈希值，用于确定其在位图中的位置。
实现细节
数据结构：
bits []byte：存储位图的实际数据，每个字节包含 8 位。
size int：表示位图的总位数（len(bits) * 8）。
核心方法：
Set 方法：
使用 hash(id) 计算字符串的哈希值。
通过取模操作确定目标位的位置（字节索引和位索引）。
使用按位或操作将目标位设置为 1。
IsSet 方法：
同样使用 hash(id) 确定目标位的位置。
使用按位与操作检查目标位是否为 1。
哈希函数：
使用种子值 131313 和逐字符累加的方式生成哈希值。
通过 hash & 0x7FFFFFFF 确保哈希值为正数，避免负数导致的错误。
*/
// Bitmap is a bitmap implementation.
type Bitmap struct {
	bits []byte
	size int
}

// NewBitmap returns a new Bitmap.
func NewBitmap(size int) *Bitmap {
	// init size
	if size == 0 {
		size = 250
	}
	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8,
	}
}

// [0,0,0,0,0,0,0,0]
func (b *Bitmap) Set(id string) {
	// 计算id在bitmap中的位置
	idx := hash(id) % b.size
	// 计算在那个字节
	byteIdx := idx / 8
	// 计算在字节中的那个 位 位置
	bitIdx := idx % 8

	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	// 计算id在bitmap中的位置
	idx := hash(id) % b.size
	// 计算在那个字节
	byteIdx := idx / 8
	// 计算在字节中的那个 位 位置
	bitIdx := idx % 8

	return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}

// Export 将 Bitmap 的字节数组导出
func (b *Bitmap) Export() []byte {
	return b.bits
}

// Import 从字节数组导入 Bitmap
func Load(data []byte) *Bitmap {
	if len(data) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{
		bits: data,
		size: len(data) * 8,
	}
}

// hash 该函数实现了一个简单的字符串哈希算法，目的是将字符串 id 转换为一个整数哈希值：
//
//	种子值的作用：seed 是一个较大的质数（131313），用于减少哈希冲突的概率，提高分布均匀性。
//	逐步计算哈希值：通过遍历字符串的每个字符，将其 ASCII 值累加到哈希值中，利用乘法和加法结合种子值，确保不同字符串生成的哈希值差异较大。
//	按位与操作：hash & 0x7FFFFFFF 将结果限制为正数（去掉符号位），避免负数哈希值的出现。
func hash(id string) int {
	seed := 131313
	hash := 0
	for _, c := range id {
		hash = hash*seed + int(c)
	}
	return hash & 0x7FFFFFFF
}
