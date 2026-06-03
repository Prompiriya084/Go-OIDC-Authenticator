package ports_authentications

type RandomNumberGenerator interface {
	ToBase64String() (string, error) // 💡 ใน Go เมธอดที่เกี่ยวกับ I/O หรือ Crypto ควรคืนค่า error ด้วยครับ
}
