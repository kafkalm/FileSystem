package encrypt

import (
	"encoding/base64"
	"strings"
)

var (
	Ip_table = []int{
		58,50,42,34,26,18,10,2,
		60,52,44,36,28,20,12,4,
		62,54,46,38,30,22,14,6,
		64,56,48,40,32,24,16,8,
		57,49,41,33,25,17,9,1,
		59,51,43,35,27,19,11,3,
		61,53,45,37,29,21,13,5,
		63,55,47,39,31,23,15,7,
	}
	Ip_table_1 = []int{
		40,8,48,16,56,24,64,32,
		39,7,47,15,55,23,63,31,
		38,6,46,14,54,22,62,30,
		37,5,45,13,53,21,61,29,
		36,4,44,12,52,20,60,28,
		35,3,43,11,51,19,59,27,
		34,2,42,10,50,18,58,26,
		33,1,41,9,49,17,57,25,
	}
	Pc_1_table = []int{
		57,49, 41, 33, 25, 17, 9,
		1, 58, 50, 42, 34, 26, 18,
		10, 2, 59, 51, 43,35, 27,
		19, 11, 3, 60, 52, 44, 36,
		63, 55, 47, 39, 31, 23, 15,
		7, 62, 54, 46, 38, 30, 22,
		14, 6, 61, 53, 45, 37, 29,
		21, 13, 5, 28, 20, 12, 4,
	}
	Pc_2_table = []int{
		14, 17, 11, 24, 1, 5,
		3, 28, 15, 6, 21, 10,
		23, 19, 12, 4, 26, 8,
		16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55,
		30, 40, 51, 45, 33, 48,
		44, 49, 39, 56, 34, 53,
		46, 42, 50, 36, 29, 32,
	}
	S_box = [8][64]int{
		{
			14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7,
			0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8,
			4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0,
			15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13,
		},
		{
			15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10,
			3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5,
			0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15,
			13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9,
		},
		{
			10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8,
			13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1,
			13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7,
			1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12,
		},
		{
			7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15,
			13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9,
			10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4,
			3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14,
		},
		{
			2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9,
			14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6,
			4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14,
			11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3,
		},
		{
			12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11,
			10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8,
			9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6,
			4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13,
		},
		{
			4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1,
			13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6,
			1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2,
			6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12,
		},
		{
			13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7,
			1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2,
			7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8,
			2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11,
		},
	}
	P_table = []int{
		16, 7, 20, 21,
		29, 12, 28, 17,
		1, 15, 23, 26, 5,
		18, 31, 10, 2, 8,
		24, 14, 32, 27,
		3, 9, 19, 13, 30,
		6, 22, 11, 4, 25,
	}
)

// 字符串转换为 二进制8bit形式
func Str2Bit(s string) []string {
	bytes := []rune(s)
	var bits []string
	for _,b := range bytes {
		bits = append(bits,Int2Bit(b)...)
	}
	return bits
}

// 一个整数返回8位二进制字符串切片
// 如 整数 10 会返回 [ "0" "0" "0" "0" "1" "0" "1" "0"]
func Int2Bit(b rune) []string {
	var bits []string
	for i:=0 ; i<8 ; i+=1 {
		if b % 2 == 0 {
			bits = append(bits,"0")
		} else {
			bits = append(bits,"1")
		}
		b /= 2
	}
	// 导致
	for i:=0 ; i<4 ; i+=1 {
		temp := bits[i]
		bits[i] = bits[7-i]
		bits[7-i] = temp
	}
	return bits
}

// 二进制字符串转十进制整数
func Bin2Dec(bin string) int {
	num := 0
	l := len(bin)
	for i := l - 1; i >= 0; i-- {
		num += (int(bin[l-i-1]) & 0xf) << uint8(i)
	}
	return num
}

// 8个bit转成字符
func Bin2Str(bits []string) string {
	num := 0
	for i:=0 ; i < len(bits) ; i++ {
		if bits[i] == "1" {
			num += 1 << uint8(len(bits)-i-1)
		}
	}
	return string(num)
}

//奇校验
func Key(key string) []string {
	newKey := Str2Bit(key)
	var newKeySlice []string
	for i:=0 ; i<8 ; i++ {
		count := 0
		for k:=7*i ; k < 7*i + 7 ; k++ {
			if newKey[k] == "1" {
				count += 1
			}
		}
		flag := "0"
		if count % 2 == 0 {
			flag = "1"
		}
		newKeySlice = append(newKeySlice,newKey[7*i:7*i+7]...)
		newKeySlice = append(newKeySlice,flag)
	}
	return newKeySlice
}

// 明文/密文 分为 64bit 一组 并填充
func Divide_64(s string) []string {
	new_s := Str2Bit(s)
	var plaintext []string
	for i:=0 ; i < len(new_s) / 64 ; i++ {
		plaintext = append(plaintext,new_s[64*i:64*i+64]...)
	}
	m := (64 - (len(new_s[64*(len(new_s)/64):]))) % 64
	if m != 0 {
		plaintext = append(plaintext,new_s[64*(len(new_s)/64):]...)
		for i:=0 ; i < 56 - len(new_s[64*(len(new_s) / 64):]) ; i++ {
			plaintext = append(plaintext,"0")
		}
		plaintext = append(plaintext,Int2Bit(rune(m-8))...)
	}
	return plaintext
}

// 置换函数
func Display(key []string,table []int) []string {
	var key_after []string
	for i:=0 ; i < len(table) ; i++ {
		key_after = append(key_after,key[table[i]-1])
	}
	return key_after
}

// 左移函数
func LS(key []string,i int) []string {
	// 注意深拷贝
	var C,D []string
	for i:= 0 ; i < 56 ; i++ {
		if i < 28 {
			C = append(C,key[i])
		} else {
			D = append(D,key[i])
		}
	}
	if i == 1 || i == 2 || i == 9 || i == 16 {
		C = append(C[1:],C[0])
		D = append(D[1:],D[0])
	} else {
		C = append(C[2:],C[0],C[1])
		D = append(D[2:],D[0],D[1])
	}
	return append(C,D...)
}

// 密钥扩展
func KeyExpansion(key []string) [][]string {
	var key_list [][]string
	key_pc_1 := Display(key,Pc_1_table)
	key_after_ls := key_pc_1
	for i:=1 ; i < 17 ; i++ {
		key_after_ls = LS(key_after_ls,i)
		key_pc_2 := Display(key_after_ls,Pc_2_table)
		key_list = append(key_list,key_pc_2)
	}
	return key_list
}

// E变换
func Echange(R []string) []string {
	var E []string
	E = append(E,R[31])
	E = append(E,R[0:5]...)
	for i:=0 ; i < 6 ; i++ {
		E = append(E,R[3 + 4*i:9+4*i]...)
	}
	E = append(E,R[27:32]...)
	E = append(E,R[0])
	return E
}

// S-盒变换
func Schange(E []string) []string {
	var S []string
	for k:=0 ; k < 8 ; k++ {
		i := Bin2Dec(E[0+k*6] + E[5+k*6])
		j := Bin2Dec(E[1+k*6] + E[2+k*6] + E[3+k*6] + E[4+k*6])
		S = append(S,Int2Bit(rune(S_box[k][i*16+j]))[4:]...)
	}
	return S
}

// 异或操作
func XOR(k1,k2 []string) []string {
	var E []string
	for i:=0 ; i < len(k1) ; i++ {
		if k1[i] == k2[i] {
			E = append(E,"0")
		} else {
			E = append(E,"1")
		}
	}
	return E
}

// f函数
func f(R,K []string) []string {
	E := Echange(R)
	E_XOR_K := XOR(E,K)
	S := Schange(E_XOR_K)
	P := Display(S,P_table)
	return P
}

// 加密函数
func MyDESEncrypt(plaintext ,key string) string {
	var ciphertext_list []string
	var ciphertext string = ""

	new_key := Key(key)
	new_keys := KeyExpansion(new_key)

	plaintext_byte := []byte(plaintext)	//转换为字节 base64转码
	plaintext = base64.StdEncoding.EncodeToString(plaintext_byte)
	plaintext_list := Divide_64(plaintext)
	for i:=0 ; i < len(plaintext_list) / 64 ; i++ {
		plaintext_ip := Display(plaintext_list[i*64:i*64+64], Ip_table)
		var L, R [][]string
		L = append(L, plaintext_ip[0:32])
		R = append(R, plaintext_ip[32:64])
		for j := 1; j < 17; j++ {
			L = append(L, R[j-1])
			R = append(R, XOR(L[j-1], f(R[j-1], new_keys[j-1])))
		}
		plaintext_r_l := append(R[16], L[16]...)
		plaintext_r_l_1 := Display(plaintext_r_l, Ip_table_1)
		ciphertext_list = append(ciphertext_list, plaintext_r_l_1...)
	}
	for i:=0 ; i < len(ciphertext_list) ; i+=8 {
		ciphertext += Bin2Str(ciphertext_list[i:i+8])
	}
	return base64.StdEncoding.EncodeToString([]byte(ciphertext))
}

// 统计切片中字符个数
func SliceCount(slice []string,c string) int {
	num := 0
	for i:=0 ; i < len(slice);i++ {
		if slice[i] == c {
			num += 1
		}
	}
	return num
}

// 解密函数
func MyDESDecrypt(ciphertext ,key string) string {
	var plaintext_list []string
	var plaintext string = ""
	new_key := Key(key)
	new_keys := KeyExpansion(new_key)

	ciphertext_byte,_ := base64.StdEncoding.DecodeString(ciphertext)
	ciphertext = string(ciphertext_byte)
	ciphertext_list := Divide_64(ciphertext)

	for i:=0 ; i < len(ciphertext_list) / 64 ; i++ {
		ciphertext_ip := Display(ciphertext_list[i*64:i*64 + 64],Ip_table)
		var L, R [][]string
		L = append(L,ciphertext_ip[0:32])
		R = append(R,ciphertext_ip[32:64])
		for j:=1 ; j < 17 ; j++ {
			L = append(L,R[j-1])
			R = append(R,XOR(L[j-1],f(R[j-1],new_keys[16-j])))
		}
		ciphertext_r_l := append(R[16],L[16]...)
		ciphertext_r_l_1 := Display(ciphertext_r_l,Ip_table_1)
		plaintext_list = append(plaintext_list,ciphertext_r_l_1...)
	}
	// 填充0的个数
	n := Bin2Dec(strings.Join(plaintext_list[len(plaintext_list)-8:len(plaintext_list)],""))
	if n == 0 {
		// 直接去掉最后8bit
		for i:=0 ; i < len(plaintext_list) - 8 ; i+=8 {
			plaintext += Bin2Str(plaintext_list[i:i+8])
		}
	} else {
		// 统计 0 出现的个数
		num := 0
		for i:= len(plaintext_list) - 8 ; i >= 8 ;i -= 8 {
			if SliceCount(plaintext_list[i-8:i],"0") != 8 {
				break
			} else {
				num += 8
			}
		}
		// 0 的数量匹配 删去添加的0
		if num == n {
			for i:=0 ; i < len(plaintext_list) -8 - num ; i+= 8 {
				plaintext += Bin2Str(plaintext_list[i:i+8])
			}
		} else {
			for i:=0 ; i < len(plaintext_list) ; i+= 8 {
				plaintext += Bin2Str(plaintext_list[i:i+8])
			}
		}
	}
	decode_plaintext,_ := base64.StdEncoding.DecodeString(plaintext)
	return string(decode_plaintext)
}