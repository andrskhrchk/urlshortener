package shortener

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(n int) string {
	if n == 0 {
		return string(charset[0])
	}

	var result []byte

	for n > 0 {
		remainder := n % 62
		result = (append(result, charset[remainder]))
		n = n / 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
