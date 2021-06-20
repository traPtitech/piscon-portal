# PISCON-PORTAL

piscon用のポータルサイト

## ディレクトリ構成(isucon8-qualify)
```
./piscon-portal
././isucon9-qualify




## 環境変数
直下に`.env`を置いてそこに配置


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
| `AWS_SUBNET_ID` | VPCのサブネットIP |


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

```


## memo 


## docker
```=shell
 sudo apt-get update

 sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

 echo \
  "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

   sudo apt-get update

    sudo apt-get install docker-ce docker-ce-cli containerd.io
```

## clone
```
git clone https://github.com/traPtitech/piscon-portal

cd piscon-portal
git checkout aws





