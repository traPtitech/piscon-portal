# PISCON-PORTAL

piscon用のポータルサイト

## 環境変数
| key                 | value                            |
| ------------------- | -------------------------------- |
| `ENV`               | `prod`or `""`                    |
| `ACCESS_ID`         | AWSのアクセスID                  |
| `ACCESS_SECRET_KEY` | AWSのシークレットキー            |
| `MARIADB_USERNAME`  | DBのユーザーネーム               |
| `MARIADB_PASSWORD`  | DBのパスワード                   |
| `MARIADB_DATABASE`  | DBの名前(デフォルトだと`isucon`) |

### AWSのみ
| key             | value             |
| --------------- | ----------------- |
| `AWS_SUBNET_IP` | VPCのサブネットIP |


### 謎
| key               | value        |
| ----------------- | ------------ |
| `ISUCON_PASSWORD` | 存在するが謎 |

### ISUCON9-qualify以外を使いたい場合

ここを書き換えてください
```go
func formatCommand(ip string) string {
	return fmt.Sprintf("/home/isucon/isucari/bin/benchmarker "+
		"-data-dir \"/home/isucon/isucari/initial-data\" "+
		"-payment-url \"http://172.16.0.1:5555\" "+
		"-shipment-url \"http://172.16.0.1:7000\" "+
		"-static-dir \"/home/isucon/isucari/webapp/public/static\" "+
		"-target-host \"%s\" "+
		"-target-url \"http://%s\"", ip, ip)
}

````