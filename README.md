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
| `AWS_SUBNET_ID` | VPCのサブネットIP |


### 謎
| key               | value        |
| ----------------- | ------------ |
| `ISUCON_PASSWORD` | 存在するが謎 |